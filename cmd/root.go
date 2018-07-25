package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "hostupdater",
	Short: "hostupdater is a simple binary that manages your hostfile",
	Long: `hostupdater aims to provide just the right amount of glue between your docker setup and a hostsfile.
Use the binary to patch your hostfile with the routing information collected form popular local development solutions like minicube, docker-compose, traefik, vagrant or directly from the docker socket.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
