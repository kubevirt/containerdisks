package rhcos

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
			releaseString: "4.9",
			mockFile:      "testdata/rhcos-4.9.checksum",
		}, want: want{
			artifactDetails: &api.ArtifactDetails{
				SHA256Sum:   "3466690807fb710102559ea57daac0484c59ed4d914996882d601b8bb7a7ada8",
				DownloadURL: "https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/4.9/latest/rhcos-openstack.x86_64.qcow2.gz",
				Compression: "gzip",
			},
			metadata: &api.Metadata{
				Name:                    "rhcos",
				Version:                 "4.9",
				ExampleCloudInitPayload: docs.Ignition(),
				Description:             description,
			},
		},
			wantErr: false,
		},
		{name: "rhcos:4.8", fields: fields{
			releaseString: "4.8",
			mockFile:      "testdata/rhcos-4.8.checksum",
		}, want: want{
			artifactDetails: &api.ArtifactDetails{
				SHA256Sum:   "99da4ed945b391d452e55a3a7809c799e4c74f69acbee1ecaec78f368c4e369e",
				DownloadURL: "https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/4.8/latest/rhcos-openstack.x86_64.qcow2.gz",
				Compression: "gzip",
			},
			metadata: &api.Metadata{
				Name:                    "rhcos",
				Version:                 "4.8",
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
