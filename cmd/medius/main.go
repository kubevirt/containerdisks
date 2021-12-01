package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"kubevirt.io/containerdisks/cmd/medius/common"
	"kubevirt.io/containerdisks/cmd/medius/docs"
	"kubevirt.io/containerdisks/cmd/medius/images"
)

func main() {

	options := &common.Options{
		AllowInsecureRegistry: false,
		Registry:              "quay.io/containerdisks",
		DryRun:                true,
	}

	rootCmd := &cobra.Command{
		Use:   "medius",
		Short: "medius determines if new OS images are released and publishes them as containerdisks",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	imagesCmd := &cobra.Command{
		Use: "images",
		Run: func(cmd *cobra.Command, args []string) {
			os.Exit(1)
		},
	}
	docsCmd := &cobra.Command{
		Use: "docs",
		Run: func(cmd *cobra.Command, args []string) {
			os.Exit(1)
		},
	}
	rootCmd.AddCommand(imagesCmd)
	rootCmd.AddCommand(docsCmd)

	imagesCmd.AddCommand(images.NewPublishImagesCommand(options))
	docsCmd.AddCommand(docs.NewPublishDocsCommand(options))

	rootCmd.PersistentFlags().StringVar(&options.Registry, "registry", options.Registry, "target registry for the containerdisks")
	rootCmd.PersistentFlags().BoolVar(&options.AllowInsecureRegistry, "insecure-skip-tls", options.AllowInsecureRegistry, "allow connecting to insecure registries")
	rootCmd.PersistentFlags().BoolVar(&options.DryRun, "dry-run", options.DryRun, "don't publish anything")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
