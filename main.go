package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"

	"github.com/docker/cli/cli"
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
			if len(args) == 0 {
				return cmd.Help()
			}
			_ = cmd.Help()
			return cli.StatusError{
				StatusCode: 1,
				Status:     fmt.Sprintf("ERROR: unknown command: %q", args[0]),
			}
		},
	}

	addCommands(cmd)

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

func addCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		deployCommand(),
	)
}

func deployCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a stack",
		Long:  "Deploy a stack to the Docker Swarm cluster",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env := os.Environ()

			namespace := args[0]

			// Prepare the command to execute
			execArgv := []string{"docker", "stack", "deploy"}
			if composeFile, err := cmd.Flags().GetString("compose-file"); err == nil {
				execArgv = append(execArgv, "--compose-file="+composeFile)
			}
			if detach, err := cmd.Flags().GetBool("detach"); err == nil {
				execArgv = append(execArgv, "--detach="+fmt.Sprintf("%t", detach))
			}
			if prune, err := cmd.Flags().GetBool("prune"); err == nil {
				execArgv = append(execArgv, "--prune="+fmt.Sprintf("%t", prune))
			}
			if quiet, err := cmd.Flags().GetString("quiet"); err == nil {
				execArgv = append(execArgv, "--quiet="+quiet)
			}
			if resolveImage, err := cmd.Flags().GetString("resolve-image"); err == nil {
				execArgv = append(execArgv, "--resolve-image="+resolveImage)
			}
			if withRegistryAuth, err := cmd.Flags().GetBool("with-registry-auth"); err == nil {
				execArgv = append(execArgv, "--with-registry-auth="+fmt.Sprintf("%t", withRegistryAuth))
			}
			execArgv = append(execArgv, namespace)

			// Generate a random number for the RANDOM environment variable
			r := rand.New(rand.NewSource(99))
			env = append(env, fmt.Sprintf("RANDOM=%d", r.Uint32()))

			// Set the DOCKER_STACK_NAMESPACE environment variable
			env = append(env, fmt.Sprintf("DOCKER_STACK_NAMESPACE=%s", namespace))

			command := exec.Cmd{
				Path:   "/usr/local/bin/docker",
				Args:   execArgv,
				Env:    env,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}

			return command.Run()
		},
	}

	cmd.Flags().StringP("compose-file", "c", "docker-stack.yml", "Path to a Compose file, or \"-\" to read from stdin")
	cmd.Flags().BoolP("detach", "d", true, "Exit immediately instead of waiting for the stack services to converge")
	cmd.Flags().Bool("prune", false, "Prune services that are no longer referenced")
	cmd.Flags().BoolP("quiet", "q", false, "Suppress progress output")
	cmd.Flags().String("resolve-image", "always", "Query the registry to resolve image digest and supported platforms (\"always\", \"changed\", \"never\")")
	cmd.Flags().Bool("with-registry-auth", false, "Send registry authentication details to Swarm agents")

	return cmd
}
