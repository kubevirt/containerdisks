package fedora

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/common"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/pkg/http"
	"kubevirt.io/containerdisks/testutil"
)

var _ = Describe("Fedora", func() {
	DescribeTable("Inspect should be able to parse releases files",
		func(release, mockFile string, details *api.ArtifactDetails, envVariables map[string]string, metadata *api.Metadata) {
			c := New(release, envVariables)
			c.getter = testutil.NewMockGetter(mockFile)
			got, err := c.Inspect()
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(details))
			Expect(c.Metadata()).To(Equal(metadata))
		},
		Entry("fedora:35", "35", "testdata/releases.json",
			&api.ArtifactDetails{
				SHA256Sum:            "fe84502779b3477284a8d4c86731f642ca10dd3984d2b5eccdf82630a9ca2de6",
				DownloadURL:          "https://download.fedoraproject.org/pub/fedora/linux/releases/35/Cloud/x86_64/images/Fedora-Cloud-Base-35-1.2.x86_64.qcow2", //nolint:lll
				AdditionalUniqueTags: []string{"35-1.2"},
				ImageArchitecture:    "amd64",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.small",
				common.DefaultPreferenceEnv:   "fedora",
			},
			&api.Metadata{
				Name:        "fedora",
				Version:     "35",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "fedora",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.small",
					common.DefaultPreferenceEnv:   "fedora",
				},
			},
		),
		Entry("fedora:34", "34", "testdata/releases.json",
			&api.ArtifactDetails{
				SHA256Sum:            "b9b621b26725ba95442d9a56cbaa054784e0779a9522ec6eafff07c6e6f717ea",
				DownloadURL:          "https://download.fedoraproject.org/pub/fedora/linux/releases/34/Cloud/x86_64/images/Fedora-Cloud-Base-34-1.2.x86_64.qcow2", //nolint:lll
				AdditionalUniqueTags: []string{"34-1.2"},
				ImageArchitecture:    "amd64",
			},
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.small",
				common.DefaultPreferenceEnv:   "fedora",
			},
			&api.Metadata{
				Name:        "fedora",
				Version:     "34",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "fedora",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.small",
					common.DefaultPreferenceEnv:   "fedora",
				},
			},
		),
	)

	It("Gather should be able to parse releases files", func() {
		artifacts := [][]api.Artifact{
			{
				&fedora{
					Version: "36",
					Arch:    "x86_64",
					Variant: "Cloud",
					getter:  &http.HTTPGetter{},
					EnvVariables: map[string]string{
						common.DefaultInstancetypeEnv: "u1.small",
						common.DefaultPreferenceEnv:   "fedora",
					},
				},
			},
			{
				&fedora{
					Version: "35",
					Arch:    "x86_64",
					Variant: "Cloud",
					getter:  &http.HTTPGetter{},
					EnvVariables: map[string]string{
						common.DefaultInstancetypeEnv: "u1.small",
						common.DefaultPreferenceEnv:   "fedora",
					},
				},
			},
		}

		c := NewGatherer()
		c.getter = testutil.NewMockGetter("testdata/releases.json")
		got, err := c.Gather()
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal(artifacts))
	})
})

func TestFedora(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fedora Suite")
}
