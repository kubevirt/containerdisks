package fedora

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
		{name: "fedora:35", fields: fields{
			releaseString: "35",
			mockFile:      "testdata/release.json",
		}, want: want{
			artifactDetails: &api.ArtifactDetails{
				SHA256Sum:   "fe84502779b3477284a8d4c86731f642ca10dd3984d2b5eccdf82630a9ca2de6",
				DownloadURL: "https://download.fedoraproject.org/pub/fedora/linux/releases/35/Cloud/x86_64/images/Fedora-Cloud-Base-35-1.2.x86_64.qcow2",
				Compression: "",
			},
			metadata: &api.Metadata{
				Name:    "fedora",
				Version: "35",
			},
		},
			wantErr: false,
		},
		{name: "fedora:34", fields: fields{
			releaseString: "34",
			mockFile:      "testdata/release.json",
		}, want: want{
			artifactDetails: &api.ArtifactDetails{
				SHA256Sum:   "b9b621b26725ba95442d9a56cbaa054784e0779a9522ec6eafff07c6e6f717ea",
				DownloadURL: "https://download.fedoraproject.org/pub/fedora/linux/releases/34/Cloud/x86_64/images/Fedora-Cloud-Base-34-1.2.x86_64.qcow2",
				Compression: "",
			},
			metadata: &api.Metadata{
				Name:    "fedora",
				Version: "34",
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
