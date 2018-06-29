package traefik

import (
	"net/http"
	"io/ioutil"
	"errors"
	"encoding/json"
	"strings"
	"github.com/nstapelbroek/hostupdater/domain"
	"net"
)

func GetHosts() (hosts []*domain.Hostname, err error) {
	ip := net.ParseIP("127.0.0.1")
	var endpoint = "http://127.0.0.1:8080/api/providers" // @todo make Traefik endpoint configurable
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("No valid response status")
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
				if (!strings.HasPrefix(routes.Rule, "Host:")) {
					continue
				}

				var newHost = domain.NewHostname(strings.TrimPrefix(routes.Rule, "Host:"), ip)
				hosts = append(hosts, newHost)
			}
		}
	}

	return
}
