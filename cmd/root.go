package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"
	"os"
)

var Verbose bool
var rootCmd = &cobra.Command{
	Use:   "hostupdater",
	Short: "hostupdater is a simple binary that manages your hostfile",
	Long: `hostupdater aims to provide just the right amount of glue between your docker setup and a hostsfile.
Use the binary to patch your hostfile with the routing information collected form popular loadbalancers like traefik.`,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf(err.Error())
		os.Exit(1)
	}
}
