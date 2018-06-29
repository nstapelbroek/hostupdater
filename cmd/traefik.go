package cmd

import (
	"github.com/spf13/cobra"
	"github.com/nstapelbroek/hostupdater/traefik"
	"fmt"
)

func init() {
	rootCmd.AddCommand(traefikCmd)
}

var traefikCmd = &cobra.Command{
	Use:   "traefik",
	Short: "Retrieve host information from a traefik loadbalancer",
	Long:  `Retrieve host information from a traefik loadbalancer.
			If no --source flag is passed, it will fetch from http://localhost:8080`,
	Run: func(cmd *cobra.Command, args []string) {
		domains, _ := traefik.GetHosts()
		fmt.Println(domains)
	},
}

