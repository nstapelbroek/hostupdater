package cmd

import (
	"github.com/spf13/cobra"
	"github.com/nstapelbroek/hostupdater/traefik"
	"github.com/nstapelbroek/hostupdater/helper"
	"github.com/sirupsen/logrus"
	"github.com/cbednarski/hostess"
	"time"
	"regexp"
)

func init() {
	rootCmd.AddCommand(traefikCmd)
	traefikCmd.Flags().String("address", "127.0.0.1", "The IP of the Traefik server we're trying to fetch the frontend configuration from.")
	traefikCmd.Flags().Int16("port", 8080, "The port where the Traefik host is serving it's API.")
	traefikCmd.Flags().Int8("interval", 0, "Update every X seconds, use the default value of 0 for a single execution.")
	traefikCmd.Flags().Int8("wait", 0, "Wait an amount of X seconds before execution allowing services to pass health-checks and register")
	traefikCmd.Flags().String("filter", ".*", "Only update the hosts who pass this regular expression, useful when using multiple loadbalancers")
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
		filterExpression, err := regexp.Compile(filter)
		if err != nil {
			logrus.WithFields(logrus.Fields{"expression": filter}).Error("Invalid regex passed to filter")
			return
		}

		err = updateHostsFromTraefikApi(TraefikAddress, filterExpression)
		if err != nil || interval == 0 {
			return
		}

		// The interval argument is set, startup the timer
		ticker := time.NewTicker(time.Second * time.Duration(interval))
		quit := make(chan struct{})
		for {
			select {
			case <-ticker.C:
				if err := updateHostsFromTraefikApi(TraefikAddress, filterExpression); err != nil {
					return err
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	},
}

func updateHostsFromTraefikApi(address traefik.Address, filterExpression *regexp.Regexp) (err error) {
	frontendHosts, err := traefik.GetFrontendHosts(address.IP, address.PortNumber)
	if err != nil {
		logrus.WithFields(logrus.Fields{"ip": address.IP, "port": address.PortNumber}).Error(err)
		return
	}

	filteredHosts := make([]*hostess.Hostname, len(frontendHosts))
	for i, host := range frontendHosts {
		if !filterExpression.MatchString(host) {
			continue
		}

		hostName, _ := hostess.NewHostname(host, address.IP.String(), true)
		filteredHosts = append(filteredHosts[:i], hostName)
	}

	err = writeHostsToFile(filteredHosts)
	if err != nil {
		logrus.Errorln("failed persist changes to hostsfile.", err)
	}

	return
}

func writeHostsToFile(hosts []*hostess.Hostname) (err error) {
	hostfile, errors := hostess.LoadHostfile()
	if len(errors) > 0 {
		err = errors[0]
		return
	}

	shouldSave := false
	for _, host := range hosts {
		if host == nil || hostfile.Hosts.Contains(host) {
			continue
		}

		// Both duplicate and conflicts return errors so you are aware of them
		_ = hostfile.Hosts.Add(host)

		shouldSave = true
		logrus.WithFields(logrus.Fields{"domain": host.Domain, "ip": host.IP.String(),}).Info("created or updated record")
	}

	if !shouldSave {
		logrus.Debugln("no changes to hostfile needed")
		return
	}

	return hostfile.Save()
}
