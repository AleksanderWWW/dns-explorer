package resolver

import (
	"net"
	"strings"
)

var rootServers []net.IP

func init() {
	rootServers = getRootServers()
}

type ResolutionInfo struct {
	Hostname  string `json:"HostName"`
	IpAddress string `json:"IpAddress"`
}

func Resolve_hostname(hostname string) (ResolutionInfo, error) {
	return ResolutionInfo{
		Hostname: hostname,
	}, nil
}

func getRootServers() []net.IP {
	rootServers := []net.IP{}
	for _, rootServer := range strings.Split(ROOT_SERVERS, ",") {
		rootServers = append(rootServers, net.ParseIP(rootServer))
	}
	return rootServers
}
