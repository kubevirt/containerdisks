package generic

import (
	v1 "kubevirt.io/api/core/v1"
	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/docs"
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

func (c *generic) VM(name, imgRef, _ string) *v1.VirtualMachine {
	return docs.BasicVM(
		name,
		imgRef,
	)
}

func (c *generic) UserData(_ *docs.UserData) string {
	return ""
}

func New(artifactDetails *api.ArtifactDetails, metadata *api.Metadata) *generic {
	return &generic{artifactDetails: artifactDetails, metadata: metadata}
}
