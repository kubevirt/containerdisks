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

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/ulikunitz/xz"
	"go.podman.io/image/v5/pkg/compression/types"

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

			focusMatched, resultsChan, workerErr := spawnWorkers(cmd.Context(), options, func(e *common.Entry) (*api.ArtifactResult, error) {
				errString := ""
				artifact := e.Artifacts[0]

				b := buildAndPublish{
					Ctx:     cmd.Context(),
					Log:     common.Logger(artifact),
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

			if !focusMatched {
				logrus.Fatalf("no artifact was processed, focus '%s' did not match", options.Focus)
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
	metadata := entry.Artifacts[0].Metadata()
	artifactInfo, err := entry.Artifacts[0].Inspect()
	if err != nil {
		return nil, fmt.Errorf("error introspecting artifact %q: %v", metadata.Describe(), err)
	}

	rebuildNeeded, err := b.rebuildNeeded(entry)
	if err != nil {
		return nil, err
	}
	if !rebuildNeeded && !b.Options.PublishImagesOptions.ForceBuild {
		b.Log.Info("Nothing to do.")
		return nil, nil
	}
	if errors.Is(b.Ctx.Err(), context.Canceled) {
		return nil, b.Ctx.Err()
	}

	images, artifacts, err := b.buildImages(entry)
	if err != nil {
		return nil, err
	}
	defer cleanupArtifacts(artifacts)

	names := prepareTags(timestamp, b.Options.PublishImagesOptions.TargetRegistry, entry, artifactInfo)
	for _, name := range names {
		if len(images) > 1 {
			containerDiskIndex, err := build.ContainerDiskIndex(images)
			if err != nil {
				return nil, fmt.Errorf("error creating the containerdisk index : %v", err)
			}
			if err := b.pushImageIndex(containerDiskIndex, name); err != nil {
				return nil, err
			}
		} else if len(images) == 1 {
			if err := b.pushImage(images[0], name); err != nil {
				return nil, err
			}
		}
		if errors.Is(b.Ctx.Err(), context.Canceled) {
			return nil, b.Ctx.Err()
		}
	}

	return prepareTags(timestamp, "", entry, artifactInfo), nil
}

func (b *buildAndPublish) getImageChecksum(description, arch string) (imageChecksum string, err error) {
	imageName := path.Join(b.Options.PublishImagesOptions.SourceRegistry, description)
	imageInfo, err := b.Repo.ImageMetadata(imageName, arch, b.Options.AllowInsecureRegistry)
	if err != nil {
		err = b.handleMetadataError(imageName, err)
	} else {
		b.Log.Infof("Latest containerdisk checksum: %q", imageInfo.Labels[build.LabelShaSum])
		imageChecksum = imageInfo.Labels[build.LabelShaSum]
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
	case repository.IsArchUnknownError(err):
		b.Log.Info("Image with arch does not exist yet, it will be created")
	default:
		return fmt.Errorf("error introspecting image %q: %v", imageName, err)
	}

	return nil
}

func (b *buildAndPublish) getArtifact(artifactInfo *api.ArtifactDetails) (string, error) {
	artifactReader, err := b.getArtifactReader(artifactInfo)
	if err != nil {
		return "", err
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
	if checksum != artifactInfo.Checksum {
		return "", fmt.Errorf("expected checksum %q but got %q", artifactInfo.Checksum, checksum)
	}

	return file, nil
}

func (b *buildAndPublish) getArtifactReader(artifactInfo *api.ArtifactDetails) (http.ReadCloserWithChecksum, error) {
	var artifactReader http.ReadCloserWithChecksum
	var err error
	const retries = 3
	for range retries {
		artifactReader, err = b.Getter.GetWithChecksumAndContext(b.Ctx, artifactInfo.DownloadURL, artifactInfo.ChecksumHash)
		if err == nil {
			return artifactReader, nil
		}
		b.Log.Infof("Artifact download verification failed, retrying...")
	}
	return nil, fmt.Errorf("error opening a connection to the specified download location: %v", err)
}

func (b *buildAndPublish) readArtifact(artifactReader http.ReadCloserWithChecksum, compression string) (string, error) {
	var err error

	// Initialize reader with the artifactReader for the case where no compression is used
	var reader io.Reader = artifactReader

	switch compression {
	case types.GzipAlgorithmName:
		reader, err = gzip.NewReader(artifactReader)
		if err != nil {
			return "", fmt.Errorf("error creating a gunzip reader for the specified download location: %v", err)
		}
	case types.XzAlgorithmName:
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

func (b *buildAndPublish) buildImages(entry *common.Entry) ([]v1.Image, []string, error) {
	var images []v1.Image
	var artifacts []string

	for i := range entry.Artifacts {
		metadata := entry.Artifacts[i].Metadata()
		artifactInfo, err := entry.Artifacts[i].Inspect()
		if err != nil {
			return nil, nil, fmt.Errorf("error introspecting artifact %q: %v", metadata.Describe(), err)
		}

		b.Log.Infof("Rebuild needed, downloading %q ...", artifactInfo.DownloadURL)
		file, err := b.getArtifact(artifactInfo)
		if err != nil {
			return nil, nil, err
		}
		artifacts = append(artifacts, file)

		b.Log.Info("Building containerdisk ...")
		image, err := build.ContainerDisk(file,
			artifactInfo.ImageArchitecture,
			build.ContainerDiskConfig(artifactInfo.Checksum, metadata.EnvVariables))
		if err != nil {
			return nil, nil, fmt.Errorf("error creating the containerdisk : %v", err)
		}
		if errors.Is(b.Ctx.Err(), context.Canceled) {
			return nil, nil, b.Ctx.Err()
		}
		images = append(images, image)
	}

	return images, artifacts, nil
}

func (b *buildAndPublish) rebuildNeeded(entry *common.Entry) (bool, error) {
	if len(entry.Artifacts) == 0 {
		err := errors.New("entry has no artifacts to check for rebuild")
		b.Log.Error(err)
		return false, err
	}

	for i := range entry.Artifacts {
		metadata := entry.Artifacts[i].Metadata()
		artifactInfo, err := entry.Artifacts[i].Inspect()
		if err != nil {
			return false, fmt.Errorf("error introspecting artifact %q: %v", metadata.Describe(), err)
		}
		imageChecksum, err := b.getImageChecksum(metadata.Describe(), artifactInfo.ImageArchitecture)
		if err != nil {
			return false, err
		}
		if imageChecksum != artifactInfo.Checksum {
			return true, nil
		}
	}

	return false, nil
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

func (b *buildAndPublish) pushImageIndex(containerDiskIndex v1.ImageIndex, name string) error {
	if !b.Options.DryRun {
		b.Log.Infof("Pushing %s image index", name)
		if err := b.Repo.PushImageIndex(b.Ctx, containerDiskIndex, name); err != nil {
			b.Log.WithError(err).Error("Failed to push image image")
			return err
		}
	} else {
		b.Log.Infof("Dry run enabled, not pushing %s image index", name)
	}

	return nil
}

func prepareTags(timestamp time.Time, registry string, entry *common.Entry, artifactDetails *api.ArtifactDetails) []string {
	metadata := entry.Artifacts[0].Metadata()
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

func cleanupArtifacts(artifacts []string) {
	for _, file := range artifacts {
		os.Remove(file)
	}
}
