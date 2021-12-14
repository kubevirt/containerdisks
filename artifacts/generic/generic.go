package generic

import (
	"kubevirt.io/containerdisks/pkg/api"
)

type generic struct {
	artifactDetails *api.ArtifactDetails
	metadata        *api.Metadata
}

func (c *generic) Metadata() *api.Metadata {
	return c.metadata
}

func (c *generic) Inspect() (*api.ArtifactDetails, error) {
	return c.artifactDetails, nil
}

func New(artifactDetails *api.ArtifactDetails, metadata *api.Metadata) *generic {
	return &generic{artifactDetails: artifactDetails, metadata: metadata}
}
