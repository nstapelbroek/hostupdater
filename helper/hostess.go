package helper

import (
	"github.com/cbednarski/hostess"
	"github.com/sirupsen/logrus"
)

func WriteHostsToFile(hosts []*hostess.Hostname) (err error) {
	hostfile, errors := hostess.LoadHostfile()
	if len(errors) > 0 {
		err = errors[0]
		return
	}

	shouldSave := false
	for _, host := range hosts {
		if hostfile.Hosts.Contains(host) {
			continue
		}

		shouldSave = true
		hostfile.Hosts.Add(host)
		logrus.WithFields(logrus.Fields{"domain": host.Domain, "ip": host.IP.String(),}).Info("created or updated record")
	}

	if !shouldSave {
		logrus.Infoln("no changes to hostfile needed")
		return
	}

	err = hostfile.Save()
	return
}