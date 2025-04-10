package resolver

import (
	"testing"

	"golang.org/x/net/dns/dnsmessage"
)

func TestOutgoingDNSQuery(t *testing.T) {
	question := dnsmessage.Question{
		Name:  dnsmessage.MustNewName("com."),
		Type:  dnsmessage.TypeNS,
		Class: dnsmessage.ClassINET,
	}
	_, err := runOutgoingDnsQuery(
		getRootServers(), question,
	)

	if err != nil {
		t.Error(err)
	}
}
