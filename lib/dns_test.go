package lib

import (
	"fmt"
	"testing"
	"time"
)

func printAns(ans Answer) {
	if err := ans.Err; err != nil {
		println(err)
		return
	}
	for _, ip := range ans.Hosts() {
		println(ip)
	}
}
func TestQuery(t *testing.T) {
	dnsserver := "114.114.114.114:53"
	domain := "kfd.me"
	timeout := time.Millisecond * 100

	typ := "A"
	t.Run(fmt.Sprintf("%s %s", domain, typ), func(t *testing.T) {
		ans := Question(dnsserver, domain, typ, timeout)
		printAns(ans)
	})

	typ = "AAAA"
	t.Run(fmt.Sprintf("%s %s", domain, typ), func(t *testing.T) {
		ans := Question(dnsserver, domain, typ, timeout)
		printAns(ans)
	})

	typ = "MX"
	t.Run(fmt.Sprintf("%s %s", domain, typ), func(t *testing.T) {
		ans := Question(dnsserver, domain, typ, timeout)
		printAns(ans)
	})

	typ = "TXT"
	t.Run(fmt.Sprintf("%s %s", domain, typ), func(t *testing.T) {
		ans := Question(dnsserver, domain, typ, timeout)
		printAns(ans)
	})
}
