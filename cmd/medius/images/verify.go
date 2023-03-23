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
	"k8s.io/utils/pointer"
	v1 "kubevirt.io/api/core/v1"
	kvirtcli "kubevirt.io/client-go/kubecli"
	kvirtlog "kubevirt.io/client-go/log"
	"kubevirt.io/containerdisks/cmd/medius/common"
	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/docs"
)

const (
	VerifyUsername = "verify"
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

			resultsChan, err := spawnWorkers(cmd.Context(), options, func(e *common.Entry) (*api.ArtifactResult, error) {
				description := e.Artifact.Metadata().Describe()
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
				err := verifyArtifact(cmd.Context(), e.Artifact, r, options, client)
				if err != nil {
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

			if err := writeResultsFile(options.ImagesOptions.ResultsFile, results); err != nil {
				logrus.Fatal(err)
			}

			if err != nil {
				if options.VerifyImagesOptions.NoFail {
					logrus.Warn(err)
				} else {
					logrus.Fatal(err)
				}
			}
		},
	}
	verifyCmd.Flags().StringVar(&options.VerifyImagesOptions.Registry, "registry", options.VerifyImagesOptions.Registry, "Registry that contains containerdisks to verify")
	verifyCmd.Flags().StringVar(&options.VerifyImagesOptions.Namespace, "namespace", options.VerifyImagesOptions.Namespace, "Namespace to run verify in")
	verifyCmd.Flags().BoolVar(&options.VerifyImagesOptions.NoFail, "no-fail", options.VerifyImagesOptions.NoFail, "Return success even if a worker fails")
	verifyCmd.Flags().IntVar(&options.VerifyImagesOptions.Timeout, "timeout", options.VerifyImagesOptions.Timeout, "Maximum seconds to wait for VM to be running")
	verifyCmd.Flags().AddGoFlagSet(kvirtcli.FlagSet())

	err := verifyCmd.MarkFlagRequired("registry")
	if err != nil {
		logrus.Fatal(err)
	}

	return verifyCmd
}

func verifyArtifact(ctx context.Context, artifact api.Artifact, result api.ArtifactResult, options *common.Options, client kvirtcli.KubevirtClient) error {
	log := common.Logger(artifact)

	if len(result.Tags) == 0 {
		err := errors.New("no containerdisks to verify")
		log.Error(err)
		return err
	}

	imgRef := path.Join(options.VerifyImagesOptions.Registry, result.Tags[0])
	vm, privateKey, err := createVM(artifact, imgRef)
	if err != nil {
		log.WithError(err).Error("Failed to create VM object")
		return err
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return ctx.Err()
	}

	vmClient := client.VirtualMachine(options.VerifyImagesOptions.Namespace)
	log.Info("Creating VM")
	if vm, err = vmClient.Create(vm); err != nil {
		log.WithError(err).Error("Failed to create VM")
		return err
	}

	defer func() {
		if err = vmClient.Delete(vm.Name, &metav1.DeleteOptions{GracePeriodSeconds: pointer.Int64(0)}); err != nil {
			log.WithError(err).Error("Failed to delete VM")
		}
	}()

	if errors.Is(ctx.Err(), context.Canceled) {
		return ctx.Err()
	}

	log.Info("Waiting for VM to be ready")
	if err = waitVMReady(ctx, vm.Name, vmClient, options.VerifyImagesOptions.Timeout); err != nil {
		if errors.Is(ctx.Err(), context.Canceled) {
			return ctx.Err()
		}

		log.WithError(err).Error("VM not ready")
		return err
	}

	vmi, err := client.VirtualMachineInstance(options.VerifyImagesOptions.Namespace).Get(vm.Name, &metav1.GetOptions{})
	if err != nil {
		log.WithError(err).Error("Failed to get VMI")
		return err
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return ctx.Err()
	}

	log.Info("Running tests on VMI")
	for _, testFn := range artifact.Tests() {
		if err = testFn(ctx, vmi, &api.ArtifactTestParams{Username: VerifyUsername, PrivateKey: privateKey}); err != nil {
			log.WithError(err).Error("Failed to verify containerdisk")
			return err
		}
		if errors.Is(ctx.Err(), context.Canceled) {
			return ctx.Err()
		}
	}

	log.Info("Tests successful")
	return nil
}

func createVM(artifact api.Artifact, imgRef string) (*v1.VirtualMachine, ed25519.PrivateKey, error) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	publicKey, err := marshallPublicKey(&privateKey)
	if err != nil {
		return nil, nil, err
	}

	userData := artifact.UserData(
		&docs.UserData{
			Username:       VerifyUsername,
			AuthorizedKeys: []string{publicKey},
		},
	)

	name := randName(artifact.Metadata().Name)
	vm := artifact.VM(name, imgRef, userData)
	vm.Spec.Template.Spec.TerminationGracePeriodSeconds = pointer.Int64(0)
	return vm, privateKey, nil
}

func marshallPublicKey(key *ed25519.PrivateKey) (string, error) {
	sshKey, err := ssh.NewPublicKey(key.Public())
	if err != nil {
		return "", err
	}

	marshalled := string(ssh.MarshalAuthorizedKey(sshKey))
	return marshalled[:len(marshalled)-1], nil
}

func randName(name string) string {
	return name + "-" + urand.String(5)
}

func waitVMReady(ctx context.Context, name string, client kvirtcli.VirtualMachineInterface, timeout int) error {
	return wait.PollImmediateWithContext(ctx, time.Second, time.Duration(timeout)*time.Second, func(_ context.Context) (bool, error) {
		vm, err := client.Get(name, &metav1.GetOptions{})

		if err != nil {
			return false, err
		}

		return vm.Status.Ready, nil
	})
}
