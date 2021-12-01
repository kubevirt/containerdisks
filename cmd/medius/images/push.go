package images

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sync"
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
		ForceBuild: false,
		Focus:      "",
		Workers:    1,
	}

	publishCmd := &cobra.Command{
		Use:   "push",
		Short: "Determine if containerdisks need an update and push an update to the target registry if needed",
		Run: func(cmd *cobra.Command, args []string) {
			errChan := make(chan error, options.PublishImagesOptions.Workers)
			jobChan := make(chan api.Artifact, options.PublishImagesOptions.Workers)

			wg := &sync.WaitGroup{}
			wg.Add(options.PublishImagesOptions.Workers)
			for x := 0; x < options.PublishImagesOptions.Workers; x++ {
				go worker(wg, jobChan, options, errChan)
			}

			for i, desc := range common.Registry {
				if options.PublishImagesOptions.Focus != "" && options.PublishImagesOptions.Focus != desc.Artifact.Metadata().Describe() {
					continue
				}
				jobChan <- common.Registry[i].Artifact
			}
			close(jobChan)

			wg.Wait()

			select {
			case <-errChan:
				os.Exit(1)
			default:
				os.Exit(0)
			}
		},
	}
	publishCmd.Flags().BoolVar(&options.PublishImagesOptions.ForceBuild, "force", options.PublishImagesOptions.ForceBuild, "Force a rebuild and push")
	publishCmd.Flags().StringVar(&options.PublishImagesOptions.Focus, "focus", options.PublishImagesOptions.Focus, "Only build a specific containerdisk")
	publishCmd.Flags().IntVar(&options.PublishImagesOptions.Workers, "workers", options.PublishImagesOptions.Workers, "Number of parallel workers")

	return publishCmd
}

func worker(wg *sync.WaitGroup, job chan api.Artifact, options *common.Options, errChan chan error) {
	defer wg.Done()
	for a := range job {
		if err := buildAndPublish(a, options, time.Now()); err != nil {
			common.Logger(a).Error(err)
			errChan <- err
		}
	}
}

func buildAndPublish(artifact api.Artifact, options *common.Options, timestamp time.Time) error {
	metadata := artifact.Metadata()
	log := common.Logger(artifact)

	imageName := path.Join(options.Registry, metadata.Describe())
	artifactInfo, err := artifact.Inspect()
	if err != nil {
		return fmt.Errorf("error introspecting artifact %q: %v", metadata.Describe(), err)
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
			return fmt.Errorf("error introspecting image %q: %v", imageName, err)
		}
	} else {
		log.Infof("Latest containerdisk checksum: %q", imageInfo.Labels["shasum"])
		imageSha = imageInfo.Labels["shasum"]
	}
	if artifactInfo.SHA256Sum == imageSha && !options.PublishImagesOptions.ForceBuild {
		log.Info("Nothing to do.")
		return nil
	}
	log.Infof("Rebuild needed, downloading %q ...", artifactInfo.DownloadURL)
	getter := &http.HTTPGetter{}
	artifactReader, err := getter.GetWithChecksum(artifactInfo.DownloadURL)
	if err != nil {
		return fmt.Errorf("error opening a connection to the specified download location: %v", err)
	}
	defer artifactReader.Close()

	var reader io.Reader = artifactReader
	if artifactInfo.Compression == types.GzipAlgorithmName {
		reader, err = gzip.NewReader(artifactReader)
		if err != nil {
			return fmt.Errorf("error creating a gunzip reader for the specified download location: %v", err)
		}
	} else if artifactInfo.Compression == types.XzAlgorithmName {
		reader, err = xz.NewReader(artifactReader)
		if err != nil {
			return fmt.Errorf("error creating a lzma reader for the specified download location: %v", err)
		}
	}

	file, err := ioutil.TempFile("", "containerdisks")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("error writing the image to the destination file: %v", err)
	}
	file.Close()
	checksum := artifactReader.Checksum()

	if checksum != artifactInfo.SHA256Sum {
		return fmt.Errorf("expected checksum %q but got %q", artifactInfo.SHA256Sum, checksum)
	}
	log.Info("Building containerdisk ...")
	containerDisk, err := build.BuildContainerDisk(file.Name(), checksum)
	if err != nil {
		return fmt.Errorf("error creating the containerdisk : %v", err)
	}
	names := prepareTags(timestamp, options.Registry, metadata, artifactInfo)
	for _, name := range names {
		logrus.Infof("Pushing %q", name)
		if !options.DryRun {
			if err := build.PushImage(context.Background(), containerDisk, name); err != nil {
				return err
			}
		}
	}

	return nil
}

func prepareTags(timestamp time.Time, registry string, metadata *api.Metadata, artifactDetails *api.ArtifactDetails) []string {
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
	return names
}
