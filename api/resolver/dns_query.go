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
	Parser *dnsmessage.Parser
	Header *dnsmessage.Header
}

func RunOutgoingDnsQuery(
	servers []net.IP,
	question dnsmessage.Question,
) (*dnsQueryResult, error) {

	conn, err := getServerConn(servers)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	return doDNSQuery(conn, question)
}

func doDNSQuery(conn net.Conn, question dnsmessage.Question) (*dnsQueryResult, error) {
	message, err := newDNSMessage(question)
	if err != nil {
		return nil, err
	}

	buffer, err := message.Pack()
	if err != nil {
		return &dnsQueryResult{}, nil
	}

	_, err = conn.Write(buffer)
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

	questions, err := p.AllQuestions()
	if err != nil {
		return &dnsQueryResult{}, err
	}
	if len(questions) != len(message.Questions) {
		return &dnsQueryResult{}, fmt.Errorf("answer packet doesn't have the same amount of questions")
	}
	err = p.SkipAllQuestions()
	if err != nil {
		return &dnsQueryResult{}, err
	}

	return &dnsQueryResult{
		Parser: &p,
		Header: &headers,
	}, nil
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
