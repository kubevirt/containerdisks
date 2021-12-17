package docs

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"kubevirt.io/containerdisks/pkg/quay"

	"github.com/ghodss/yaml"
	"kubevirt.io/containerdisks/cmd/medius/common"
	"kubevirt.io/containerdisks/pkg/docs"
)

func NewPublishDocsCommand(options *common.Options) *cobra.Command {
	options.PublishDocsOptions = common.PublishDocsOptions{
		TokenFile: "",
	}

	publishCmd := &cobra.Command{
		Use:   "publish",
		Short: "Synchronize container disk descriptions with quay.io",
		RunE: func(cmd *cobra.Command, args []string) error {
			success := true

			elements := strings.Split(options.Registry, "/")
			if len(elements) != 2 || elements[0] != "quay.io" || elements[1] == "" {
				return fmt.Errorf("error determining quay.io organiation from %v, this command only works with quay.io", options.Registry)
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
				vm := docs.BasicVirtualMachine(metadata.Name, path.Join(options.Registry, metadata.Name)+":"+metadata.Version, metadata.ExampleCloudInitPayload)
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

	publishCmd.Flags().StringVar(&options.PublishDocsOptions.TokenFile, "quay-token-file", options.PublishDocsOptions.TokenFile, "quay.io oauth token file")
	publishCmd.MarkFlagRequired("quay-token-file")
	return publishCmd
}
