package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"kubevirt.io/containerdisks/artifacts/fedora"
	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/build"
	"kubevirt.io/containerdisks/pkg/http"
	"kubevirt.io/containerdisks/pkg/repository"
)

var registry = map[string]api.Artifact{
	fedora.NewFedora("35").Metadata().Describe(): fedora.NewFedora("35"),
	//fedora.NewFedora("34").Metadata().Describe(): fedora.NewFedora("34"), # example
	//fedora.NewFedora("33").Metadata().Describe(): fedora.NewFedora("33"), # example
}

type PublishOptions struct {
	ForceBuild bool
	Focus      string
	Workers    int
}

func NewPublishCommand(options *Options) *cobra.Command {
	options.PublishOptions = PublishOptions{
		ForceBuild: false,
		Focus:      "",
		Workers:    1,
	}

	publishCmd := &cobra.Command{
		Use:   "publish",
		Short: "Determine if containerdisks need an update and publish all which need one",
		Run: func(cmd *cobra.Command, args []string) {
			errChan := make(chan error, options.PublishOptions.Workers)
			jobChan := make(chan api.Artifact, options.PublishOptions.Workers)

			wg := &sync.WaitGroup{}
			wg.Add(options.PublishOptions.Workers)
			for x := 0; x < options.PublishOptions.Workers; x++ {
				go worker(wg, jobChan, options, errChan)
			}

			for desc := range registry {
				if options.PublishOptions.Focus != "" && options.PublishOptions.Focus != desc {
					continue
				}
				jobChan <- registry[desc]
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
	publishCmd.Flags().BoolVar(&options.PublishOptions.ForceBuild, "force", options.PublishOptions.ForceBuild, "Force a rebuild and push")
	publishCmd.Flags().StringVar(&options.PublishOptions.Focus, "focus", options.PublishOptions.Focus, "Only build a specific containerdisk")
	publishCmd.Flags().IntVar(&options.PublishOptions.Workers, "workers", options.PublishOptions.Workers, "Number of parallel workers")

	return publishCmd
}

func worker(wg *sync.WaitGroup, job chan api.Artifact, options *Options, errChan chan error) {
	defer wg.Done()
	for a := range job {
		if err := buildAndPublish(a, options, time.Now()); err != nil {
			logger(a).Error(err)
			errChan <- err
		}
	}
}

func buildAndPublish(artifact api.Artifact, options *Options, timestamp time.Time) error {
	metadata := artifact.Metadata()
	log := logger(artifact)

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
		} else {
			return fmt.Errorf("error introspecting image %q: %v", imageName, err)
		}
	} else {
		log.Infof("Latest containerdisk checksum: %q", imageInfo.Labels["shasum"])
		imageSha = imageInfo.Labels["shasum"]
	}
	if artifactInfo.SHA256Sum == imageSha && !options.PublishOptions.ForceBuild {
		log.Info("Nothing to do.")
		return nil
	}
	log.Infof("Rebuild needed, downloading %q ...", artifactInfo.DownloadURL)
	getter := &http.HTTPGetter{}
	reader, err := getter.GetWithChecksum(artifactInfo.DownloadURL)
	if err != nil {
		return fmt.Errorf("error opening a connection to the specified download location: %v", err)
	}
	defer reader.Close()

	file, err := ioutil.TempFile("", "containerdisks")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("error writing the image to the destination file: %v", err)
	}
	file.Close()
	checksum := reader.Checksum()

	if checksum != artifactInfo.SHA256Sum {
		return fmt.Errorf("expected checksum %q but got %q", artifactInfo.SHA256Sum, checksum)
	}
	log.Info("Building containerdisk ...")
	containerDisk, err := build.BuildContainerDisk(file.Name(), checksum)
	if err != nil {
		return fmt.Errorf("error creating the containerdisk : %v", err)
	}
	names := []string{fmt.Sprintf("%s-%s", imageName, timestamp.Format("0601021504")), imageName}
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

func logger(artifact api.Artifact) *logrus.Entry {
	metadata := artifact.Metadata()
	return logrus.WithFields(
		logrus.Fields{
			"name":    metadata.Name,
			"version": metadata.Version,
		},
	)
}
