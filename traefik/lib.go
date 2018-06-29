package traefik

import (
	"net/http"
	"io/ioutil"
	"errors"
	"encoding/json"
	"strings"
	"fmt"
)

func GetHosts(traefikAddress string) (hosts []string, err error) {
	var endpoint = fmt.Sprintf("http://%s/api/providers", traefikAddress)
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

				hosts = append(hosts, strings.TrimPrefix(routes.Rule, "Host:"))
			}
		}
	}

	return
}
