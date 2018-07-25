package traefik

import (
	"net/http"
	"io/ioutil"
	"errors"
	"encoding/json"
	"strings"
	"fmt"
	"net"
)

func GetHosts(traefikIp net.IP, traefikPort int16) (hosts []string, err error) {
	// We'll make the traefikPort argument optional by overwriting it's default value with the default HTTP traefikPort
	if traefikPort == 0 {
		traefikPort = 80
	}

	var endpoint = fmt.Sprintf("http://%s:%d/api/providers", traefikIp.String(), traefikPort)
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
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

				hosts = append(hosts, strings.TrimPrefix(routes.Rule, "Host:"))
			}
		}
	}

	return
}
