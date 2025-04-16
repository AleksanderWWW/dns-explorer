package resolver

import (
	"testing"
	"strings"
	"net"

	"golang.org/x/net/dns/dnsmessage"
)

func TestOutgoingDnsQuery(t *testing.T) {
	question := dnsmessage.Question{
		Name:  dnsmessage.MustNewName("com."),
		Type:  dnsmessage.TypeNS,
		Class: dnsmessage.ClassINET,
	}
	rootServers := strings.Split(ROOT_SERVERS, ",")
	if len(rootServers) == 0 {
		t.Fatalf("No root servers found")
	}
	servers := []net.IP{net.ParseIP(rootServers[0])}
	result, err := RunOutgoingDnsQuery(servers, question)
	if err != nil {
		t.Fatalf("outgoingDnsQuery error: %s", err)
	}
	if result.Header == nil {
		t.Fatalf("No header found")
	}

	if result.Header.RCode != dnsmessage.RCodeSuccess {
		t.Fatalf("response was not succesful (maybe the DNS server has changed?)")
	}

	err = result.Parser.SkipAllAnswers()
	if err != nil {
		t.Fatalf("SkipAllAnswers error: %s", err)
	}
	parsedAuthorities, err := result.Parser.AllAuthorities()
	if err != nil {
		t.Fatalf("Error getting answers")
	}
	if len(parsedAuthorities) == 0 {
		t.Fatalf("No answers received")
	}
}
