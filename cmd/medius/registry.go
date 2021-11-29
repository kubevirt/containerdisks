package main

import (
	"kubevirt.io/containerdisks/artifacts/centos"
	"kubevirt.io/containerdisks/artifacts/fedora"
	"kubevirt.io/containerdisks/artifacts/rhcos"
	"kubevirt.io/containerdisks/pkg/api"
)

var registry = []api.Artifact{
	fedora.New("35"),
	rhcos.New("4.9"),
	centos.New("8.4"),
	centos.New("7-2009"),
}
