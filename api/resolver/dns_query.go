package resolver

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"net"

	"golang.org/x/net/dns/dnsmessage"
)

type dnsQueryResult struct {
	parser *dnsmessage.Parser
	header *dnsmessage.Header
}

func runOutgoingDnsQuery(
	servers []net.IP,
	question dnsmessage.Question,
) (*dnsQueryResult, error) {

	conn, err := getServerConn(servers)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	buf, err := newMessageBuffer(question)
	if err != nil {
		return nil, err
	}

	return doDNSQuery(conn, buf)
}

func doDNSQuery(conn net.Conn, buffer []byte) (*dnsQueryResult, error) {
	_, err := conn.Write(buffer)
	if err != nil {
		return nil, err
	}

	answer := make([]byte, 512)
	n, err := bufio.NewReader(conn).Read(answer)
	if err != nil {
		return nil, err
	}

	var p dnsmessage.Parser

	headers, err := p.Start(answer[:n])
	if err != nil {
		return nil, err
	}

	// questions, err := p.AllQuestions()
	// if err != nil {
	// 	return nil, err
	// }

	// if len(questions) != len(message.Questions) {
	// 	return nil, fmt.Errorf(
	// 		"answer packet has %d questions, expected %d",
	// 		len(questions),
	// 		len(message.Questions),
	// 	)
	// }

	return &dnsQueryResult{
		parser: &p,
		header: &headers,
	}, nil
}

func newMessageBuffer(question dnsmessage.Question) ([]byte, error) {
	message, err := newDNSMessage(question)
	if err != nil {
		return nil, err
	}

	return message.Pack()
}

func newDNSMessage(question dnsmessage.Question) (dnsmessage.Message, error) {
	max := ^uint16(0)
	randomNum, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return dnsmessage.Message{}, err
	}

	return dnsmessage.Message{
		Header: dnsmessage.Header{
			ID: uint16(randomNum.Int64()),
		},
		Questions: []dnsmessage.Question{question},
	}, nil
}

func getServerConn(servers []net.IP) (net.Conn, error) {
	var conn net.Conn
	var err error

	for _, server := range servers {
		conn, err = net.Dial("udp", server.String()+":53")

		if err == nil {
			return conn, nil
		}
	}

	if conn == nil {
		return nil, fmt.Errorf("failed to make connection, %s", err)
	}

	return conn, nil
}
