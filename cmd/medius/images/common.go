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

func spawnWorkers(ctx context.Context, o *common.Options, fn func(*common.Entry) (*api.ArtifactResult, error)) (chan workerResult, error) {
	registry := common.NewRegistry()
	count := len(registry)
	errChan := make(chan error, count)
	jobChan := make(chan *common.Entry, count)
	resultsChan := make(chan workerResult, count)
	defer close(resultsChan)

	if o.ImagesOptions.Workers > count {
		logrus.Warnf("Limiting workers to number of artifacts: %d", count)
		o.ImagesOptions.Workers = count
	}

	wg := &sync.WaitGroup{}
	wg.Add(o.ImagesOptions.Workers)
	for x := 0; x < o.ImagesOptions.Workers; x++ {
		go func() {
			defer wg.Done()
			for e := range jobChan {
				result, err := fn(e)
				if result != nil {
					resultsChan <- workerResult{
						Key:   e.Artifact.Metadata().Describe(),
						Value: *result,
					}
				}
				if err != nil && !errors.Is(err, context.Canceled) {
					common.Logger(e.Artifact).Error(err)
					errChan <- err
				}
				if errors.Is(ctx.Err(), context.Canceled) {
					return
				}
			}
		}()
	}

	fillJobChan(jobChan, registry, o.Focus)
	close(jobChan)

	wg.Wait()

	select {
	case err := <-errChan:
		return resultsChan, err
	default:
		return resultsChan, nil
	}
}

func fillJobChan(jobChan chan *common.Entry, registry []common.Entry, focus string) {
	for i, desc := range registry {
		if focus == "" && desc.SkipWhenNotFocused {
			continue
		}

		if focus != "" && focus != desc.Artifact.Metadata().Describe() {
			continue
		}

		jobChan <- &registry[i]
	}
}

func writeResultsFile(fileName string, results map[string]api.ArtifactResult) error {
	logrus.Info("Writing results file")

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	const permissionUserReadWrite = 0600
	err = os.WriteFile(fileName, data, permissionUserReadWrite)
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
