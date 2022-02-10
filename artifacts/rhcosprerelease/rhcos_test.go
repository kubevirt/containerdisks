package rhcosprerelease

import (
	"testing"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/testutil"

	. "github.com/onsi/gomega"
)

func Test_Inspect(t *testing.T) {
	type fields struct {
		releaseString string
		mockFile      string
	}
	type want struct {
		artifactDetails *api.ArtifactDetails
		metadata        *api.Metadata
	}
	tests := []struct {
		name    string
		fields  fields
		want    want
		wantErr bool
	}{
		{name: "rhcos:4.9", fields: fields{
			releaseString: "latest-4.9",
			mockFile:      "testdata/rhcos-latest-4.9-prerelease.checksum",
		}, want: want{
			artifactDetails: &api.ArtifactDetails{
				SHA256Sum:            "3466690807fb710102559ea57daac0484c59ed4d914996882d601b8bb7a7ada8",
				DownloadURL:          "https://mirror.openshift.com/pub/openshift-v4/x86_64/dependencies/rhcos/pre-release/latest-4.9/rhcos-openstack.x86_64.qcow2.gz",
				Compression:          "gzip",
				AdditionalUniqueTags: []string{"4.9.0-rc.7", "3466690807fb710102559ea57daac0484c59ed4d914996882d601b8bb7a7ada8"},
			},
			metadata: &api.Metadata{
				Name:                    "rhcos",
				Version:                 "4.9-pre-release",
				ExampleCloudInitPayload: docs.Ignition(),
				Description:             description,
			},
		},
			wantErr: false,
		},
		{name: "rhcos:latest", fields: fields{
			releaseString: "latest",
			mockFile:      "testdata/rhcos-latest-prerelease.checksum",
		}, want: want{
			artifactDetails: &api.ArtifactDetails{
				SHA256Sum:            "f581896eee37216021bfce9ddd5e1fd8289c366ca0d1db25221c77688de85fd7",
				DownloadURL:          "https://mirror.openshift.com/pub/openshift-v4/x86_64/dependencies/rhcos/pre-release/latest/rhcos-openstack.x86_64.qcow2.gz",
				Compression:          "gzip",
				AdditionalUniqueTags: []string{"4.10.0-rc.1", "f581896eee37216021bfce9ddd5e1fd8289c366ca0d1db25221c77688de85fd7"},
			},
			metadata: &api.Metadata{
				Name:                    "rhcos",
				Version:                 "latest-pre-release",
				ExampleCloudInitPayload: docs.Ignition(),
				Description:             description,
			},
		},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			c := New(tt.fields.releaseString)
			c.getter = testutil.NewMockGetter(tt.fields.mockFile)
			got, err := c.Inspect()
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).NotTo(HaveOccurred())
			}
			g.Expect(got).To(Equal(tt.want.artifactDetails))
			g.Expect(c.Metadata()).To(Equal(tt.want.metadata))
		})
	}
}
