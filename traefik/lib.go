package traefik

import (
	"net/http"
	"io/ioutil"
	"errors"
	"encoding/json"
	"strings"
	"fmt"
	"net"
	"github.com/cbednarski/hostess"
)

func GetHosts(traefikIp net.IP, traefikPort int16) (hosts []*hostess.Hostname, err error) {
	var endpoint = fmt.Sprintf("http://%s:%d/api/providers", traefikIp.String(), traefikPort)
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("api providers endpoint responded with a non 200 status code")
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	var providers ProvidersResponse
	err = json.Unmarshal(responseBody, &providers)
	if err != nil {
		return
	}

	for _, provider := range providers {
		for _, frontend := range provider.Frontends {
			for _, routes := range frontend.Routes {
				if !strings.HasPrefix(routes.Rule, "Host:") {
					continue
				}

				hostname, err := hostess.NewHostname(strings.TrimPrefix(routes.Rule, "Host:"), traefikIp.String(), true)
				if err != nil {
					continue
				}

				hosts = append(hosts, hostname)
			}
		}
	}

	return
}
