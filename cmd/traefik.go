package cmd

import (
	"github.com/spf13/cobra"
	"github.com/nstapelbroek/hostupdater/traefik"
	"github.com/cbednarski/hostess"

	"fmt"
	)

func init() {
	rootCmd.AddCommand(traefikCmd)
}

var traefikCmd = &cobra.Command{
	Use:   "traefik",
	Short: "Retrieve host information from a traefik loadbalancer",
	Long:  `Retrieve host information from a traefik loadbalancer.
			If no --source flag is passed, it will fetch from http://127.0.0.1:8080`,
	Run: func(cmd *cobra.Command, args []string) {
		domains, _ := traefik.GetHosts()

		hostsfile, _:= hostess.LoadHostfile()

		for _, domain := range domains {
			hostname, _ := hostess.NewHostname(domain.Name, domain.Address.String(), true)
			hostsfile.Hosts.ContainsDomain(hostname.Domain)
			hostsfile.Hosts.Add(hostname)
			//fmt.Printf("%s", hostsfile.Format())
			hostsfile.Save()
		}

		fmt.Println(domains)
	},
}

