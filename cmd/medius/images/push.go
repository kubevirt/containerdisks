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
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/ulikunitz/xz"
	"kubevirt.io/containerdisks/cmd/medius/common"
	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/build"
	"kubevirt.io/containerdisks/pkg/http"
	"kubevirt.io/containerdisks/pkg/repository"
)

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

			resultsChan, err := spawnWorkers(cmd.Context(), options, func(e *common.Entry) (*api.ArtifactResult, error) {
				errString := ""
				tags, err := buildAndPublish(cmd.Context(), e, options, time.Now())
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

			if err != nil {
				if options.PublishImagesOptions.NoFail {
					logrus.Warn(err)
				} else {
					logrus.Fatal(err)
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

func buildAndPublish(ctx context.Context, entry *common.Entry, options *common.Options, timestamp time.Time) ([]string, error) {
	description := entry.Artifact.Metadata().Describe()
	log := common.Logger(entry.Artifact)

	imageName := path.Join(options.PublishImagesOptions.SourceRegistry, description)
	artifactInfo, err := entry.Artifact.Inspect()
	if err != nil {
		return nil, fmt.Errorf("error introspecting artifact %q: %v", description, err)
	}
	log.Infof("Remote artifact checksum: %q", artifactInfo.SHA256Sum)
	repo := repository.RepositoryImpl{}
	imageSha := ""
	imageInfo, err := repo.ImageMetadata(imageName, options.AllowInsecureRegistry)
	if err != nil {
		if repository.IsRepositoryUnknownError(err) {
			log.Info("Repository does not yet exist, it will be created")
		} else if repository.IsManifestUnknownError(err) {
			log.Info("Tag does not yet exist, it will be created")
		} else if repository.IsTagUnknownError(err) {
			log.Info("Tag is gone but seems to have existed already, it will be created")
		} else {
			return nil, fmt.Errorf("error introspecting image %q: %v", imageName, err)
		}
	} else {
		log.Infof("Latest containerdisk checksum: %q", imageInfo.Labels[build.LabelShaSum])
		imageSha = imageInfo.Labels[build.LabelShaSum]
	}
	if artifactInfo.SHA256Sum == imageSha && !options.PublishImagesOptions.ForceBuild {
		log.Info("Nothing to do.")
		return nil, nil
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return nil, ctx.Err()
	}

	log.Infof("Rebuild needed, downloading %q ...", artifactInfo.DownloadURL)
	getter := &http.HTTPGetter{}
	artifactReader, err := getter.GetWithChecksum(artifactInfo.DownloadURL)
	if err != nil {
		return nil, fmt.Errorf("error opening a connection to the specified download location: %v", err)
	}
	defer artifactReader.Close()

	var reader io.Reader = artifactReader
	if artifactInfo.Compression == types.GzipAlgorithmName {
		reader, err = gzip.NewReader(artifactReader)
		if err != nil {
			return nil, fmt.Errorf("error creating a gunzip reader for the specified download location: %v", err)
		}
	} else if artifactInfo.Compression == types.XzAlgorithmName {
		reader, err = xz.NewReader(artifactReader)
		if err != nil {
			return nil, fmt.Errorf("error creating a lzma reader for the specified download location: %v", err)
		}
	}

	file, err := os.CreateTemp("", "containerdisks")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	if _, err := io.Copy(file, reader); err != nil {
		return nil, fmt.Errorf("error writing the image to the destination file: %v", err)
	}
	file.Close()
	if errors.Is(ctx.Err(), context.Canceled) {
		return nil, ctx.Err()
	}

	checksum := artifactReader.Checksum()
	if checksum != artifactInfo.SHA256Sum {
		return nil, fmt.Errorf("expected checksum %q but got %q", artifactInfo.SHA256Sum, checksum)
	}
	log.Info("Building containerdisk ...")
	containerDisk, err := build.BuildContainerDisk(file.Name(), checksum)
	if err != nil {
		return nil, fmt.Errorf("error creating the containerdisk : %v", err)
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return nil, ctx.Err()
	}

	names := prepareTags(timestamp, options.PublishImagesOptions.TargetRegistry, entry, artifactInfo)
	for _, name := range names {
		if !options.DryRun {
			log.Infof("Pushing %s", name)
			if err := repo.PushImage(ctx, containerDisk, name); err != nil {
				log.WithError(err).Error("Failed to push image")
				return nil, err
			}
		} else {
			log.Infof("Dry run enabled, not pushing %s", name)
		}

		if errors.Is(ctx.Err(), context.Canceled) {
			return nil, ctx.Err()
		}
	}

	return prepareTags(timestamp, "", entry, artifactInfo), nil
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
