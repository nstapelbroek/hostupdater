package cmd

import (
	"github.com/spf13/cobra"
	"github.com/nstapelbroek/hostupdater/traefik"
	"github.com/nstapelbroek/hostupdater/helper"
	"github.com/sirupsen/logrus"
	"time"
	"net"
)

func init() {
	rootCmd.AddCommand(traefikCmd)
	traefikCmd.Flags().String("address", "127.0.0.1", "The IP of the Traefik server we're trying to fetch the frontend configuration from.")
	traefikCmd.Flags().Int16("port", 8080, "The port where the Traefik host is serving it's API.")
	traefikCmd.Flags().Int8("interval", 0, "Update every X seconds, use the default value of 0 for a single execution.")
}

var traefikCmd = &cobra.Command{
	Use:   "traefik",
	Short: "Retrieve frontend routing information from a Traefik router",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if v, _ := rootCmd.PersistentFlags().GetBool("verbose"); v {
			logrus.SetLevel(logrus.DebugLevel)
		}

		address, _ := cmd.Flags().GetString("address")
		port, _ := cmd.Flags().GetInt16("port")
		interval, _ := cmd.Flags().GetInt8("interval")

		traefikIp, err := helper.AddressToIp(address)
		if err != nil {
			logrus.WithFields(logrus.Fields{"address": address}).Error("failed resolving address to a usable IP")
			return
		}

		err = updateHostsFromTraefikApi(traefikIp, port)
		if err != nil || interval == 0 {
			return
		}

		ticker := time.NewTicker(time.Second * time.Duration(interval))
		quit := make(chan struct{})
		for {
			select {
			case <-ticker.C:
				if err := updateHostsFromTraefikApi(traefikIp, port); err != nil {
					return err
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	},
}

func updateHostsFromTraefikApi(traefikIp net.IP, port int16) (err error) {
	hosts, err := traefik.GetHosts(traefikIp, port)
	if err != nil {
		logrus.WithFields(logrus.Fields{"ip": traefikIp, "port": port}).Error(err)
		return
	}

	err = helper.WriteHostsToFile(hosts)
	if err != nil {
		logrus.Errorln("failed persist changes to hostsfile.", err)
	}

	return
}
