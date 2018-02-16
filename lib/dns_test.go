package lib

import (
	"fmt"
	"testing"
)

func TestQuery(t *testing.T) {
	domain := "kfd.me"
	ans := Question(domain, "A")
	if ans.Error() != nil {
		fmt.Println(ans.Error())
		return
	}
	fmt.Printf("%v\n", ans)
}
