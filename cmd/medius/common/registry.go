package common

import (
	"crypto/sha256"
	"strings"

	"github.com/sirupsen/logrus"

	"kubevirt.io/containerdisks/artifacts/centosstream"
	"kubevirt.io/containerdisks/artifacts/debian"
	"kubevirt.io/containerdisks/artifacts/fedora"
	"kubevirt.io/containerdisks/artifacts/generic"
	"kubevirt.io/containerdisks/artifacts/opensuse/leap"
	"kubevirt.io/containerdisks/artifacts/opensuse/tumbleweed"
	"kubevirt.io/containerdisks/artifacts/ubuntu"
	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/common"
	"kubevirt.io/containerdisks/pkg/docs"
)

type Entry struct {
	Artifacts          []api.Artifact
	UseForDocs         bool
	UseForLatest       bool
	SkipWhenNotFocused bool
}

var staticRegistry = []Entry{
	{
		Artifacts: []api.Artifact{
			centosstream.New("10", "x86_64", &docs.UserData{Username: "cloud-user"}, defaultEnvVariables("u1.medium", "centos.stream10")),
			centosstream.New("10", "aarch64", &docs.UserData{Username: "cloud-user"}, defaultEnvVariables("u1.medium", "centos.stream10")),
			centosstream.New("10", "s390x", &docs.UserData{Username: "cloud-user"}, defaultEnvVariables("u1.medium", "centos.stream10")),
		},
		UseForDocs: true,
	},
	{
		Artifacts: []api.Artifact{
			centosstream.New("9", "x86_64", &docs.UserData{Username: "cloud-user"}, defaultEnvVariables("u1.medium", "centos.stream9")),
			centosstream.New("9", "aarch64", &docs.UserData{Username: "cloud-user"}, defaultEnvVariables("u1.medium", "centos.stream9")),
			centosstream.New("9", "s390x", &docs.UserData{Username: "cloud-user"}, defaultEnvVariables("u1.medium", "centos.stream9")),
		},
		UseForDocs: true,
	},
	{
		Artifacts: []api.Artifact{
			ubuntu.New("24.04", "x86_64", defaultEnvVariables("u1.medium", "ubuntu")),
			ubuntu.New("24.04", "aarch64", defaultEnvVariables("u1.medium", "ubuntu")),
			ubuntu.New("24.04", "s390x", defaultEnvVariables("u1.medium", "ubuntu")),
		},
		UseForDocs: true,
	},
	{
		Artifacts: []api.Artifact{
			ubuntu.New("22.04", "x86_64", defaultEnvVariables("u1.medium", "ubuntu")),
			ubuntu.New("22.04", "aarch64", defaultEnvVariables("u1.medium", "ubuntu")),
			ubuntu.New("22.04", "s390x", defaultEnvVariables("u1.medium", "ubuntu")),
		},
		UseForDocs: false,
	},
	{
		Artifacts: []api.Artifact{
			ubuntu.New("20.04", "x86_64", defaultEnvVariables("u1.medium", "ubuntu")),
			ubuntu.New("20.04", "aarch64", defaultEnvVariables("u1.medium", "ubuntu")),
			ubuntu.New("20.04", "s390x", defaultEnvVariables("u1.medium", "ubuntu")),
		},
		UseForDocs: false,
	},
	{
		Artifacts: []api.Artifact{
			tumbleweed.New("x86_64", defaultEnvVariables("u1.medium", "opensuse.tumbleweed")),
		},
		UseForDocs: true,
	},
	{
		Artifacts: []api.Artifact{
			leap.New("x86_64", "15.6", defaultEnvVariables("u1.medium", "opensuse.leap")),
			leap.New("aarch64", "15.6", defaultEnvVariables("u1.medium", "opensuse.leap")),
		},
		UseForDocs: true,
	},
	{
		Artifacts: []api.Artifact{
			leap.New("x86_64", "15.5", defaultEnvVariables("u1.medium", "opensuse.leap")),
			leap.New("aarch64", "15.5", defaultEnvVariables("u1.medium", "opensuse.leap")),
		},
	},
	{
		Artifacts: []api.Artifact{
			debian.New("11", "bullseye", "x86_64", &docs.UserData{Username: "debian"}, nil),
			debian.New("11", "bullseye", "aarch64", &docs.UserData{Username: "debian"}, nil),
		},
	},
	{
		Artifacts: []api.Artifact{
			debian.New("12", "bookworm", "x86_64", &docs.UserData{Username: "debian"}, nil),
			debian.New("12", "bookworm", "aarch64", &docs.UserData{Username: "debian"}, nil),
		},
		UseForDocs:   true,
		UseForLatest: true,
	},
	// for testing only
	{
		Artifacts: []api.Artifact{
			generic.New(
				&api.ArtifactDetails{
					Checksum:          "cc704ab14342c1c8a8d91b66a7fc611d921c8b8f1aaf4695f9d6463d913fa8d1",
					ChecksumHash:      sha256.New,
					DownloadURL:       "https://download.cirros-cloud.net/0.6.1/cirros-0.6.1-x86_64-disk.img",
					ImageArchitecture: "amd64",
				},
				&api.Metadata{
					Name:    "cirros",
					Version: "6.1",
				},
			),
			generic.New(
				&api.ArtifactDetails{
					Checksum:          "db9420c481c11dee17860aa46fb1a3efa05fa4fb152726d6344e24da03cb0ccf",
					ChecksumHash:      sha256.New,
					DownloadURL:       "https://download.cirros-cloud.net/0.6.1/cirros-0.6.1-aarch64-disk.img",
					ImageArchitecture: "arm64",
				},
				&api.Metadata{
					Name:    "cirros",
					Version: "6.1",
				},
			),
		},
		SkipWhenNotFocused: true,
		UseForDocs:         false,
	},
}

func gatherArtifacts(registry *[]Entry, gatherers []api.ArtifactsGatherer) {
	for _, gatherer := range gatherers {
		artifacts, err := gatherer.Gather()
		if err != nil {
			logrus.Warn("Failed to gather artifacts", err)
		} else {
			for i := range artifacts {
				*registry = append(*registry, Entry{
					Artifacts:    artifacts[i],
					UseForDocs:   i == 0,
					UseForLatest: i == 0,
				})
			}
		}
	}
}

func defaultEnvVariables(defaultInstancetype, defaultPreference string) map[string]string {
	return map[string]string{
		common.DefaultInstancetypeEnv: defaultInstancetype,
		common.DefaultPreferenceEnv:   defaultPreference,
	}
}

func NewRegistry() []Entry {
	registry := make([]Entry, len(staticRegistry))
	copy(registry, staticRegistry)

	gatherers := []api.ArtifactsGatherer{fedora.NewGatherer()}
	gatherArtifacts(&registry, gatherers)

	return registry
}

func ShouldSkip(focus string, entry *Entry) bool {
	if focus == "" {
		return entry.SkipWhenNotFocused
	}

	if len(entry.Artifacts) == 0 {
		return true
	}

	focusSplit := strings.Split(focus, ":")
	wildcardFocus := len(focusSplit) == 2 && focusSplit[1] == "*"

	if wildcardFocus {
		return focusSplit[0] != entry.Artifacts[0].Metadata().Name
	}

	return focus != entry.Artifacts[0].Metadata().Describe()
}
