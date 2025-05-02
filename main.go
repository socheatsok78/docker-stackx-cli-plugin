package main

import (
	"fmt"
	"os"

	"github.com/docker/cli/cli-plugins/metadata"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

var (
	Version = "0.1.0"
	Vendor  = "github.com/socheatsok78/docker-stackx-cli-plugin"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cmd := &cobra.Command{
		Use:   "stackx",
		Short: "Docker Stack Extended",
		Long:  "Extended Docker Stack CLI plugin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cli, err := command.NewDockerCli()
	if err != nil {
		return fmt.Errorf("failed to create docker cli: %w", err)
	}

	return plugin.RunPlugin(cli, cmd, metadata.Metadata{
		SchemaVersion: "0.1.0",
		Vendor:        Vendor,
		Version:       Version,
	})
}
