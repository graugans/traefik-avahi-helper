/*
Copyright © 2024 Christian Ege <ch@ege.io>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func rootCmdRunE(cmd *cobra.Command, args []string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
	}

	labelRe, err := regexp.Compile(`traefik\.http\.routers\.(.*)\.rule`) // error if regexp invalid
	if err != nil {
		return err
	}
	domainRe, err := regexp.Compile(`(?P<domain>[^\x60]*?\.local)`) // error if regexp invalid
	if err != nil {
		return err
	}
	var hostnames []string = []string{}
	for _, container := range containers {
		fmt.Printf("Checking container ID: %s, Name: %s\n", container.ID[:10], container.Image)
		if container.Labels["traefik.enable"] == "true" {
			for key, value := range container.Labels {
				if labelRe.Match([]byte(key)) {
					match := domainRe.FindStringSubmatch(value)
					if len(match) > 0 {
						hostname := match[0]
						hostnames = append(hostnames, hostname)
						fmt.Printf("   Found %s\n", hostname)
					}
				}
			}
		}
	}

	fmt.Println("Found hostnames:")
	for i := range hostnames {
		fmt.Printf("    - %s\n", hostnames[i])
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
