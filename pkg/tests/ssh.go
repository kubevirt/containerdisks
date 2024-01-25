package tests

import (
	"context"
	"fmt"

	"golang.org/x/crypto/ssh"
	v1 "kubevirt.io/api/core/v1"
	kvirtcli "kubevirt.io/client-go/kubecli"

	"kubevirt.io/containerdisks/pkg/api"
)

func SSH(ctx context.Context, vmi *v1.VirtualMachineInstance, params *api.ArtifactTestParams) error {
	kvirtClient, err := kvirtcli.GetKubevirtClient()
	if err != nil {
		return err
	}

	signer, err := ssh.NewSignerFromKey(params.PrivateKey)
	if err != nil {
		return err
	}

	// Test SSH while deliberately ignoring insecure host keys
	config := &ssh.ClientConfig{
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec
		User:            params.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	return retryTest(ctx, func() error {
		return testSSH(vmi, kvirtClient, config)
	})
}

func testSSH(vmi *v1.VirtualMachineInstance, kvirtClient kvirtcli.KubevirtClient, config *ssh.ClientConfig) error {
	const sshPort = 22
	tunnel, err := kvirtClient.VirtualMachineInstance(vmi.Namespace).PortForward(vmi.Name, sshPort, "tcp")
	if err != nil {
		return fmt.Errorf("failed to forward ssh port: %w", err)
	}

	conn := tunnel.AsConn()
	addr := fmt.Sprintf("vmi/%s.%s:22", vmi.Name, vmi.Namespace)
	sshConn, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
	if err != nil {
		return err
	}

	session, err := ssh.NewClient(sshConn, chans, reqs).NewSession()
	if err != nil {
		return err
	}

	err = session.Run("echo hello")
	if err != nil {
		return err
	}

	return nil
}
