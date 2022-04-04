package centosstream

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
		{name: "centos-stream:8", fields: fields{
			releaseString: "8",
			mockFile:      "testdata/centos-stream8.checksum",
		}, want: want{
			artifactDetails: &api.ArtifactDetails{
				SHA256Sum:            "8e22e67687b81e38c7212fc30c47cb24cbc4935c0f2459ed139f498397d1e7cd",
				DownloadURL:          "https://cloud.centos.org/centos/8-stream/x86_64/images/CentOS-Stream-GenericCloud-8-20210603.0.x86_64.qcow2",
				Compression:          "",
				AdditionalUniqueTags: []string{"8-20210603.0"},
			},
			metadata: &api.Metadata{
				Name:                   "centos-stream",
				Version:                "8",
				Description:            description,
				ExampleUserDataPayload: docs.CloudInit(&docs.UserData{}),
			},
		},
			wantErr: false,
		},
		{name: "centos-stream:9", fields: fields{
			releaseString: "9",
			mockFile:      "testdata/centos-stream9.checksum",
		}, want: want{
			artifactDetails: &api.ArtifactDetails{
				SHA256Sum:            "bcebdc00511d6e18782732570056cfbc7cba318302748bfc8f66be9c0db68142",
				DownloadURL:          "https://cloud.centos.org/centos/9-stream/x86_64/images/CentOS-Stream-GenericCloud-9-20211222.0.x86_64.qcow2",
				Compression:          "",
				AdditionalUniqueTags: []string{"9-20211222.0"},
			},
			metadata: &api.Metadata{
				Name:                   "centos-stream",
				Version:                "9",
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
