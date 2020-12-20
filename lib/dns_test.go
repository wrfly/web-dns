package lib

import (
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	for _, typ := range []string{"A", "AAAA", "MX", "NS", "TXT"} {
		t.Logf("dig kfd.me %s", typ)
		ans := Question(QueryOption{
			NSServer: "8.8.8.8:53",
			Domain:   "kfd.me",
			Type:     typ,
			Timeout:  time.Millisecond * 200,
		})
		if err := ans.Err; err != nil {
			t.Errorf("err: %s", err)
			continue
		}
		for _, ip := range ans.Hosts() {
			t.Logf("IP: %s", ip)
		}
	}

}
