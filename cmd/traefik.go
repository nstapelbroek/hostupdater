package cmd

import (
	"github.com/spf13/cobra"
	"github.com/nstapelbroek/hostupdater/traefik"
	"github.com/nstapelbroek/hostupdater/helper"
	"github.com/sirupsen/logrus"
	"time"
)

func init() {
	rootCmd.AddCommand(traefikCmd)
	traefikCmd.Flags().String("address", "127.0.0.1", "The IP of the Traefik server we're trying to fetch the frontend configuration from.")
	traefikCmd.Flags().Int16("port", 8080, "The port where the Traefik host is serving it's API.")
	traefikCmd.Flags().Int8("interval", 0, "Update every X seconds, use the default value of 0 for a single execution.")
	traefikCmd.Flags().Int8("wait", 0, "Wait an amount of X seconds before execution allowing services to pass health-checks and register")
	traefikCmd.Flags().String("filter", "*", "Only update the hosts who pass this regular expression, useful when using multiple loadbalancers")
}

var traefikCmd = &cobra.Command{
	Use:   "traefik",
	Short: "Retrieve frontend routing information from a Traefik router",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if v, _ := rootCmd.PersistentFlags().GetBool("verbose"); v {
			logrus.SetLevel(logrus.DebugLevel)
		}

		waitTime, _ := cmd.Flags().GetInt8("wait")
		address, _ := cmd.Flags().GetString("address")
		port, _ := cmd.Flags().GetInt16("port")
		interval, _ := cmd.Flags().GetInt8("interval")
		filter, _ := cmd.Flags().GetString("filter")

		time.Sleep(time.Duration(waitTime) * time.Second)

		traefikIp, err := helper.AddressToIp(address)
		if err != nil {
			logrus.WithFields(logrus.Fields{"address": address}).Error("failed resolving address to a usable IP")
			return
		}

		TraefikAddress := traefik.Address{IP: traefikIp, PortNumber: port}

		err = updateHostsFromTraefikApi(TraefikAddress, filter)
		if err != nil || interval == 0 {
			return
		}

		// The interval argument is set, startup the timer
		ticker := time.NewTicker(time.Second * time.Duration(interval))
		quit := make(chan struct{})
		for {
			select {
			case <-ticker.C:
				if err := updateHostsFromTraefikApi(TraefikAddress, filter); err != nil {
					return err
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	},
}

func updateHostsFromTraefikApi(address traefik.Address, filter string) (err error) {
	hosts, err := traefik.GetHosts(address.IP, address.PortNumber)
	if err != nil {
		logrus.WithFields(logrus.Fields{"ip": address.IP, "port": address.PortNumber}).Error(err)
		return
	}

	// Filter the results here

	err = helper.WriteHostsToFile(hosts)
	if err != nil {
		logrus.Errorln("failed persist changes to hostsfile.", err)
	}

	return
}
