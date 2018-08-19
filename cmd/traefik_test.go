package cmd

import (
	"testing"
	"regexp"
	"github.com/stretchr/testify/assert"
	"github.com/cbednarski/hostess"
)

type filterHostTestSetup struct {
	regex     *regexp.Regexp
	hosts     []string
	traefikIP string
}

func newFilterHostTestSetup(regex *regexp.Regexp, hosts []string) filterHostTestSetup {
	return filterHostTestSetup{regex: regex, hosts: hosts, traefikIP: "127.0.0.1"}
}

func TestFilterHostWithNilValues(t *testing.T) {
	setup := newFilterHostTestSetup(nil, nil)

	r := filterHosts(setup.hosts, setup.regex, setup.traefikIP)

	assert.Empty(t, r)
}

func TestFilterHostsWithNilRegex(t *testing.T) {
	hosts := []string{"somedomain.com", "subdomain.somedomain.com"}
	setup := newFilterHostTestSetup(nil, hosts)
	expectedHostOne, _ := hostess.NewHostname(setup.hosts[0], setup.traefikIP, true)
	expectedHostTwo, _ := hostess.NewHostname(setup.hosts[1], setup.traefikIP, true)

	r := filterHosts(setup.hosts, setup.regex, setup.traefikIP)

	assert.NotEmpty(t, r)
	assert.Len(t, r, 2)
	assert.Equal(t, expectedHostOne, r[0])
	assert.Equal(t, expectedHostTwo, r[1])
}

func TestFilterHostsWithEmptyStringRegex(t *testing.T) {
	hosts := []string{"somedomain.com"}
	regex, _ := regexp.Compile("")
	setup := newFilterHostTestSetup(regex, hosts)

	r := filterHosts(setup.hosts, setup.regex, setup.traefikIP)

	assert.Len(t, r, 1)
}

func TestFilterHostWithNoMatchers(t *testing.T) {
	hosts := []string{"somedomain.com", "subdomain.somedomain.com"}
	regex, _ := regexp.Compile("unabletomatch")
	setup := newFilterHostTestSetup(regex, hosts)

	r := filterHosts(setup.hosts, setup.regex, setup.traefikIP)

	assert.Empty(t, r, )
	assert.Len(t, r, 0)
}

func TestFilterHostsWithTopLevelDomainRegex(t *testing.T) {
	hosts := []string{"somedomain.co.uk", "somedomain.com", "subdomain.somedomain.com", "somedomain.biz", "mydomain.com"}
	regex, _ := regexp.Compile(".*.com")
	setup := newFilterHostTestSetup(regex, hosts)
	expectedHostOne, _ := hostess.NewHostname(setup.hosts[1], setup.traefikIP, true)
	expectedHostTwo, _ := hostess.NewHostname(setup.hosts[2], setup.traefikIP, true)
	expectedHostThree, _ := hostess.NewHostname(setup.hosts[4], setup.traefikIP, true)

	r := filterHosts(setup.hosts, setup.regex, setup.traefikIP)

	assert.NotEmpty(t, r)
	assert.Len(t, r, 3)
	assert.Equal(t, expectedHostOne, r[0])
	assert.Equal(t, expectedHostTwo, r[1])
	assert.Equal(t, expectedHostThree, r[2])
}

func TestFilterHostsWithRegularDomainRegex(t *testing.T) {
	hosts := []string{"somedomain.co.uk", "somedomain.com", "subdomain.somedomain.com", "somedomain.biz", "mydomain.com"}
	regex, _ := regexp.Compile(".*somedomain.*")
	setup := newFilterHostTestSetup(regex, hosts)
	expectedHostOne, _ := hostess.NewHostname(setup.hosts[0], setup.traefikIP, true)
	expectedHostTwo, _ := hostess.NewHostname(setup.hosts[1], setup.traefikIP, true)
	expectedHostThree, _ := hostess.NewHostname(setup.hosts[2], setup.traefikIP, true)
	expectedHostFour, _ := hostess.NewHostname(setup.hosts[3], setup.traefikIP, true)

	r := filterHosts(setup.hosts, setup.regex, setup.traefikIP)

	assert.NotEmpty(t, r)
	assert.Len(t, r, 4)
	assert.Equal(t, expectedHostOne, r[0])
	assert.Equal(t, expectedHostTwo, r[1])
	assert.Equal(t, expectedHostThree, r[2])
	assert.Equal(t, expectedHostFour, r[3])
}
