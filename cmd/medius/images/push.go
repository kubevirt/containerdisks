package images

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/containers/image/v5/pkg/compression/types"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/ulikunitz/xz"
	"kubevirt.io/containerdisks/cmd/medius/common"
	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/build"
	"kubevirt.io/containerdisks/pkg/http"
	"kubevirt.io/containerdisks/pkg/repository"
)

type buildAndPublish struct {
	Ctx     context.Context
	Log     *logrus.Entry
	Options *common.Options
	Repo    repository.Repository
	Getter  http.Getter
}

func NewPublishImagesCommand(options *common.Options) *cobra.Command {
	options.PublishImagesOptions = common.PublishImageOptions{
		SourceRegistry: "quay.io/containerdisks",
	}

	publishCmd := &cobra.Command{
		Use:   "push",
		Short: "Determine if containerdisks need an update and push an update to the target registry if needed",
		Run: func(cmd *cobra.Command, args []string) {
			if options.PublishImagesOptions.TargetRegistry == "" {
				options.PublishImagesOptions.TargetRegistry = options.PublishImagesOptions.SourceRegistry
			}

			resultsChan, workerErr := spawnWorkers(cmd.Context(), options, func(e *common.Entry) (*api.ArtifactResult, error) {
				errString := ""

				b := buildAndPublish{
					Ctx:     cmd.Context(),
					Log:     common.Logger(e.Artifact),
					Options: options,
					Repo:    &repository.RepositoryImpl{},
					Getter:  &http.HTTPGetter{},
				}
				tags, err := b.Do(e, time.Now())
				if err != nil {
					errString = err.Error()
				}

				if tags == nil && err == nil {
					return nil, nil
				}

				return &api.ArtifactResult{
					Tags:  tags,
					Stage: StagePush,
					Err:   errString,
				}, err
			})

			results := map[string]api.ArtifactResult{}
			for result := range resultsChan {
				results[result.Key] = result.Value
			}

			if !options.DryRun {
				if err := writeResultsFile(options.ImagesOptions.ResultsFile, results); err != nil {
					logrus.Fatal(err)
				}
			}

			if workerErr != nil {
				if options.PublishImagesOptions.NoFail {
					logrus.Warn(workerErr)
				} else {
					logrus.Fatal(workerErr)
				}
			}
		},
	}
	publishCmd.Flags().BoolVar(&options.PublishImagesOptions.ForceBuild, "force",
		options.PublishImagesOptions.ForceBuild, "Force a rebuild and push")
	publishCmd.Flags().BoolVar(&options.PublishImagesOptions.NoFail, "no-fail",
		options.PublishImagesOptions.NoFail, "Return success even if a worker fails")
	publishCmd.Flags().StringVar(&options.PublishImagesOptions.SourceRegistry, "source-registry",
		options.PublishImagesOptions.SourceRegistry, "Registry to check if updates are needed")
	publishCmd.Flags().StringVar(&options.PublishImagesOptions.TargetRegistry, "target-registry",
		options.PublishImagesOptions.TargetRegistry, "Registry to push built containerdisks to")

	return publishCmd
}

func (b *buildAndPublish) Do(entry *common.Entry, timestamp time.Time) ([]string, error) {
	description := entry.Artifact.Metadata().Describe()
	artifactInfo, err := entry.Artifact.Inspect()
	if err != nil {
		return nil, fmt.Errorf("error introspecting artifact %q: %v", description, err)
	}
	b.Log.Infof("Remote artifact checksum: %q", artifactInfo.SHA256Sum)

	imageSha, err := b.getImageSha(description)
	if err != nil {
		return nil, err
	}
	if imageSha == artifactInfo.SHA256Sum && !b.Options.PublishImagesOptions.ForceBuild {
		b.Log.Info("Nothing to do.")
		return nil, nil
	}
	if errors.Is(b.Ctx.Err(), context.Canceled) {
		return nil, b.Ctx.Err()
	}

	b.Log.Infof("Rebuild needed, downloading %q ...", artifactInfo.DownloadURL)
	file, err := b.getArtifact(artifactInfo)
	if err != nil {
		return nil, err
	}
	defer os.Remove(file)

	b.Log.Info("Building containerdisk ...")
	containerDisk, err := build.ContainerDisk(file, artifactInfo.SHA256Sum)
	if err != nil {
		return nil, fmt.Errorf("error creating the containerdisk : %v", err)
	}
	if errors.Is(b.Ctx.Err(), context.Canceled) {
		return nil, b.Ctx.Err()
	}

	names := prepareTags(timestamp, b.Options.PublishImagesOptions.TargetRegistry, entry, artifactInfo)
	for _, name := range names {
		if err := b.pushImage(containerDisk, name); err != nil {
			return nil, err
		}
		if errors.Is(b.Ctx.Err(), context.Canceled) {
			return nil, b.Ctx.Err()
		}
	}

	return prepareTags(timestamp, "", entry, artifactInfo), nil
}

func (b *buildAndPublish) getImageSha(description string) (imageSha string, err error) {
	imageName := path.Join(b.Options.PublishImagesOptions.SourceRegistry, description)
	imageInfo, err := b.Repo.ImageMetadata(imageName, b.Options.AllowInsecureRegistry)
	if err != nil {
		err = b.handleMetadataError(imageName, err)
	} else {
		b.Log.Infof("Latest containerdisk checksum: %q", imageInfo.Labels[build.LabelShaSum])
		imageSha = imageInfo.Labels[build.LabelShaSum]
	}

	return
}

func (b *buildAndPublish) handleMetadataError(imageName string, err error) error {
	switch {
	case repository.IsRepositoryUnknownError(err):
		b.Log.Info("Repository does not yet exist, it will be created")
	case repository.IsManifestUnknownError(err):
		b.Log.Info("Tag does not yet exist, it will be created")
	case repository.IsTagUnknownError(err):
		b.Log.Info("Tag is gone but seems to have existed already, it will be created")
	default:
		return fmt.Errorf("error introspecting image %q: %v", imageName, err)
	}

	return nil
}

func (b *buildAndPublish) getArtifact(artifactInfo *api.ArtifactDetails) (string, error) {
	artifactReader, err := b.Getter.GetWithChecksumAndContext(b.Ctx, artifactInfo.DownloadURL)
	if err != nil {
		return "", fmt.Errorf("error opening a connection to the specified download location: %v", err)
	}
	defer artifactReader.Close()

	file, err := b.readArtifact(artifactReader, artifactInfo.Compression)
	if err != nil {
		return "", err
	}
	if errors.Is(b.Ctx.Err(), context.Canceled) {
		return "", b.Ctx.Err()
	}

	checksum := artifactReader.Checksum()
	if checksum != artifactInfo.SHA256Sum {
		return "", fmt.Errorf("expected checksum %q but got %q", artifactInfo.SHA256Sum, checksum)
	}

	return file, nil
}

func (b *buildAndPublish) readArtifact(artifactReader http.ReadCloserWithChecksum, compression string) (string, error) {
	var err error

	// Initialize reader with the artifactReader for the case where no compression is used
	var reader io.Reader = artifactReader
	if compression == types.GzipAlgorithmName {
		reader, err = gzip.NewReader(artifactReader)
		if err != nil {
			return "", fmt.Errorf("error creating a gunzip reader for the specified download location: %v", err)
		}
	} else if compression == types.XzAlgorithmName {
		reader, err = xz.NewReader(artifactReader)
		if err != nil {
			return "", fmt.Errorf("error creating a lzma reader for the specified download location: %v", err)
		}
	}

	file, err := os.CreateTemp("", "containerdisks")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Uncompress disks in chunks up to size defined below
	const chunkSize = 1024 * 1024 * 50 // MiB
	for {
		_, err := io.CopyN(file, reader, chunkSize)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("error writing the image to the destination file: %v", err)
		}
		if errors.Is(b.Ctx.Err(), context.Canceled) {
			return "", b.Ctx.Err()
		}
	}

	return file.Name(), nil
}

func (b *buildAndPublish) pushImage(containerDisk v1.Image, name string) error {
	if !b.Options.DryRun {
		b.Log.Infof("Pushing %s", name)
		if err := b.Repo.PushImage(b.Ctx, containerDisk, name); err != nil {
			b.Log.WithError(err).Error("Failed to push image")
			return err
		}
	} else {
		b.Log.Infof("Dry run enabled, not pushing %s", name)
	}

	return nil
}

func prepareTags(timestamp time.Time, registry string, entry *common.Entry, artifactDetails *api.ArtifactDetails) []string {
	metadata := entry.Artifact.Metadata()
	imageName := path.Join(registry, metadata.Describe())

	names := []string{fmt.Sprintf("%s-%s", imageName, timestamp.Format("0601021504"))}
	for _, tag := range artifactDetails.AdditionalUniqueTags {
		if tag == "" {
			continue
		}
		names = append(names, fmt.Sprintf("%s:%s", path.Join(registry, metadata.Name), tag))
	}
	// the least specific tag is last
	names = append(names, imageName)

	if entry.UseForLatest {
		names = append(names, fmt.Sprintf("%s:%s", path.Join(registry, metadata.Name), "latest"))
	}

	return names
}
