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
		func(release, arch, mockFile string, details *api.ArtifactDetails, metadata *api.Metadata) {
			c := New(release, arch)
			c.getter = testutil.NewMockGetter(mockFile)
			got, err := c.Inspect()
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(details))
			Expect(c.Metadata()).To(Equal(metadata))
		},
		Entry("fedora:35 x86_64", "35", "x86_64", "testdata/releases.json",
			&api.ArtifactDetails{
				SHA256Sum:            "fe84502779b3477284a8d4c86731f642ca10dd3984d2b5eccdf82630a9ca2de6",
				DownloadURL:          "https://download.fedoraproject.org/pub/fedora/linux/releases/35/Cloud/x86_64/images/Fedora-Cloud-Base-35-1.2.x86_64.qcow2", //nolint:lll
				AdditionalUniqueTags: []string{"35-1.2"},
				ImageArchitecture:    "amd64",
			},
			&api.Metadata{
				Name:        "fedora",
				Version:     "35",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "fedora",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: defaultInstancetypeX86_64,
					common.DefaultPreferenceEnv:   defaultPreferenceX86_64,
				},
			},
		),
		Entry("fedora:35 aarch64", "35", "aarch64", "testdata/releases.json",
			&api.ArtifactDetails{
				SHA256Sum:            "c71f2e6ce75b516d565e2c297ea9994c69b946cb3eaa0a4bbea400dbd6f59ae6",
				DownloadURL:          "https://download.fedoraproject.org/pub/fedora/linux/releases/35/Cloud/aarch64/images/Fedora-Cloud-Base-35-1.2.aarch64.qcow2", //nolint:lll
				AdditionalUniqueTags: []string{"35-1.2"},
				ImageArchitecture:    "arm64",
			},
			&api.Metadata{
				Name:        "fedora",
				Version:     "35",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "fedora",
				},
			},
		),
		Entry("fedora:34 x86_64", "34", "x86_64", "testdata/releases.json",
			&api.ArtifactDetails{
				SHA256Sum:            "b9b621b26725ba95442d9a56cbaa054784e0779a9522ec6eafff07c6e6f717ea",
				DownloadURL:          "https://download.fedoraproject.org/pub/fedora/linux/releases/34/Cloud/x86_64/images/Fedora-Cloud-Base-34-1.2.x86_64.qcow2", //nolint:lll
				AdditionalUniqueTags: []string{"34-1.2"},
				ImageArchitecture:    "amd64",
			},
			&api.Metadata{
				Name:        "fedora",
				Version:     "34",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "fedora",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: defaultInstancetypeX86_64,
					common.DefaultPreferenceEnv:   defaultPreferenceX86_64,
				},
			},
		),
		Entry("fedora:34 aarch64", "34", "aarch64", "testdata/releases.json",
			&api.ArtifactDetails{
				SHA256Sum:            "141f16f52bfbe159947267658a0dbfbbe96fd5b988a95d1271f9c9ed61156da2",
				DownloadURL:          "https://download.fedoraproject.org/pub/fedora/linux/releases/34/Cloud/aarch64/images/Fedora-Cloud-Base-34-1.2.aarch64.qcow2", //nolint:lll
				AdditionalUniqueTags: []string{"34-1.2"},
				ImageArchitecture:    "arm64",
			},
			&api.Metadata{
				Name:        "fedora",
				Version:     "34",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "fedora",
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
						common.DefaultInstancetypeEnv: defaultInstancetypeX86_64,
						common.DefaultPreferenceEnv:   defaultPreferenceX86_64,
					},
				},
				&fedora{
					Version: "36",
					Arch:    "aarch64",
					Variant: "Cloud",
					getter:  &http.HTTPGetter{},
				},
			},
			{
				&fedora{
					Version: "35",
					Arch:    "x86_64",
					Variant: "Cloud",
					getter:  &http.HTTPGetter{},
					EnvVariables: map[string]string{
						common.DefaultInstancetypeEnv: defaultInstancetypeX86_64,
						common.DefaultPreferenceEnv:   defaultPreferenceX86_64,
					},
				},
				&fedora{
					Version: "35",
					Arch:    "aarch64",
					Variant: "Cloud",
					getter:  &http.HTTPGetter{},
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
