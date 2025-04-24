package resolver

import "github.com/miekg/dns"

func ResolveDNS(name string) (CachedDNSResponse, uint32, error) {
	var records []string
	var minTTL uint32 = 3600 // default fallback TTL

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(name), dns.TypeA)

	c := new(dns.Client)
	in, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		return CachedDNSResponse{}, 0, err
	}

	for _, ans := range in.Answer {
		if a, ok := ans.(*dns.A); ok {
			records = append(records, a.A.String())
			if a.Hdr.Ttl < minTTL {
				minTTL = a.Hdr.Ttl
			}
		}
	}

	return CachedDNSResponse{Records: records}, minTTL, nil
}
