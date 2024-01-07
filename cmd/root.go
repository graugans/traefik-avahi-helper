/*
Copyright © 2024 Christian Ege <ch@ege.io>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"slices"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/graugans/traefik-avahi-helper/internal"
	"github.com/graugans/traefik-avahi-helper/internal/avahi"
	"github.com/spf13/cobra"
)

type hostnameStatus int

const (
	hostIsAdded   = 0
	hostIsRemoved = 1
)

type hostNameStatus struct {
	name  string
	state hostnameStatus
}

func hostNameReceiver(publisher *avahi.Publisher, ch <-chan hostNameStatus) {
	hosts := []string{}
	for {
		host := <-ch
		if host.state == hostIsAdded {
			fmt.Printf("Add CNAME: %s\n", host.name)
			if !slices.Contains(hosts, host.name) {
				hosts = append(hosts, host.name)
			}

		} else {
			fmt.Printf("Remove CNAME: %s\n", host.name)
			if slices.Contains(hosts, host.name) {
				index := slices.Index(hosts, host.name)
				hosts = append(hosts[:index], hosts[index+1:]...)
			}
			fmt.Printf("TODO remove CNAME form AVAHI")
		}
		if len(hosts) > 0 {
			fmt.Println("The host names we manage:")
			for i := range hosts {
				fmt.Println("   - ", hosts[i])
			}
		}
		err := publisher.PublishCNAMES(hosts, 600)
		if err != nil {
			fmt.Println("Error while Publishing CNAMES: ", err)
		}
	}
}

func handleAlreadyRunningContainers(cli *client.Client, hostNameChannel chan hostNameStatus) error {
	fmt.Println("Scanning already running containers with \"traefik.enable=true\" label ...")
	parser := internal.NewLabelParser()
	ctxWithBackground := context.Background()
	containers, err := cli.ContainerList(ctxWithBackground, types.ContainerListOptions{})
	if err != nil {
		return err
	}

	for _, container := range containers {
		if parser.IsTraefikEnabled(container.Labels) {
			fmt.Println("Found container with \"traefik.enable=true\" label")
			fmt.Printf("    - ID  : %s\n    - Name: %s\n",
				container.ID[:10],
				container.Image,
			)
			hostname, err := parser.FindLinkLocalHostName(container.Labels)
			if err != nil {
				// No Hostname found
				continue
			}
			hostNameChannel <- hostNameStatus{name: hostname, state: hostIsAdded}
		}
	}
	return nil
}

func handleAttachContainerEvents(cli *client.Client, hostNameChannel chan hostNameStatus) {
	fmt.Println("Waiting for containers events with the \"traefik.enable=true\" label ...")
	parser := internal.NewLabelParser()
	msgs, errs := cli.Events(context.Background(), types.EventsOptions{})
	for {
		select {
		case err := <-errs:
			fmt.Println("Error: ", err)
		case msg := <-msgs:
			// Ignore containers without Traefik being enabled
			if parser.IsTraefikEnabled(msg.Actor.Attributes) &&
				(msg.Action == "start" || msg.Action == "stop") {
				fmt.Println("Found container with \"traefik.enable=true\" label")
				fmt.Printf(
					"    - ID    : %s\n    - Name  : %s\n    - Action: %s\n",
					msg.Actor.ID[:10],
					msg.Actor.Attributes["image"],
					msg.Action,
				)
				if msg.Type == "container" && msg.Action == "start" {
					hostname, err := parser.FindLinkLocalHostName(msg.Actor.Attributes)
					if err != nil {
						// No Link Local Hostname found
						continue
					}
					hostNameChannel <- hostNameStatus{name: hostname, state: hostIsAdded}
				}
				if msg.Type == "container" && msg.Action == "stop" {
					hostname, err := parser.FindLinkLocalHostName(msg.Actor.Attributes)
					if err != nil {
						// No Link Local Hostname found
						continue
					}
					hostNameChannel <- hostNameStatus{name: hostname, state: hostIsRemoved}
				}
			}
		}
	}
}

func rootCmdRunE(cmd *cobra.Command, args []string) error {
	var wg sync.WaitGroup
	var err error
	var cli *client.Client

	fmt.Println("Creating publisher")
	publisher, err := avahi.NewPublisher()
	if err != nil {
		return fmt.Errorf("failed to create publisher: %w", err)
	}

	fqdn := publisher.Fqdn()
	fmt.Printf("FQDN from Avahi: %s\n", fqdn)

	hostNameChannel := make(chan hostNameStatus)
	wg.Add(1)
	go func() {
		hostNameReceiver(publisher, hostNameChannel)
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
