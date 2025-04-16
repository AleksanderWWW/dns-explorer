package resolver

import (
	"net"
	"strings"
)

type ResolutionInfo struct {
	Hostname  string `json:"HostName"`
	IpAddress string `json:"IpAddress"`
}

func ResolveHostname(hostname string) (ResolutionInfo, error) {
	return ResolutionInfo{
		Hostname: hostname,
	}, nil
}

func GetRootServers() []net.IP {
	rootServers := []net.IP{}
	for _, rootServer := range strings.Split(ROOT_SERVERS, ",") {
		rootServers = append(rootServers, net.ParseIP(rootServer))
	}
	return rootServers
}
