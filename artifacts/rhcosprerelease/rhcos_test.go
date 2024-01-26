package rhcosprerelease

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"kubevirt.io/api/instancetype"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/testutil"
)

var _ = Describe("RhcosPrerelease", func() {
	DescribeTable("Inspect should be able to parse checksum files",
		func(release, mockFile string, details *api.ArtifactDetails, additionalLabels map[string]string, metadata *api.Metadata) {
			c := New(release, additionalLabels)
			c.getter = testutil.NewMockGetter(mockFile)
			got, err := c.Inspect()
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(details))
			Expect(c.Metadata()).To(Equal(metadata))
		},
		Entry("rhcos:4.9", "latest-4.9", "testdata/rhcos-latest-4.9-prerelease.checksum",
			&api.ArtifactDetails{
				SHA256Sum:            "3466690807fb710102559ea57daac0484c59ed4d914996882d601b8bb7a7ada8",
				DownloadURL:          "https://mirror.openshift.com/pub/openshift-v4/x86_64/dependencies/rhcos/pre-release/latest-4.9/rhcos-openstack.x86_64.qcow2.gz", //nolint:lll
				Compression:          "gzip",
				AdditionalUniqueTags: []string{"4.9.0-rc.7", "3466690807fb710102559ea57daac0484c59ed4d914996882d601b8bb7a7ada8"},
			},
			map[string]string{
				instancetype.DefaultInstancetypeLabel: "u1.small",
				instancetype.DefaultPreferenceLabel:   "rhel.8",
			},
			&api.Metadata{
				Name:        "rhcos",
				Version:     "4.9-pre-release",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "core",
				},
				AdditionalLabels: map[string]string{
					instancetype.DefaultInstancetypeLabel: "u1.small",
					instancetype.DefaultPreferenceLabel:   "rhel.8",
				},
			},
		),
		Entry("rhcos:latest", "latest", "testdata/rhcos-latest-prerelease.checksum",
			&api.ArtifactDetails{
				SHA256Sum:            "f581896eee37216021bfce9ddd5e1fd8289c366ca0d1db25221c77688de85fd7",
				DownloadURL:          "https://mirror.openshift.com/pub/openshift-v4/x86_64/dependencies/rhcos/pre-release/latest/rhcos-openstack.x86_64.qcow2.gz", //nolint:lll
				Compression:          "gzip",
				AdditionalUniqueTags: []string{"4.10.0-rc.1", "f581896eee37216021bfce9ddd5e1fd8289c366ca0d1db25221c77688de85fd7"},
			},
			map[string]string{
				instancetype.DefaultInstancetypeLabel: "u1.small",
				instancetype.DefaultPreferenceLabel:   "rhel.8",
			},
			&api.Metadata{
				Name:        "rhcos",
				Version:     "latest-pre-release",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "core",
				},
				AdditionalLabels: map[string]string{
					instancetype.DefaultInstancetypeLabel: "u1.small",
					instancetype.DefaultPreferenceLabel:   "rhel.8",
				},
			},
		),
	)
})

func TestRhcosPrerelease(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RhcosPrerelease Suite")
}
