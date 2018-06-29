package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "hostupdater",
	Short: "HostUpdater is a simple binary that manages your hostfile",
	Long: `Todo: needs a long description. See https://nstapelbroek.com`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}