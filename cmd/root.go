/*
Copyright © 2024 Christian Ege <ch@ege.io>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/graugans/traefik-avahi-helper/internal"
	"github.com/spf13/cobra"
)

func hostNameReceiver(ch <-chan string) {
	for {
		hostname := <-ch
		fmt.Printf("Found %s\n", hostname)
	}
}

func handleAlreadyRunningContainers(cli *client.Client, hostNameChannel chan string) error {
	fmt.Println("Scanning already running containers...")
	parser := internal.NewLabelParser()
	ctxWithBackground := context.Background()
	containers, err := cli.ContainerList(ctxWithBackground, types.ContainerListOptions{})
	if err != nil {
		return err
	}

	for _, container := range containers {
		fmt.Printf("Checking container ID: %s, Name: %s\n", container.ID[:10], container.Image)
		if parser.IsTraefikEnabled(container.Labels) {
			hostname, err := parser.FindLinkLocalHostName(container.Labels)
			if err != nil {
				// No Hostname found
				continue
			}
			hostNameChannel <- hostname
		}
	}
	return nil
}

func handleAttachContainerEvents(cli *client.Client, hostNameChannel chan string) {
	fmt.Println("Waiting for new Containers being attached...")
	parser := internal.NewLabelParser()
	msgs, errs := cli.Events(context.Background(), types.EventsOptions{})
	for {
		select {
		case err := <-errs:
			fmt.Println("Error: ", err)
		case msg := <-msgs:
			if msg.Type == "container" && msg.Action == "attach" {
				fmt.Println("Type: ", msg.Type, "Action: ", msg.Action)
				if parser.IsTraefikEnabled(msg.Actor.Attributes) {
					hostname, err := parser.FindLinkLocalHostName(msg.Actor.Attributes)
					if err != nil {
						// No Hostname found
						continue
					}
					hostNameChannel <- hostname
				}
			}
		}
	}
}

func rootCmdRunE(cmd *cobra.Command, args []string) error {
	var wg sync.WaitGroup
	var err error
	var cli *client.Client

	hostNameChannel := make(chan string)
	wg.Add(1)
	go func() {
		hostNameReceiver(hostNameChannel)
		wg.Done()
	}()

	cli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	err = handleAlreadyRunningContainers(cli, hostNameChannel)
	if err != nil {
		return err
	}
	// Give the hostNameReceiver a chance to do its job
	runtime.Gosched()

	// Handle Attach events of new containers being added
	wg.Add(1)
	go func() {
		handleAttachContainerEvents(cli, hostNameChannel)
		wg.Done()
	}()

	wg.Wait()
	return nil
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
