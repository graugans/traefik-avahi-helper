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

func HostNameReceiver(ch <-chan string) {
	for {
		hostname := <-ch
		fmt.Printf("Found %s\n", hostname)
	}
}

func rootCmdRunE(cmd *cobra.Command, args []string) error {

	hostNameChannel := make(chan string)
	go HostNameReceiver(hostNameChannel)

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	ctxWithBackground := context.Background()
	containers, err := cli.ContainerList(ctxWithBackground, types.ContainerListOptions{})
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
	for _, container := range containers {
		fmt.Printf("Checking container ID: %s, Name: %s\n", container.ID[:10], container.Image)
		if container.Labels["traefik.enable"] == "true" {
			for key, value := range container.Labels {
				if labelRe.Match([]byte(key)) {
					match := domainRe.FindStringSubmatch(value)
					if len(match) > 0 {
						hostname := match[0]
						hostNameChannel <- hostname
					}
				}
			}
		}
	}

	msgs, errs := cli.Events(context.Background(), types.EventsOptions{})

	for {
		select {
		case err := <-errs:
			fmt.Println("Error: ", err)
		case msg := <-msgs:
			if msg.Type == "container" && msg.Action == "attach" {
				fmt.Println("Type: ", msg.Type, "Action: ", msg.Action)
				if msg.Actor.Attributes["traefik.enable"] == "true" {
					for key, value := range msg.Actor.Attributes {
						if labelRe.Match([]byte(key)) {
							match := domainRe.FindStringSubmatch(value)
							if len(match) > 0 {
								hostname := match[0]
								hostNameChannel <- hostname
							}
						}
					}
				}
			}
		}
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
