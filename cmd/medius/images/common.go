package images

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"kubevirt.io/containerdisks/cmd/medius/common"
	"kubevirt.io/containerdisks/pkg/api"
)

const (
	StagePush    = "push"
	StageVerify  = "verify"
	StagePromote = "promote"
)

type workerResult struct {
	Key   string
	Value api.ArtifactResult
}

func spawnWorkers(ctx context.Context, options *common.Options, workerFn func(api.Artifact) (*api.ArtifactResult, error)) (chan workerResult, error) {
	count := len(common.Registry)
	errChan := make(chan error, count)
	jobChan := make(chan api.Artifact, count)
	resultsChan := make(chan workerResult, count)
	defer close(resultsChan)

	if options.ImagesOptions.Workers > count {
		logrus.Warnf("Limiting workers to number of artifacts: %d", count)
		options.ImagesOptions.Workers = count
	}

	wg := &sync.WaitGroup{}
	wg.Add(options.ImagesOptions.Workers)
	for x := 0; x < options.ImagesOptions.Workers; x++ {
		go func() {
			defer wg.Done()
			for a := range jobChan {
				result, err := workerFn(a)
				if result != nil {
					resultsChan <- workerResult{
						Key:   a.Metadata().Describe(),
						Value: *result,
					}
				}
				if err != nil && !errors.Is(err, context.Canceled) {
					common.Logger(a).Error(err)
					errChan <- err
				}
				if errors.Is(ctx.Err(), context.Canceled) {
					return
				}
			}
		}()
	}

	fillJobChan(jobChan, options.Focus)
	close(jobChan)

	wg.Wait()

	select {
	case err := <-errChan:
		return resultsChan, err
	default:
		return resultsChan, nil
	}
}

func fillJobChan(jobChan chan api.Artifact, focus string) {
	for i, desc := range common.Registry {
		if focus == "" && desc.SkipWhenNotFocused {
			continue
		}

		if focus != "" && focus != desc.Artifact.Metadata().Describe() {
			continue
		}

		jobChan <- common.Registry[i].Artifact
	}
}

func writeResultsFile(fileName string, results map[string]api.ArtifactResult) error {
	logrus.Info("Writing results file")

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func readResultsFile(fileName string) (map[string]api.ArtifactResult, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	results := map[string]api.ArtifactResult{}
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}

	return results, nil
}
