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
			Expect(got.ChecksumHash).ToNot(BeNil())
			Expect(got.Checksum).To(Equal(details.Checksum))
			Expect(got.DownloadURL).To(Equal(details.DownloadURL))
			Expect(got.AdditionalUniqueTags).To(Equal(details.AdditionalUniqueTags))
			Expect(got.ImageArchitecture).To(Equal(details.ImageArchitecture))
			Expect(got.Compression).To(Equal(details.Compression))
			Expect(c.Metadata()).To(Equal(metadata))
		},
		Entry("fedora:40 x86_64", "40", "x86_64", "testdata/releases.json",
			&api.ArtifactDetails{
				Checksum:             "ac58f3c35b73272d5986fa6d3bc44fd246b45df4c334e99a07b3bbd00684adee",
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
					common.DefaultInstancetypeEnv: defaultInstancetype,
					common.DefaultPreferenceEnv:   defaultPreferenceX86_64,
				},
				Arch: "x86_64",
			},
		),
		Entry("fedora:40 aarch64", "40", "aarch64", "testdata/releases.json",
			&api.ArtifactDetails{
				Checksum:             "ebdce26d861a9d15072affe1919ed753ec7015bd97b3a7d0d0df6a10834f7459",
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
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: defaultInstancetype,
					common.DefaultPreferenceEnv:   defaultPreferenceAarch64,
				},
				Arch: "aarch64",
			},
		),
		Entry("fedora:40 s390x", "40", "s390x", "testdata/releases.json",
			&api.ArtifactDetails{
				Checksum:             "808226b31c6c61e08cde77fe7ba61d766f7528c857e7ae8553040c177cbda9a7",
				DownloadURL:          "https://download.fedoraproject.org/pub/fedora-secondary/releases/40/Cloud/s390x/images/Fedora-Cloud-Base-Generic.s390x-40-1.14.qcow2", //nolint:lll
				AdditionalUniqueTags: []string{"40-1.14"},
				ImageArchitecture:    "s390x",
			},
			&api.Metadata{
				Name:        "fedora",
				Version:     "40",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "fedora",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: defaultInstancetype,
					common.DefaultPreferenceEnv:   defaultPreferenceS390x,
				},
				Arch: "s390x",
			},
		),
		Entry("fedora:39 x86_64", "39", "x86_64", "testdata/releases.json",
			&api.ArtifactDetails{
				Checksum:             "ab5be5058c5c839528a7d6373934e0ce5ad6c8f80bd71ed3390032027da52f37",
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
					common.DefaultInstancetypeEnv: defaultInstancetype,
					common.DefaultPreferenceEnv:   defaultPreferenceX86_64,
				},
				Arch: "x86_64",
			},
		),
		Entry("fedora:39 aarch64", "39", "aarch64", "testdata/releases.json",
			&api.ArtifactDetails{
				Checksum:             "765996d5b77481ca02d0ac06405641bf134ac920cfc1e60d981c64d7971162dc",
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
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: defaultInstancetype,
					common.DefaultPreferenceEnv:   defaultPreferenceAarch64,
				},
				Arch: "aarch64",
			},
		),
		Entry("fedora:39 s390x", "39", "s390x", "testdata/releases.json",
			&api.ArtifactDetails{
				Checksum:             "36dec66c791c9d1225d74e8828fdb0976ad89f695e8e6f5c93269cafa8563907",
				DownloadURL:          "https://download.fedoraproject.org/pub/fedora-secondary/releases/39/Cloud/s390x/images/Fedora-Cloud-Base-39-1.5.s390x.qcow2", //nolint:lll
				AdditionalUniqueTags: []string{"39-1.5"},
				ImageArchitecture:    "s390x",
			},
			&api.Metadata{
				Name:        "fedora",
				Version:     "39",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "fedora",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: defaultInstancetype,
					common.DefaultPreferenceEnv:   defaultPreferenceS390x,
				},
				Arch: "s390x",
			},
		),
		Entry("fedora:41-beta aarch64", "41 Beta", "aarch64", "testdata/releases.json",
			&api.ArtifactDetails{
				Checksum:             "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				DownloadURL:          "https://download.fedoraproject.org/pub/fedora/linux/releases/test/41_Beta/Cloud/aarch64/images/Fedora-Cloud-Base-Generic-41_Beta-1.2.aarch64.qcow2", //nolint:lll
				AdditionalUniqueTags: []string(nil),
				ImageArchitecture:    "arm64",
			},
			&api.Metadata{
				Name:        "fedora",
				Version:     "41-beta",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "fedora",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: defaultInstancetype,
					common.DefaultPreferenceEnv:   defaultPreferenceAarch64,
				},
				Arch:         "aarch64",
				IsPrerelease: true,
			},
		),
	)

	It("Gather should be able to parse releases files", func() {
		// Stable versions should come first (sorted descending), followed by
		// prerelease versions. This ensures stable releases get the "latest" tag.
		artifacts := [][]api.Artifact{
			{
				parsedRelease("40", "40", "x86_64", defaultPreferenceX86_64),
				parsedRelease("40", "40", "aarch64", defaultPreferenceAarch64),
			},
			{
				parsedRelease("39", "39", "x86_64", defaultPreferenceX86_64),
				parsedRelease("39", "39", "aarch64", defaultPreferenceAarch64),
				parsedRelease("39", "39", "s390x", defaultPreferenceS390x),
			},
			{
				parsedRelease("41-beta", "41 Beta", "x86_64", defaultPreferenceX86_64),
				parsedRelease("41-beta", "41 Beta", "aarch64", defaultPreferenceAarch64),
			},
		}

		c := NewGatherer()
		c.getter = testutil.NewMockGetter("testdata/releases.json")
		got, err := c.Gather()
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal(artifacts))
	})

	DescribeTable("IsStableVersion",
		func(version string, expected bool) {
			Expect(IsStableVersion(version)).To(Equal(expected))
		},
		Entry("stable", "43", true),
		Entry("beta", "44 Beta", false),
		Entry("rc", "44 RC1", false),
	)

	DescribeTable("NormalizeVersion",
		func(version, expected string) {
			Expect(NormalizeVersion(version)).To(Equal(expected))
		},
		Entry("stable", "43", "43"),
		Entry("beta", "44 Beta", "44-beta"),
	)
})

func parsedRelease(version, releaseVersion, arch, defaultPreference string) api.Artifact {
	return &fedora{
		Version:        version,
		ReleaseVersion: releaseVersion,
		Arch:           arch,
		Variant:        "Cloud",
		getter:         &http.HTTPGetter{},
		EnvVariables: map[string]string{
			common.DefaultInstancetypeEnv: defaultInstancetype,
			common.DefaultPreferenceEnv:   defaultPreference,
		},
	}
}

func TestFedora(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fedora Suite")
}
