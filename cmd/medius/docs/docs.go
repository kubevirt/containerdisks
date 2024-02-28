package docs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"kubevirt.io/containerdisks/cmd/medius/common"
	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/pkg/quay"
)

func NewPublishDocsCommand(options *common.Options) *cobra.Command {
	options.PublishDocsOptions = common.PublishDocsOptions{
		Registry: "quay.io/containerdisks",
	}

	publishCmd := &cobra.Command{
		Use:   "publish",
		Short: "Synchronize container disk descriptions with quay.io",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(options)
		},
	}
	publishCmd.Flags().StringVar(&options.PublishDocsOptions.Registry, "registry",
		options.PublishDocsOptions.Registry, "target registry for the containerdisks")
	publishCmd.Flags().StringVar(&options.PublishDocsOptions.TokenFile, "quay-token-file",
		options.PublishDocsOptions.TokenFile, "quay.io oauth token file")

	err := publishCmd.MarkFlagRequired("quay-token-file")
	if err != nil {
		logrus.Fatal(err)
	}

	return publishCmd
}

func run(options *common.Options) error {
	success := true
	focusMatched := false

	quayOrg, err := getQuayOrg(options.PublishDocsOptions.Registry)
	if err != nil {
		return err
	}

	client := quay.NewQuayClient(options.PublishDocsOptions.TokenFile, quayOrg)
	registry := common.NewRegistry()
	for i, p := range registry {
		if common.ShouldSkip(options.Focus, &registry[i]) || !p.UseForDocs {
			continue
		}

		focusMatched = true

		artifact := p.Artifacts[0]
		log := common.Logger(artifact)
		name := artifact.Metadata().Name

		description, err := createDescription(artifact, options.PublishDocsOptions.Registry)
		if err != nil {
			success = false
			log.Errorf("error marshaling example for %q: %v", name, err)
			continue
		}

		log.Info("Updating description on quay.io")
		if !options.DryRun {
			if err := client.Update(context.Background(), name, description); err != nil {
				success = false
				log.Errorf("error marshaling example for for %q: %v", name, err)
			}
		}
	}

	if !focusMatched {
		return fmt.Errorf("no artifact was processed, focus '%s' did not match", options.Focus)
	}

	if !success {
		return errors.New("an error occurred during publishing of the docs")
	}

	return nil
}

func getQuayOrg(registry string) (string, error) {
	elements := strings.Split(registry, "/")
	if len(elements) != 2 || elements[0] != "quay.io" || elements[1] == "" {
		return "", fmt.Errorf(
			"error determining quay.io organization from %v, this command only works with quay.io",
			registry,
		)
	}

	return elements[1], nil
}

func createDescription(artifact api.Artifact, registry string) (string, error) {
	metadata := artifact.Metadata()
	vm := artifact.VM(
		metadata.Name,
		path.Join(registry, metadata.Describe()),
		artifact.UserData(&metadata.ExampleUserData),
	)

	example, err := yaml.Marshal(&vm)
	if err != nil {
		return "", fmt.Errorf("error marshaling example for for %q: %v", metadata.Name, err)
	}

	data := &docs.TemplateData{
		Name:        metadata.Name,
		Description: metadata.Description,
		Example:     string(example),
	}

	var result bytes.Buffer
	if err := docs.Template().Execute(&result, data); err != nil {
		return "", fmt.Errorf("error rendering template for %q: %v", metadata.Name, err)
	}

	return result.String(), nil
}
