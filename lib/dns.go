package lib

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	dns "golang.org/x/net/dns/dnsmessage"
)

var dnsBuf = sync.Pool{
	New: func() interface{} {
		return make([]byte, 2048)
	},
}

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

type Result struct {
	Host string `json:"host"`
	Type string `json:"type"`
	TTL  uint32 `json:"ttl"`
}

type Answer struct {
	Result []Result `json:"result"`
	DigAt  int64    `json:"dig"`
	Err    error    `json:"err"`
}

func (a Answer) Hosts() []string {
	hosts := []string{}
	for _, result := range a.Result {
		hosts = append(hosts, result.Host)
	}
	return hosts
}

func (a Answer) Marshal() []byte {
	bs, _ := json.Marshal(a)
	return bs
}

func Question(dnsserver, domain, typ string, timeout time.Duration) Answer {
	logrus.Debugf("dns: %s, domain: %s, type: %s",
		dnsserver, domain, typ)
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	// type and name
	dnsType, err := convertType(typ)
	if err != nil {
		return Answer{Err: err}
	}
	dnsName, err := newName(domain)
	if err != nil {
		return Answer{Err: err}
	}

	// build query message
	msg := buildQueryMessage(dnsName, dnsType)
	buf, err := msg.Pack()
	if err != nil {
		return Answer{Err: err}
	}

	t := time.Now().Add(timeout)
	u, err := net.Dial("udp", dnsserver)
	if err := u.SetWriteDeadline(t); err != nil {
		return Answer{Err: err}
	}
	u.Write(buf)

	got := dnsBuf.Get().([]byte)
	if err := u.SetReadDeadline(t); err != nil {
		return Answer{Err: err}
	}
	n, err := u.Read(got)
	if err != nil {
		return Answer{Err: err}
	}

	u.Close()
	msg.Unpack(got[:n])
	dnsBuf.Put(got)

	result, err := parseMessage(msg)
	if err != nil {
		return Answer{Err: err}
	}
	logrus.Debugf("got answer: %v", result)
	return Answer{Result: result, DigAt: time.Now().Unix()}
}

func buildQueryMessage(name dns.Name, typ dns.Type) (msg dns.Message) {
	msg = dns.Message{
		Header:    dns.Header{RecursionDesired: true},
		Questions: []dns.Question{{Name: name, Type: typ, Class: dns.ClassINET}},
	}
	return
}

func parseMessage(msg dns.Message) (results []Result, err error) {
	results = []Result{}
	var (
		host string
		typ  string
	)
	for _, resource := range msg.Answers {
		h := resource.Header
		switch h.Type {
		case dns.TypeA:
			r := resource.Body.(*dns.AResource)
			host = net.IP(r.A[:]).String()
			typ = "A"
		case dns.TypeAAAA:
			r := resource.Body.(*dns.AAAAResource)
			host = net.IP(r.AAAA[:]).String()
			typ = "AAAA"
		case dns.TypeMX:
			r := resource.Body.(*dns.MXResource)
			host = r.MX.String()
			typ = "MX"
		case dns.TypeNS:
			r := resource.Body.(*dns.NSResource)
			host = r.NS.String()
			typ = "NS"
		case dns.TypeTXT:
			r := resource.Body.(*dns.TXTResource)
			host = r.Txt
			typ = "TXT"
		case dns.TypeCNAME:
			r := resource.Body.(*dns.CNAMEResource)
			host = r.CNAME.String()
			typ = "CNAME"
		default:
			return nil, fmt.Errorf("unknown query type")
		}
		results = append(results, Result{host, typ, h.TTL})
	}

	return
}
