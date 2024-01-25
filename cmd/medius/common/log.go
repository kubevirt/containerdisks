package common

import (
	"github.com/sirupsen/logrus"

	"kubevirt.io/containerdisks/pkg/api"
)

func Logger(artifact api.Artifact) *logrus.Entry {
	metadata := artifact.Metadata()
	return logrus.WithFields(
		logrus.Fields{
			"name":    metadata.Name,
			"version": metadata.Version,
		},
	)
}
