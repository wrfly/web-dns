package dig

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestDig(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	ns := []string{
		"114.114.114.114:53",
		"8.8.8.8:53",
	}
	domain := "kfd.me"
	digger := New("mem", ns)

	t.Run("Type A of kfd.me", func(t *testing.T) {
		ips, err := digger.Dig(domain, "A")
		if err != nil {
			t.Error(err)
		}
		for _, ip := range ips {
			t.Log(ip)
		}
	})
	t.Run("Type MX of kfd.me", func(t *testing.T) {
		ips, err := digger.Dig(domain, "MX")
		if err != nil {
			t.Error(err)
		}
		for _, ip := range ips {
			t.Log(ip)
		}
	})
	t.Run("Type X of kfd.me", func(t *testing.T) {
		_, err := digger.Dig(domain, "X")
		if err.Error() != "type not support" {
			t.Error(err)
		}
	})
}
