package fedora

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/pkg/http"
	"kubevirt.io/containerdisks/pkg/instancetype"
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
		Entry("fedora:40 x86_64", "40", "x86_64", "testdata/releases.json",
			&api.ArtifactDetails{
				SHA256Sum:            "ac58f3c35b73272d5986fa6d3bc44fd246b45df4c334e99a07b3bbd00684adee",
				DownloadURL:          "https://download.fedoraproject.org/pub/fedora/linux/releases/40/Cloud/x86_64/images/Fedora-Cloud-Base-Generic.x86_64-40-1.14.qcow2", //nolint:lll
				AdditionalUniqueTags: []string{"40-1.14"},
				ImageArchitecture:    "amd64",
			},
			&api.Metadata{
				Name:        "fedora",
				Version:     "40",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "fedora",
				},
				EnvVariables: map[string]string{
					instancetype.DefaultInstancetypeEnv: defaultInstancetypeX86_64,
					instancetype.DefaultPreferenceEnv:   defaultPreferenceX86_64,
				},
			},
		),
		Entry("fedora:40 aarch64", "40", "aarch64", "testdata/releases.json",
			&api.ArtifactDetails{
				SHA256Sum:            "ebdce26d861a9d15072affe1919ed753ec7015bd97b3a7d0d0df6a10834f7459",
				DownloadURL:          "https://download.fedoraproject.org/pub/fedora/linux/releases/40/Cloud/aarch64/images/Fedora-Cloud-Base-Generic.aarch64-40-1.14.qcow2", //nolint:lll
				AdditionalUniqueTags: []string{"40-1.14"},
				ImageArchitecture:    "arm64",
			},
			&api.Metadata{
				Name:        "fedora",
				Version:     "40",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "fedora",
				},
			},
		),
		Entry("fedora:39 x86_64", "39", "x86_64", "testdata/releases.json",
			&api.ArtifactDetails{
				SHA256Sum:            "ab5be5058c5c839528a7d6373934e0ce5ad6c8f80bd71ed3390032027da52f37",
				DownloadURL:          "https://download.fedoraproject.org/pub/fedora/linux/releases/39/Cloud/x86_64/images/Fedora-Cloud-Base-39-1.5.x86_64.qcow2", //nolint:lll
				AdditionalUniqueTags: []string{"39-1.5"},
				ImageArchitecture:    "amd64",
			},
			&api.Metadata{
				Name:        "fedora",
				Version:     "39",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "fedora",
				},
				EnvVariables: map[string]string{
					instancetype.DefaultInstancetypeEnv: defaultInstancetypeX86_64,
					instancetype.DefaultPreferenceEnv:   defaultPreferenceX86_64,
				},
			},
		),
		Entry("fedora:39 aarch64", "39", "aarch64", "testdata/releases.json",
			&api.ArtifactDetails{
				SHA256Sum:            "765996d5b77481ca02d0ac06405641bf134ac920cfc1e60d981c64d7971162dc",
				DownloadURL:          "https://download.fedoraproject.org/pub/fedora/linux/releases/39/Cloud/aarch64/images/Fedora-Cloud-Base-39-1.5.aarch64.qcow2", //nolint:lll
				AdditionalUniqueTags: []string{"39-1.5"},
				ImageArchitecture:    "arm64",
			},
			&api.Metadata{
				Name:        "fedora",
				Version:     "39",
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
					Version: "40",
					Arch:    "x86_64",
					Variant: "Cloud",
					getter:  &http.HTTPGetter{},
					EnvVariables: map[string]string{
						instancetype.DefaultInstancetypeEnv: defaultInstancetypeX86_64,
						instancetype.DefaultPreferenceEnv:   defaultPreferenceX86_64,
					},
				},
				&fedora{
					Version: "40",
					Arch:    "aarch64",
					Variant: "Cloud",
					getter:  &http.HTTPGetter{},
				},
			},
			{
				&fedora{
					Version: "39",
					Arch:    "x86_64",
					Variant: "Cloud",
					getter:  &http.HTTPGetter{},
					EnvVariables: map[string]string{
						instancetype.DefaultInstancetypeEnv: defaultInstancetypeX86_64,
						instancetype.DefaultPreferenceEnv:   defaultPreferenceX86_64,
					},
				},
				&fedora{
					Version: "39",
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
