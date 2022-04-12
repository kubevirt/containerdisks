package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"kubevirt.io/containerdisks/cmd/medius/common"
	"kubevirt.io/containerdisks/cmd/medius/docs"
	"kubevirt.io/containerdisks/cmd/medius/images"
)

func main() {
	options := &common.Options{
		DryRun:   true,
		Registry: "quay.io/containerdisks",
		ImagesOptions: common.ImagesOptions{
			ResultsFile: "results.json",
			Workers:     1,
		},
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
	imagesCmd.AddCommand(images.NewVerifyImagesCommand(options))
	docsCmd.AddCommand(docs.NewPublishDocsCommand(options))

	rootCmd.PersistentFlags().StringVar(&options.Registry, "registry", options.Registry, "target registry for the containerdisks")
	rootCmd.PersistentFlags().BoolVar(&options.AllowInsecureRegistry, "insecure-skip-tls", options.AllowInsecureRegistry, "allow connecting to insecure registries")
	rootCmd.PersistentFlags().BoolVar(&options.DryRun, "dry-run", options.DryRun, "don't publish anything")
	rootCmd.PersistentFlags().StringVar(&options.Focus, "focus", options.Focus, "Focus on a specific containerdisk")
	imagesCmd.PersistentFlags().StringVar(&options.ImagesOptions.ResultsFile, "results-file", options.ImagesOptions.ResultsFile, "File to store/read results of operations")
	imagesCmd.PersistentFlags().IntVar(&options.ImagesOptions.Workers, "workers", options.ImagesOptions.Workers, "Number of parallel workers")

	ctx, cancel := getInterruptibleContext()
	defer cancel()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		logrus.Fatal(err)
	}
}

func getInterruptibleContext() (context.Context, func()) {
	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	cancel := func() {
		signal.Stop(signalChan)
		cancelCtx()
	}

	go func() {
		select {
		case <-signalChan:
			cancelCtx()
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}
