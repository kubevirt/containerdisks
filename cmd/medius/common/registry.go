package common

import (
	"kubevirt.io/containerdisks/artifacts/centos"
	"kubevirt.io/containerdisks/artifacts/fedora"
	"kubevirt.io/containerdisks/artifacts/rhcos"
	"kubevirt.io/containerdisks/pkg/api"
)

type Entry struct {
	Artifact   api.Artifact
	UseForDocs bool
}

var Registry = []Entry{
	{
		Artifact:   fedora.New("35"),
		UseForDocs: true,
	},
	{
		Artifact:   rhcos.New("4.9"),
		UseForDocs: true,
	},
	{
		Artifact:   centos.New("8.4"),
		UseForDocs: true,
	},
	{
		Artifact:   centos.New("7-2009"),
		UseForDocs: false,
	},
}
