package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Options struct {
	AllowInsecureRegistry bool
	Registry              string
	DryRun                bool
	PublishOptions        PublishOptions
}

func main() {

	options := &Options{
		AllowInsecureRegistry: false,
		Registry:              "quay.io/containerdisks",
		DryRun:                true,
	}

	rootCmd := &cobra.Command{
		Use:   "medius",
		Short: "medius determines if new OS images are released and publishes them as containerdisks",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	rootCmd.AddCommand(NewPublishCommand(options))

	rootCmd.PersistentFlags().StringVar(&options.Registry, "registry", options.Registry, "target registry for the containerdisks")
	rootCmd.PersistentFlags().BoolVar(&options.AllowInsecureRegistry, "insecure-skip-tls", options.AllowInsecureRegistry, "allow connecting to insecure registries")
	rootCmd.PersistentFlags().BoolVar(&options.DryRun, "dry-run", options.DryRun, "don't publish anything")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
