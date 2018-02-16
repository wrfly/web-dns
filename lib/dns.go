package lib

import (
	"fmt"
	"net"
	"strings"

	dns "golang.org/x/net/dns/dnsmessage"
)

func newName(name string) (dns.Name, error) {
	n, err := dns.NewName(name)
	if err != nil {
		return dns.Name{}, err
	}
	return n, nil
}

func convertType(typ string) (dns.Type, error) {
	typ = strings.ToUpper(typ)
	switch typ {
	case "A":
		return dns.TypeA, nil
	case "AAAA":
		return dns.TypeAAAA, nil
	case "MX":
		return dns.TypeMX, nil
	case "CNAME":
		return dns.TypeCNAME, nil
	case "NS":
		return dns.TypeNS, nil
	case "TXT":
		return dns.TypeTXT, nil
	default:
		return dns.TypeA, fmt.Errorf("type not support")
	}
}

type Resp struct {
	IP  string
	TTL uint32
}

type Answer struct {
	result []Resp
	err    error
}

func (a Answer) IPs() []string {
	ips := []string{}
	for _, result := range a.result {
		ips = append(ips, result.IP)
	}
	return ips
}

func (a Answer) Error() error {
	return a.err
}

func Question(domain, typ string) Answer {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	// type and name
	dnsType, err := convertType(typ)
	if err != nil {
		return Answer{err: err}
	}
	dnsName, err := newName(domain)
	if err != nil {
		return Answer{err: err}
	}

	// build query message
	msg := buildQueryMessage(dnsName, dnsType)
	buf, err := msg.Pack()
	if err != nil {
		return Answer{err: err}
	}

	u, err := net.Dial("udp", "8.8.8.8:53")
	defer u.Close()
	u.Write(buf)
	got := make([]byte, 2048)
	n, err := u.Read(got)
	msg.Unpack(got[:n])

	result, err := parseMessage(msg)
	if err != nil {
		return Answer{err: err}
	}
	return Answer{result: result}
}

func buildQueryMessage(name dns.Name, typ dns.Type) (msg dns.Message) {
	msg = dns.Message{
		Header: dns.Header{RecursionDesired: true},
		Questions: []dns.Question{
			{
				Name:  name,
				Type:  typ,
				Class: dns.ClassINET,
			},
		},
	}
	return
}

func parseMessage(msg dns.Message) (resps []Resp, err error) {
	buf, err := msg.Pack()
	if err != nil {
		return nil, err
	}
	var p dns.Parser
	if _, err := p.Start(buf); err != nil {
		return nil, err
	}

	// question first
	for {
		q, err := p.Question()
		if err == dns.ErrSectionDone {
			break
		}
		if err != nil {
			return nil, err
		}

		if q.Name.String() != msg.Questions[0].Name.String() {
			continue
		}

		if err := p.SkipAllQuestions(); err != nil {
			return nil, err
		}
		break
	}

	resps = []Resp{}
	// then parse answer
	var ip string
	for {
		h, err := p.AnswerHeader()
		if err == dns.ErrSectionDone {
			break
		}
		if err != nil {
			return nil, err
		}

		if !strings.EqualFold(h.Name.String(), msg.Questions[0].Name.String()) {
			if err := p.SkipAnswer(); err != nil {
				return nil, err
			}
			continue
		}

		switch h.Type {
		case dns.TypeA:
			r, err := p.AResource()
			if err != nil {
				return nil, err
			}
			ip = net.IP(r.A[:]).String()
		case dns.TypeAAAA:
			r, err := p.AAAAResource()
			if err != nil {
				return nil, err
			}
			ip = net.IP(r.AAAA[:]).String()
		}
		resps = append(resps, Resp{ip, h.TTL})
	}

	return
}
