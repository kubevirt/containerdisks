package docs

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"kubevirt.io/containerdisks/cmd/medius/common"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/pkg/quay"
	"sigs.k8s.io/yaml"
)

func NewPublishDocsCommand(options *common.Options) *cobra.Command {
	options.PublishDocsOptions = common.PublishDocsOptions{
		Registry: "quay.io/containerdisks",
	}

	publishCmd := &cobra.Command{
		Use:   "publish",
		Short: "Synchronize container disk descriptions with quay.io",
		RunE: func(cmd *cobra.Command, args []string) error {
			success := true

			elements := strings.Split(options.PublishDocsOptions.Registry, "/")
			if len(elements) != 2 || elements[0] != "quay.io" || elements[1] == "" {
				return fmt.Errorf(
					"error determining quay.io organization from %v, this command only works with quay.io",
					options.PublishDocsOptions.Registry,
				)
			}

			client := quay.NewQuayClient(options.PublishDocsOptions.TokenFile, elements[1])
			for _, p := range common.Registry {
				if options.Focus == "" && p.SkipWhenNotFocused {
					continue
				}
				if options.Focus != "" && options.Focus != p.Artifact.Metadata().Describe() {
					continue
				}
				if !p.UseForDocs {
					continue
				}
				log := common.Logger(p.Artifact)
				metadata := p.Artifact.Metadata()
				vm := p.Artifact.VM(
					metadata.Name,
					path.Join(options.PublishDocsOptions.Registry, p.Artifact.Metadata().Describe()),
					metadata.ExampleUserDataPayload,
				)

				example, err := yaml.Marshal(&vm)
				if err != nil {
					return fmt.Errorf("error marshaling example for for %q: %v", metadata.Name, err)
				}
				var result bytes.Buffer
				data := &docs.TemplateData{
					Name:        metadata.Name,
					Description: metadata.Description,
					Example:     string(example),
				}
				if err := docs.Template().Execute(&result, data); err != nil {
					return fmt.Errorf("error rendering template for %q: %v", metadata.Name, err)
				}

				log.Info("Updating description on quay.io")
				if options.DryRun == false {
					if err := client.Update(context.Background(), metadata.Name, result.String()); err != nil {
						success = false
						log.Errorf("error marshaling example for for %q: %v", metadata.Name, err)
					}
				}
			}
			if success == false {
				os.Exit(1)
			}
			return nil
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
