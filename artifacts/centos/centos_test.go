package centos

import (
	"testing"

	"kubevirt.io/containerdisks/pkg/api"
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
		{name: "centos:8.4", fields: fields{
			releaseString: "8.4",
			mockFile:      "testdata/centos8.checksum",
		}, want: want{
			artifactDetails: &api.ArtifactDetails{
				SHA256Sum:            "3510fc7deb3e1939dbf3fe6f65a02ab1efcc763480bc352e4c06eca2e4f7c2a2",
				DownloadURL:          "https://cloud.centos.org/centos/8/x86_64/images/CentOS-8-GenericCloud-8.4.2105-20210603.0.x86_64.qcow2",
				Compression:          "",
				AdditionalUniqueTags: []string{"8.4.2105-20210603.0", "8.4.2105"},
			},
			metadata: &api.Metadata{
				Name:    "centos",
				Version: "8.4",
			},
		},
			wantErr: false,
		},
		{name: "centos:8.3", fields: fields{
			releaseString: "8.3",
			mockFile:      "testdata/centos8.checksum",
		}, want: want{
			artifactDetails: &api.ArtifactDetails{
				SHA256Sum:            "7ec97062618dc0a7ebf211864abf63629da1f325578868579ee70c495bed3ba0",
				DownloadURL:          "https://cloud.centos.org/centos/8/x86_64/images/CentOS-8-GenericCloud-8.3.2011-20201204.2.x86_64.qcow2",
				Compression:          "",
				AdditionalUniqueTags: []string{"8.3.2011-20201204.2", "8.3.2011"},
			},
			metadata: &api.Metadata{
				Name:    "centos",
				Version: "8.3",
			},
		},
			wantErr: false,
		},
		{name: "centos:7-2009", fields: fields{
			releaseString: "7-2009",
			mockFile:      "testdata/centos7.checksum",
		}, want: want{
			artifactDetails: &api.ArtifactDetails{
				SHA256Sum:   "e38bab0475cc6d004d2e17015969c659e5a308111851b0e2715e84646035bdd3",
				DownloadURL: "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud-2009.qcow2",
				Compression: "",
			},
			metadata: &api.Metadata{
				Name:    "centos",
				Version: "7-2009",
			},
		},
			wantErr: false,
		},
		{name: "ceontos:7-1809", fields: fields{
			releaseString: "7-1809",
			mockFile:      "testdata/centos7.checksum",
		}, want: want{
			artifactDetails: &api.ArtifactDetails{
				SHA256Sum:   "42c062df8a8c36991ec0282009dd52ac488461a3f7ee114fc21a765bfc2671c2",
				DownloadURL: "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud-1809.qcow2",
				Compression: "",
			},
			metadata: &api.Metadata{
				Name:    "centos",
				Version: "7-1809",
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
