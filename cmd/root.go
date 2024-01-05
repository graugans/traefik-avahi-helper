/*
Copyright © 2024 Christian Ege <ch@ege.io>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var (
	defaultDockerHost string = "unix:///var/run/docker.sock"
)

func rootCmdRunE(cmd *cobra.Command, args []string) error {
	var err error
	cli, err := client.NewClientWithOpts(client.FromEnv)
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
	return err
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "traefik-avahi-helper",
	Short: "Register a CNAME in avahi based on Traefik Docker labels",
	RunE:  rootCmdRunE,
	// Do not print a usage message on failure
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
