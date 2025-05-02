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

var env = os.Environ()

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
		configCommand(),
		deployCommand(),
	)
}
func configCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Outputs the final config file, after doing merges and interpolations",
		Args:  cobra.RangeArgs(0, 1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Set the default namespace to "default" if no argument is provided
			// or if the argument is empty
			namespace := "default"
			if len(args) > 0 {
				namespace = args[0]
			}

			// Generate a random number for the RANDOM environment variable
			r := rand.New(rand.NewSource(99))
			env = append(env, fmt.Sprintf("RANDOM=%d", r.Uint32()))

			// if env does not contains DOCKER_REGISTRY, then set it to "docker.io"
			if _, ok := os.LookupEnv("DOCKER_REGISTRY"); !ok {
				env = append(env, "DOCKER_REGISTRY=docker.io")
			}

			// Set the DOCKER_STACK_NAMESPACE environment variable
			env = append(env, fmt.Sprintf("DOCKER_STACK_NAMESPACE=%s", namespace))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Prepare the command to execute
			execArgv := []string{"docker", "stack", "config"}
			if composeFile, err := cmd.Flags().GetString("compose-file"); err == nil {
				execArgv = append(execArgv, "--compose-file="+composeFile)
			}
			if skipInterpolation, err := cmd.Flags().GetBool("skip-interpolation"); err == nil {
				execArgv = append(execArgv, "--skip-interpolation="+fmt.Sprintf("%t", skipInterpolation))
			}

			command := exec.Cmd{
				Path:   "/usr/local/bin/docker",
				Args:   execArgv,
				Env:    env,
				Stdin:  os.Stdin,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}

			return command.Run()
		},
	}

	cmd.Flags().StringP("compose-file", "c", "docker-stack.yml", "Path to a Compose file, or \"-\" to read from stdin")
	cmd.Flags().Bool("skip-interpolation", false, "Skip interpolation and output only merged config")

	return cmd
}

func deployCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a new stack or update an existing stack",
		Args:  cobra.RangeArgs(0, 1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Set the default namespace to "default" if no argument is provided
			// or if the argument is empty
			namespace := "default"
			if len(args) > 0 {
				namespace = args[0]
			}

			// Generate a random number for the RANDOM environment variable
			r := rand.New(rand.NewSource(99))
			env = append(env, fmt.Sprintf("RANDOM=%d", r.Uint32()))

			// if env does not contains DOCKER_REGISTRY, then set it to "docker.io"
			if _, ok := os.LookupEnv("DOCKER_REGISTRY"); !ok {
				env = append(env, "DOCKER_REGISTRY=docker.io")
			}

			// Set the DOCKER_STACK_NAMESPACE environment variable
			env = append(env, fmt.Sprintf("DOCKER_STACK_NAMESPACE=%s", namespace))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the default namespace to "default" if no argument is provided
			// or if the argument is empty
			namespace := "default"
			if len(args) > 0 {
				namespace = args[0]
			}

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

			command := exec.Cmd{
				Path:   "/usr/local/bin/docker",
				Args:   execArgv,
				Env:    env,
				Stdin:  os.Stdin,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}

			fmt.Printf("Deploying stack to namespace: %s\n", namespace)

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
