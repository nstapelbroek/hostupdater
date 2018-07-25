package cmd

import (
	"github.com/spf13/cobra"
	"github.com/nstapelbroek/hostupdater/traefik"
	"github.com/cbednarski/hostess"
	"github.com/nstapelbroek/hostupdater/helper"
	"github.com/Sirupsen/logrus"
	"os"
)

func init() {
	rootCmd.AddCommand(traefikCmd)
	traefikCmd.Flags().String("address", "127.0.0.1", "The IP of the Traefik server we're trying to fetch the frontend configuration from.")
	traefikCmd.Flags().Int16("port", 8080, "The port where the Traefik host is serving it's API.")
}

var traefikCmd = &cobra.Command{
	Use:   "traefik",
	Short: "Retrieve frontend routing information from a Traefik router",
	Run: func(cmd *cobra.Command, args []string) {
		address, _ := cmd.Flags().GetString("address")
		port, _ := cmd.Flags().GetInt16("port")
		traefikIp, err := helper.AddressToIp(address)
		if err != nil {
			logrus.WithFields(logrus.Fields{"address": address}).Error("failed resolving address to a usable IP")
			os.Exit(1)
		}

		hosts, err := traefik.GetHosts(traefikIp, port)
		if err != nil {
			logrus.WithFields(logrus.Fields{"ip": traefikIp, "port": port}).Error(err)
			os.Exit(1)
		}

		hostfile, errors := hostess.LoadHostfile()
		if len(errors) > 0 {
			for _, err := range errors {
				logrus.Errorln(err)
			}
			os.Exit(1)
		}

		for _, host := range hosts {
			hostname, err := hostess.NewHostname(host, traefikIp.String(), true)
			if err != nil {
				logrus.Errorln(err)
				os.Exit(1)
			}

			// Note that we skip error checking here because hostess will error on update or duplicate
			hostfile.Hosts.Add(hostname)
			logrus.WithFields(logrus.Fields{"domain": hostname.Domain, "ip": hostname.IP.String(),}).Info("created or updated record")
		}

		err = hostfile.Save()
		if err != nil {
			logrus.Errorln("failed persist changes to hostsfile.", err)
			os.Exit(1)
		} else {
			logrus.Infoln("changes persisted into hostsfile")
		}
	},
}
