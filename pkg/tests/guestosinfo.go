package tests

import (
	"context"

	v1 "kubevirt.io/api/core/v1"
	kvirtcli "kubevirt.io/client-go/kubecli"
	"kubevirt.io/containerdisks/pkg/api"
)

func GuestOsInfo(ctx context.Context, vmi *v1.VirtualMachineInstance, _ *api.ArtifactTestParams) error {
	client, err := kvirtcli.GetKubevirtClient()
	if err != nil {
		return err
	}

	return retryTest(ctx, func() error {
		_, err := client.VirtualMachineInstance(vmi.Namespace).GuestOsInfo(vmi.Name)
		return err
	})
}
