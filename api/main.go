package main

import (
	"fmt"
	"log"

	"github.com/AleksanderWWW/dns-explorer/resolver"
)

func main() {
	host := "google.com"

	info, err := resolver.ResolveHostname(host)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Returned info\n hostname: %s\n ip: %s\n", info.Hostname, info.IpAddress)
}
