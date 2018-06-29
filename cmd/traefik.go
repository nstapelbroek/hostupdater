package cmd

import (
	"github.com/spf13/cobra"
	"github.com/nstapelbroek/hostupdater/traefik"
	"github.com/cbednarski/hostess"
	"fmt"
	"os"
	)

func init() {
	rootCmd.AddCommand(traefikCmd)
	traefikCmd.Flags().String("traefik-address",  "127.0.0.1:8080", "The IP of the traefik API server we're trying to fetch the frontend configuration from, this usually references your loadbalancer from the docker internal network" )
	traefikCmd.Flags().String("host-address",  "127.0.0.1", "The overwritten IP address where all frontends should point towards" )
}

var traefikCmd = &cobra.Command{
	Use:   "traefik",
	Short: "Retrieve host information from a traefik loadbalancer",
	Run: func(cmd *cobra.Command, args []string) {
		treafikIp, _ := cmd.Flags().GetString("traefik-address")
		outputIp, _ := cmd.Flags().GetString("host-address")

		hosts, err := traefik.GetHosts(treafikIp)
		if (err != nil) {
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
			hostname, err := hostess.NewHostname(host, outputIp, true)
			if(err != nil) {
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
