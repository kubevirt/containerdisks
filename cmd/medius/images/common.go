package images

import (
	"context"
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
	"kubevirt.io/containerdisks/cmd/medius/common"
	"kubevirt.io/containerdisks/pkg/api"
)

func spawnWorkers(ctx context.Context, options *common.Options, workerFn func(api.Artifact) error) error {
	count := len(common.Registry)
	errChan := make(chan error, count)
	jobChan := make(chan api.Artifact, count)

	if options.Workers > count {
		logrus.Warnf("Limiting workers to number of artifacts: %d", count)
		options.Workers = count
	}

	wg := &sync.WaitGroup{}
	wg.Add(options.Workers)
	for x := 0; x < options.Workers; x++ {
		go func() {
			defer wg.Done()
			for a := range jobChan {
				if err := workerFn(a); err != nil {
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
		return err
	default:
		return nil
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
