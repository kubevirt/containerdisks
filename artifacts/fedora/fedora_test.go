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
					common.DefaultInstancetypeEnv: defaultInstancetypeX86_64,
					common.DefaultPreferenceEnv:   defaultPreferenceX86_64,
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
		Entry("fedora:38 x86_64", "38", "x86_64", "testdata/releases.json",
			&api.ArtifactDetails{
				SHA256Sum:            "d334670401ff3d5b4129fcc662cf64f5a6e568228af59076cc449a4945318482",
				DownloadURL:          "https://download.fedoraproject.org/pub/fedora/linux/releases/38/Cloud/x86_64/images/Fedora-Cloud-Base-38-1.6.x86_64.qcow2", //nolint:lll
				AdditionalUniqueTags: []string{"38-1.6"},
				ImageArchitecture:    "amd64",
			},
			&api.Metadata{
				Name:        "fedora",
				Version:     "38",
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
		Entry("fedora:38 aarch64", "38", "aarch64", "testdata/releases.json",
			&api.ArtifactDetails{
				SHA256Sum:            "ad71d22104a16e4f9efa93e61e8c7bce28de693f59c802586abbe85e9db55a65",
				DownloadURL:          "https://download.fedoraproject.org/pub/fedora/linux/releases/38/Cloud/aarch64/images/Fedora-Cloud-Base-38-1.6.aarch64.qcow2", //nolint:lll
				AdditionalUniqueTags: []string{"38-1.6"},
				ImageArchitecture:    "arm64",
			},
			&api.Metadata{
				Name:        "fedora",
				Version:     "38",
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
					Version: "39",
					Arch:    "x86_64",
					Variant: "Cloud",
					getter:  &http.HTTPGetter{},
					EnvVariables: map[string]string{
						common.DefaultInstancetypeEnv: defaultInstancetypeX86_64,
						common.DefaultPreferenceEnv:   defaultPreferenceX86_64,
					},
				},
				&fedora{
					Version: "39",
					Arch:    "aarch64",
					Variant: "Cloud",
					getter:  &http.HTTPGetter{},
				},
			},
			{
				&fedora{
					Version: "38",
					Arch:    "x86_64",
					Variant: "Cloud",
					getter:  &http.HTTPGetter{},
					EnvVariables: map[string]string{
						common.DefaultInstancetypeEnv: defaultInstancetypeX86_64,
						common.DefaultPreferenceEnv:   defaultPreferenceX86_64,
					},
				},
				&fedora{
					Version: "38",
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
