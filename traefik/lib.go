package traefik

import (
	"net/http"
		"io/ioutil"
		"errors"
)

type ProvidersResponse map[string]Provider
type BackendCollection map[string]Backend
type FrontendCollection map[string]Frontend

type Provider struct {
	Backends  BackendCollection
	Frontends FrontendCollection
}

type Backend struct {
	Servers map[string]Server
}

type Server struct {
	Url    string
	Weight int
}

type Frontend struct {
	Routes map[string]Route
}

type Route struct {
	Rule string
}

func GetDomains() (output string, err error) {
	var endpoint = "http://localhost:8080/api/providers"
	output = ""

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
		return output, errors.New("No valid response status")
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	return string(responseBody), nil
}
