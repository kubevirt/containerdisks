package images

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
	"path"
	"slices"
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
	"kubevirt.io/containerdisks/pkg/architecture"
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

			// Get client node architecture
			nodeArchitecture, err := retrieveNodeArch(client)
			if err != nil {
				logrus.Fatal(err)
			}
			logrus.Infof("Node Architecture: %s\n", nodeArchitecture)

			focusMatched, resultsChan, workerErr := spawnWorkers(cmd.Context(), options, func(e *common.Entry) (*api.ArtifactResult, error) {
				artifact, err := retrieveArchitectureArtifact(nodeArchitecture, e)
				if err != nil {
					return nil, err
				}
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
				err = verifyArtifact(cmd.Context(), artifact, r, options, client)
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

func retrieveNodeArch(client kvirtcli.KubevirtClient) (string, error) {
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	if len(nodes.Items) == 0 {
		return "", fmt.Errorf("no nodes found")
	}
	return nodes.Items[0].Status.NodeInfo.Architecture, nil
}

func retrieveArchitectureArtifact(nodeArchitecture string, e *common.Entry) (api.Artifact, error) {
	archIndex := slices.IndexFunc(e.Artifacts, func(a api.Artifact) bool {
		return architecture.GetImageArchitecture(a.Metadata().Arch) == nodeArchitecture
	})
	if archIndex == -1 {
		return nil, fmt.Errorf("no artifact found for node architecture %s", nodeArchitecture)
	}
	return e.Artifacts[archIndex], nil
}

func verifyArtifact(ctx context.Context, a api.Artifact, res api.ArtifactResult, o *common.Options, client kvirtcli.KubevirtClient) error {
	log := common.Logger(a)

	if len(res.Tags) == 0 {
		err := errors.New("no containerdisks to verify")
		log.Error(err)
		return err
	}

	imgRef := path.Join(o.VerifyImagesOptions.Registry, res.Tags[0])
	vm, username, privateKey, err := createVM(a, imgRef)
	if err != nil {
		log.WithError(err).Error("Failed to create VM object")
		return err
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return ctx.Err()
	}

	vmClient := client.VirtualMachine(o.VerifyImagesOptions.Namespace)
	log.Info("Creating VM")
	if vm, err = vmClient.Create(ctx, vm, metav1.CreateOptions{}); err != nil {
		log.WithError(err).Error("Failed to create VM")
		return err
	}

	defer func() {
		if err = vmClient.Delete(ctx, vm.Name, metav1.DeleteOptions{GracePeriodSeconds: ptr.To[int64](0)}); err != nil {
			log.WithError(err).Error("Failed to delete VM")
		}
	}()

	if errors.Is(ctx.Err(), context.Canceled) {
		return ctx.Err()
	}

	log.Info("Waiting for VM to be ready")
	if err = waitVMReady(ctx, vm.Name, vmClient, o.VerifyImagesOptions.Timeout); err != nil {
		if errors.Is(ctx.Err(), context.Canceled) {
			return ctx.Err()
		}

		log.WithError(err).Error("VM not ready")
		return err
	}

	vmi, err := client.VirtualMachineInstance(o.VerifyImagesOptions.Namespace).Get(ctx, vm.Name, metav1.GetOptions{})
	if err != nil {
		log.WithError(err).Error("Failed to get VMI")
		return err
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return ctx.Err()
	}

	log.Info("Running tests on VMI")
	for _, testFn := range a.Tests() {
		if err = testFn(ctx, vmi, &api.ArtifactTestParams{Username: username, PrivateKey: privateKey}); err != nil {
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

func createVM(artifact api.Artifact, imgRef string) (*v1.VirtualMachine, string, ed25519.PrivateKey, error) {
	metadata := artifact.Metadata()
	username := metadata.ExampleUserData.Username

	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, "", nil, err
	}

	publicKey, err := marshallPublicKey(&privateKey)
	if err != nil {
		return nil, "", nil, err
	}

	userData := artifact.UserData(
		&docs.UserData{
			Username:       username,
			AuthorizedKeys: []string{publicKey},
		},
	)

	name := randName(metadata.Name)
	vm := artifact.VM(name, imgRef, userData)
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
