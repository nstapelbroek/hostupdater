package cmd

import (
	"github.com/spf13/cobra"
	"github.com/nstapelbroek/hostupdater/traefik"
	"github.com/cbednarski/hostess"
	"fmt"
	"os"
	"github.com/nstapelbroek/hostupdater/helper"
)

func init() {
	rootCmd.AddCommand(traefikCmd)
	traefikCmd.Flags().String("address", "127.0.0.1", "The IP of the traefik loadbalancer we're trying to fetch the frontend configuration from")
	traefikCmd.Flags().Int16("port", 8080, "The port where the traefik host is serving it's API. We need this API to fetch the hosts")
}

var traefikCmd = &cobra.Command{
	Use:   "traefik",
	Short: "Retrieve host information from a traefik loadbalancer",
	Run: func(cmd *cobra.Command, args []string) {
		address, _ := cmd.Flags().GetString("address")
		port, _ := cmd.Flags().GetInt16("port")
		traefikIp, err := helper.AddressToIp(address)
		if err != nil {
			fmt.Errorf("%s", err)
			os.Exit(1)
		}

		hosts, err := traefik.GetHosts(traefikIp, port)
		if err != nil {
			fmt.Errorf("%s", err)
			os.Exit(1)
		}

		hostfile, errors := hostess.LoadHostfile()
		if len(errors) > 0 {
			for _, err := range errors {
				fmt.Errorf("%s", err)
			}
			os.Exit(1)
		}

		for _, host := range hosts {
			hostname, err := hostess.NewHostname(host, traefikIp.String(), true)
			if err != nil {
				fmt.Errorf("%s", err)
				os.Exit(1)
			}

			// Note that we skip error checking here because hostess will error on update or duplicate
			hostfile.Hosts.Add(hostname)
			fmt.Sprintln("Created or updated record for %s to %s", hostname.Domain, hostname.IP.String())
		}

		hostfile.Save()
	},
}
