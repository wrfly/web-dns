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

// Result ...
type Result struct {
	Host string `json:"host"`
	Type string `json:"type"`
	TTL  uint32 `json:"ttl"`

	DNSType dns.Type `json:"-"`
}

// Answer ...
type Answer struct {
	Result []Result `json:"result"`
	DigAt  int64    `json:"dig"`
	Err    error    `json:"err"`
}

// Hosts ...
func (a Answer) Hosts(typ ...dns.Type) []string {
	hosts := []string{}
	for _, result := range a.Result {
		if result.DNSType == dns.TypeA {
			hosts = append(hosts, result.Host)
		}
	}
	return hosts
}

// Marshal ...
func (a Answer) Marshal() []byte {
	bs, _ := json.Marshal(a)
	return bs
}

// QueryOption ...
type QueryOption struct {
	NSServer string
	Domain   string
	Type     string

	Timeout time.Duration
}

// Question ...
func Question(opt QueryOption) Answer {
	logrus.Debugf("dns: %s, domain: %s, type: %s", opt.NSServer, opt.Domain, opt.Type)
	if !strings.HasSuffix(opt.Domain, ".") {
		opt.Domain += "."
	}

	// type and name
	dnsType, err := convertType(opt.Type)
	if err != nil {
		return Answer{Err: err}
	}
	dnsName, err := newName(opt.Domain)
	if err != nil {
		return Answer{Err: err}
	}

	// build query message
	msg := buildQueryMessage(dnsName, dnsType)
	buf, err := msg.Pack()
	if err != nil {
		return Answer{Err: err}
	}

	u, err := net.Dial("udp", opt.NSServer)
	if err != nil {
		return Answer{Err: err}
	}

	err = u.SetWriteDeadline(time.Now().Add(opt.Timeout))
	if err != nil {
		return Answer{Err: err}
	}
	u.Write(buf)

	got := dnsBuf.Get().([]byte)
	err = u.SetReadDeadline(time.Now().Add(opt.Timeout))
	if err != nil {
		return Answer{Err: err}
	}
	n, err := u.Read(got)
	if err != nil {
		return Answer{Err: err}
	}

	u.Close()
	msg.Unpack(got[:n])
	dnsBuf.Put(got)

	result, err := parsednsMessage(msg)
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

func parsednsMessage(msg dns.Message) ([]Result, error) {
	results := make([]Result, 0)
	for _, resource := range msg.Answers {
		h := resource.Header
		hosts := make([]string, 0)
		switch h.Type {
		case dns.TypeA:
			r := resource.Body.(*dns.AResource)
			hosts = append(hosts, net.IP(r.A[:]).String())
		case dns.TypeAAAA:
			r := resource.Body.(*dns.AAAAResource)
			hosts = append(hosts, net.IP(r.AAAA[:]).String())
		case dns.TypeMX:
			r := resource.Body.(*dns.MXResource)
			hosts = append(hosts, r.MX.String())
		case dns.TypeNS:
			r := resource.Body.(*dns.NSResource)
			hosts = append(hosts, r.NS.String())
		case dns.TypeTXT:
			r := resource.Body.(*dns.TXTResource)
			hosts = r.TXT
		case dns.TypeCNAME:
			r := resource.Body.(*dns.CNAMEResource)
			hosts = append(hosts, r.CNAME.String())
		default:
			return nil, fmt.Errorf("unknown query type")
		}
		for _, host := range hosts {
			results = append(results, Result{
				Host:    host,
				Type:    strings.TrimPrefix(h.Type.String(), "Type"),
				TTL:     h.TTL,
				DNSType: h.Type,
			})
		}
	}

	return results, nil
}
