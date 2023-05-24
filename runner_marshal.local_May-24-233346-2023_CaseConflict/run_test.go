package Runner

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	r := Runner{}
	err := r.GetLatestHash()
	if err != nil {

	}
	fmt.Println(r.Data)
}
