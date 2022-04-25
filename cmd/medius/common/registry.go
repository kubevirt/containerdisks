package common

import (
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
	SkipWhenNotFocused bool
}

var Registry = []Entry{
	{
		Artifact:   fedora.New("35"),
		UseForDocs: true,
	},
	{
		Artifact:   rhcos.New("4.9"),
		UseForDocs: false,
	},
	{
		Artifact:   rhcos.New("4.10"),
		UseForDocs: true,
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
		UseForDocs: true,
	},
	{
		Artifact:   rhcosprerelease.New("latest"),
		UseForDocs: false,
	},
	{
		Artifact:   centos.New("8.4"),
		UseForDocs: true,
	},
	{
		Artifact:   centos.New("7-2009"),
		UseForDocs: false,
	},
	{
		Artifact: centosstream.New("9"),
	},
	{
		Artifact:   centosstream.New("8"),
		UseForDocs: true,
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
	//for testing only
	{
		Artifact: generic.New(
			&api.ArtifactDetails{
				SHA256Sum:   "932fcae93574e242dc3d772d5235061747dfe537668443a1f0567d893614b464",
				DownloadURL: "https://download.cirros-cloud.net/0.5.2/cirros-0.5.2-x86_64-disk.img",
			},
			&api.Metadata{
				Name:    "cirros",
				Version: "5.2",
			},
		),
		SkipWhenNotFocused: true,
		UseForDocs:         false,
	},
}
