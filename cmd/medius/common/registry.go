package common

import (
	"github.com/sirupsen/logrus"
	"kubevirt.io/containerdisks/artifacts/centos"
	"kubevirt.io/containerdisks/artifacts/centosstream"
	"kubevirt.io/containerdisks/artifacts/fedora"
	"kubevirt.io/containerdisks/artifacts/generic"
	"kubevirt.io/containerdisks/artifacts/rhcos"
	"kubevirt.io/containerdisks/artifacts/rhcosprerelease"
	"kubevirt.io/containerdisks/artifacts/ubuntu"
	"kubevirt.io/containerdisks/pkg/api"
)

type Entry struct {
	Artifact           api.Artifact
	UseForDocs         bool
	UseForLatest       bool
	SkipWhenNotFocused bool
}

var Registry = []Entry{
	{
		Artifact:   rhcos.New("4.9", true),
		UseForDocs: false,
	},
	{
		Artifact:   rhcos.New("4.10", true),
		UseForDocs: false,
	},
	{
		Artifact:   rhcos.New("4.11", true),
		UseForDocs: false,
	},
	{
		Artifact:   rhcos.New("4.12", true),
		UseForDocs: true,
	},
	{
		Artifact:   rhcos.New("latest", false),
		UseForDocs: false,
	},
	{
		Artifact:   rhcosprerelease.New("latest-4.9"),
		UseForDocs: false,
	},
	{
		Artifact:   rhcosprerelease.New("latest-4.10"),
		UseForDocs: false,
	},
	{
		Artifact:   rhcosprerelease.New("latest-4.11"),
		UseForDocs: false,
	},
	{
		Artifact:   rhcosprerelease.New("latest-4.12"),
		UseForDocs: false,
	},
	{
		Artifact:   rhcosprerelease.New("latest-4.13"),
		UseForDocs: false,
	},
	{
		Artifact:   rhcosprerelease.New("latest"),
		UseForDocs: false,
	},
	{
		Artifact:   centos.New("8.4"),
		UseForDocs: false,
	},
	{
		Artifact:   centos.New("7-2009"),
		UseForDocs: true,
	},
	{
		Artifact:   centosstream.New("9"),
		UseForDocs: true,
	},
	{
		Artifact:   centosstream.New("8"),
		UseForDocs: false,
	},
	{
		Artifact:   ubuntu.New("22.04"),
		UseForDocs: true,
	},
	{
		Artifact:   ubuntu.New("20.04"),
		UseForDocs: false,
	},
	{
		Artifact:   ubuntu.New("18.04"),
		UseForDocs: false,
	},
	// for testing only
	{
		Artifact: generic.New(
			&api.ArtifactDetails{
				SHA256Sum:   "cc704ab14342c1c8a8d91b66a7fc611d921c8b8f1aaf4695f9d6463d913fa8d1",
				DownloadURL: "https://download.cirros-cloud.net/0.6.1/cirros-0.6.1-x86_64-disk.img",
			},
			&api.Metadata{
				Name:    "cirros",
				Version: "6.1",
			},
		),
		SkipWhenNotFocused: true,
		UseForDocs:         false,
	},
}

func init() {
	for _, gatherer := range []api.ArtifactsGatherer{fedora.NewGatherer()} {
		artifacts, err := gatherer.Gather()
		if err != nil {
			logrus.Warn("Failed to gather artifacts", err)
		} else {
			for i := range artifacts {
				Registry = append(Registry, Entry{
					Artifact:     artifacts[i],
					UseForDocs:   i == 0,
					UseForLatest: i == 0,
				})
			}
		}
	}
}
