package lib

import (
	"fmt"
	"testing"
)

func printAns(ans Answer) {
	if err := ans.Error(); err != nil {
		println(err)
		return
	}
	for _, ip := range ans.IPs() {
		println(ip)
	}
}
func TestQuery(t *testing.T) {
	dnsserver := "114.114.114.114:53"
	domain := "kfd.me"

	typ := "A"
	t.Run(fmt.Sprintf("%s %s", domain, typ), func(t *testing.T) {
		ans := Question(dnsserver, domain, typ)
		printAns(ans)
	})

	typ = "AAAA"
	t.Run(fmt.Sprintf("%s %s", domain, typ), func(t *testing.T) {
		ans := Question(dnsserver, domain, typ)
		printAns(ans)
	})

	typ = "MX"
	t.Run(fmt.Sprintf("%s %s", domain, typ), func(t *testing.T) {
		ans := Question(dnsserver, domain, typ)
		printAns(ans)
	})

	typ = "TXT"
	t.Run(fmt.Sprintf("%s %s", domain, typ), func(t *testing.T) {
		ans := Question(dnsserver, domain, typ)
		printAns(ans)
	})
}
