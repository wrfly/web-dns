package lib

import (
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

type Resp struct {
	IP  string `json:"ip"`
	TTL uint32 `json:"ttl"`
}

type Answer struct {
	Result []Resp `json:"result"`
	DigAt  int64  `json:"-"`
	Err    error  `json:"err"`
}

func (a Answer) IPs() []string {
	ips := []string{}
	for _, result := range a.Result {
		ips = append(ips, result.IP)
	}
	return ips
}

func Question(dnsserver, domain, typ string) Answer {
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

	u, err := net.Dial("udp", dnsserver)
	u.Write(buf)
	got := dnsBuf.Get().([]byte)
	n, err := u.Read(got)
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

func parseMessage(msg dns.Message) (resps []Resp, err error) {
	resps = []Resp{}
	var ip string
	for _, resource := range msg.Answers {
		h := resource.Header
		switch h.Type {
		case dns.TypeA:
			r := resource.Body.(*dns.AResource)
			ip = net.IP(r.A[:]).String()
		case dns.TypeAAAA:
			r := resource.Body.(*dns.AAAAResource)
			ip = net.IP(r.AAAA[:]).String()
		case dns.TypeMX:
			r := resource.Body.(*dns.MXResource)
			ip = r.MX.String()
		case dns.TypeNS:
			r := resource.Body.(*dns.NSResource)
			ip = r.NS.String()
		case dns.TypeTXT:
			r := resource.Body.(*dns.TXTResource)
			ip = r.Txt
		case dns.TypeCNAME:
			r := resource.Body.(*dns.CNAMEResource)
			ip = r.CNAME.String()
		default:
			return nil, fmt.Errorf("unknown query type")
		}
		resps = append(resps, Resp{ip, h.TTL})
	}

	return
}
