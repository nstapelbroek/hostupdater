package traefik

type ProvidersResponse map[string]Provider
type BackendCollection map[string]Backend
type FrontendCollection map[string]Frontend
type RoutesCollection map[string]Route
type ServersCollection map[string]Server

type Provider struct {
	Backends  BackendCollection `json:"backends"`
	Frontends FrontendCollection `json:"frontends"`
}

type Backend struct {
	Servers ServersCollection
	// todo: add Loadbalancer entry: loadBalancer: {
	// method: "wrr"
	// }
}

type Server struct {
	Url    string
	Weight int
}

type Route struct {
	Rule string `json:"rule"`
}

type Frontend struct {
	Routes         RoutesCollection `json:"routes"`
	EntryPoints    []string `json:"entryPoints"`
	Backend        string `json:"backend"`
	PassHostHeader bool `json:"passHostHeader"`
	Priority       int32 `json:"priority"`
}


