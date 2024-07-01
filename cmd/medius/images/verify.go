//nolint:lll
package images

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	urand "k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/ptr"
	v1 "kubevirt.io/api/core/v1"
	kvirtcli "kubevirt.io/client-go/kubecli"
	kvirtlog "kubevirt.io/client-go/log"

	"kubevirt.io/containerdisks/cmd/medius/common"
	"kubevirt.io/containerdisks/pkg/api"
	pkgCommon "kubevirt.io/containerdisks/pkg/common"
	"kubevirt.io/containerdisks/pkg/docs"
)

func NewVerifyImagesCommand(options *common.Options) *cobra.Command {
	options.VerifyImagesOptions = common.VerifyImageOptions{
		Namespace: "kubevirt",
		Timeout:   600,
	}

	verifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify that containerdisks are bootable and guests are working",
		Run: func(cmd *cobra.Command, args []string) {
			results, err := readResultsFile(options.ImagesOptions.ResultsFile)
			if err != nil {
				logrus.Fatal(err)
			}

			// Silence the kubevirt client log
			kvirtlog.Log = kvirtlog.MakeLogger(kvirtlog.NullLogger{})
			client, err := kvirtcli.GetKubevirtClient()
			if err != nil {
				logrus.Fatal(err)
			}

			focusMatched, resultsChan, workerErr := spawnWorkers(cmd.Context(), options, func(e *common.Entry) (*api.ArtifactResult, error) {
				artifact := e.Artifacts[0]
				description := artifact.Metadata().Describe()
				r, ok := results[description]
				if !ok {
					return nil, nil
				}
				if r.Err != "" {
					return nil, fmt.Errorf("artifact %s failed in stage %s: %s", description, r.Stage, r.Err)
				}
				if r.Stage != StagePush {
					return nil, nil
				}

				errString := ""
				v := newVerifyArtifact(artifact, r, options, client)
				if err = v.Do(cmd.Context()); err != nil {
					errString = err.Error()
				}

				return &api.ArtifactResult{
					Tags:  r.Tags,
					Stage: StageVerify,
					Err:   errString,
				}, err
			})

			for result := range resultsChan {
				results[result.Key] = result.Value
			}

			if !focusMatched {
				logrus.Fatalf("no artifact was processed, focus '%s' did not match", options.Focus)
			}

			if err := writeResultsFile(options.ImagesOptions.ResultsFile, results); err != nil {
				logrus.Fatal(err)
			}

			if workerErr != nil {
				if options.VerifyImagesOptions.NoFail {
					logrus.Warn(workerErr)
				} else {
					logrus.Fatal(workerErr)
				}
			}
		},
	}
	verifyCmd.Flags().StringVar(&options.VerifyImagesOptions.Registry, "registry",
		options.VerifyImagesOptions.Registry, "Registry that contains containerdisks to verify")
	verifyCmd.Flags().StringVar(&options.VerifyImagesOptions.Namespace, "namespace",
		options.VerifyImagesOptions.Namespace, "Namespace to run verify in")
	verifyCmd.Flags().BoolVar(&options.VerifyImagesOptions.NoFail, "no-fail",
		options.VerifyImagesOptions.NoFail, "Return success even if a worker fails")
	verifyCmd.Flags().IntVar(&options.VerifyImagesOptions.Timeout, "timeout",
		options.VerifyImagesOptions.Timeout, "Maximum seconds to wait for VM to be running")
	verifyCmd.Flags().AddGoFlagSet(kvirtcli.FlagSet())

	err := verifyCmd.MarkFlagRequired("registry")
	if err != nil {
		logrus.Fatal(err)
	}

	return verifyCmd
}

type verifyArtifact struct {
	artifact api.Artifact
	result   api.ArtifactResult
	options  *common.Options
	client   kvirtcli.KubevirtClient
	vmClient kvirtcli.VirtualMachineInterface
	log      *logrus.Entry
}

func newVerifyArtifact(artifact api.Artifact, result api.ArtifactResult, options *common.Options, client kvirtcli.KubevirtClient) *verifyArtifact {
	log := common.Logger(artifact)
	vmClient := client.VirtualMachine(options.VerifyImagesOptions.Namespace)
	return &verifyArtifact{
		artifact: artifact,
		result:   result,
		options:  options,
		client:   client,
		vmClient: vmClient,
		log:      log,
	}
}

func (v *verifyArtifact) Do(ctx context.Context) error {
	if len(v.result.Tags) == 0 {
		err := fmt.Errorf("no containerdisks to verify")
		v.log.Error(err)
		return err
	}

	v.log.Info("Verifying documentation VirtualMachine")
	if err := v.verifyVM(ctx, v.createVM); err != nil {
		return err
	}

	if v.artifact.Metadata() != nil && v.artifact.Metadata().EnvVariables != nil {
		v.log.Info("Verifying ENV variable VirtualMachine using ", v.artifact.Metadata().EnvVariables)
		return v.verifyVM(ctx, v.createVMWithEnvVariables)
	}

	return nil
}

func (v *verifyArtifact) verifyVM(ctx context.Context, createVM func() (*v1.VirtualMachine, string, ed25519.PrivateKey, error)) error {
	vm, username, privateKey, err := createVM()
	if err != nil {
		v.log.WithError(err).Error("Failed to create VM object")
		return err
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return ctx.Err()
	}
	v.log.Info("Creating VM")
	if vm, err = v.vmClient.Create(ctx, vm, metav1.CreateOptions{}); err != nil {
		v.log.WithError(err).Error("Failed to create VM")
		return err
	}

	defer func() {
		if err = v.vmClient.Delete(ctx, vm.Name, metav1.DeleteOptions{GracePeriodSeconds: ptr.To[int64](0)}); err != nil {
			v.log.WithError(err).Error("Failed to delete VM")
		}
	}()

	if errors.Is(ctx.Err(), context.Canceled) {
		return ctx.Err()
	}

	v.log.Info("Waiting for VM to be ready")
	if err = waitVMReady(ctx, vm.Name, v.vmClient, v.options.VerifyImagesOptions.Timeout); err != nil {
		if errors.Is(ctx.Err(), context.Canceled) {
			return ctx.Err()
		}

		v.log.WithError(err).Error("VM not ready")
		return err
	}

	vmi, err := v.client.VirtualMachineInstance(v.options.VerifyImagesOptions.Namespace).Get(ctx, vm.Name, metav1.GetOptions{})
	if err != nil {
		v.log.WithError(err).Error("Failed to get VMI")
		return err
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return ctx.Err()
	}

	v.log.Info("Running tests on VMI")
	for _, testFn := range v.artifact.Tests() {
		if err = testFn(ctx, vmi, &api.ArtifactTestParams{Username: username, PrivateKey: privateKey}); err != nil {
			v.log.WithError(err).Error("Failed to verify containerdisk")
			return err
		}
		if errors.Is(ctx.Err(), context.Canceled) {
			return ctx.Err()
		}
	}

	v.log.Info("Tests successful")
	return nil
}

func (v *verifyArtifact) createMetadata() (username, userData string, privateKey ed25519.PrivateKey, err error) {
	username = v.artifact.Metadata().ExampleUserData.Username

	_, privateKey, err = ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", nil, err
	}

	publicKey, err := marshallPublicKey(&privateKey)
	if err != nil {
		return "", "", nil, err
	}

	userData = v.artifact.UserData(
		&docs.UserData{
			Username:       username,
			AuthorizedKeys: []string{publicKey},
		},
	)

	return username, userData, privateKey, nil
}

func (v *verifyArtifact) createVMWithEnvVariables() (*v1.VirtualMachine, string, ed25519.PrivateKey, error) {
	username, userData, privateKey, err := v.createMetadata()
	if err != nil {
		return nil, "", nil, err
	}

	metadata := v.artifact.Metadata()
	name := randName(metadata.Name)
	imgRef := path.Join(v.options.VerifyImagesOptions.Registry, v.result.Tags[0])
	vm := v.artifact.VM(name, imgRef, userData)
	vm.Spec.Template.Spec.TerminationGracePeriodSeconds = ptr.To[int64](0)

	if instancetype, ok := metadata.EnvVariables[pkgCommon.DefaultInstancetypeEnv]; ok {
		vm.Spec.Instancetype = &v1.InstancetypeMatcher{
			Name: instancetype,
		}
		vm.Spec.Template.Spec.Domain.Resources = v1.ResourceRequirements{}
	}

	if preference, ok := metadata.EnvVariables[pkgCommon.DefaultPreferenceEnv]; ok {
		vm.Spec.Preference = &v1.PreferenceMatcher{
			Name: preference,
		}
		vm.Spec.Template.Spec.Domain.Devices.Disks[0].Disk.Bus = ""
	}

	return vm, username, privateKey, nil
}

func (v *verifyArtifact) createVM() (*v1.VirtualMachine, string, ed25519.PrivateKey, error) {
	username, userData, privateKey, err := v.createMetadata()
	if err != nil {
		return nil, "", nil, err
	}

	name := randName(v.artifact.Metadata().Name)
	imgRef := path.Join(v.options.VerifyImagesOptions.Registry, v.result.Tags[0])
	vm := v.artifact.VM(name, imgRef, userData)
	vm.Spec.Template.Spec.TerminationGracePeriodSeconds = ptr.To[int64](0)
	return vm, username, privateKey, nil
}

func marshallPublicKey(key *ed25519.PrivateKey) (string, error) {
	sshKey, err := ssh.NewPublicKey(key.Public())
	if err != nil {
		return "", err
	}

	marshaled := string(ssh.MarshalAuthorizedKey(sshKey))
	return marshaled[:len(marshaled)-1], nil
}

func randName(name string) string {
	const randomCharCount = 5
	return name + "-" + urand.String(randomCharCount)
}

func waitVMReady(ctx context.Context, name string, client kvirtcli.VirtualMachineInterface, timeout int) error {
	return wait.PollUntilContextTimeout(ctx, time.Second, time.Duration(timeout)*time.Second, true, func(_ context.Context) (bool, error) {
		vm, err := client.Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		return vm.Status.Ready, nil
	})
}
