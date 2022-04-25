package ubuntu

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
		{name: "ubuntu:22.04", fields: fields{
			releaseString: "22.04",
			mockFile:      "testdata/SHA256SUM",
		}, want: want{
			artifactDetails: &api.ArtifactDetails{
				SHA256Sum:            "de5e632e17b8965f2baf4ea6d2b824788e154d9a65df4fd419ec4019898e15cd",
				DownloadURL:          "https://cloud-images.ubuntu.com/releases/22.04/release/ubuntu-22.04-server-cloudimg-amd64.img",
				AdditionalUniqueTags: []string{"de5e632e17b8965f2baf4ea6d2b824788e154d9a65df4fd419ec4019898e15cd"},
			},
			metadata: &api.Metadata{
				Name:                   "ubuntu",
				Version:                "22.04",
				Description:            description,
				ExampleUserDataPayload: docs.CloudInit(&docs.UserData{}),
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
