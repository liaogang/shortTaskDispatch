package pin_code_task_cache

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"sync"
	"testing"
)

func TestSyncMapStoreDiffrent(t *testing.T) {
	var m = sync.Map{}

	m.Store("a", errors.New("asdf"))
	m.Store("b", "bVal")

	fmt.Println(m)

	var r = &ReqT{
		Val: "vvvv",
		I:   2,
	}
	m.Store("3", r)

	fmt.Println(m)
	r.Val = "kkkkkkkkkkkkkkk"

	if val, ok := m.Load("3"); ok {
		var req2 = val.(*ReqT)
		fmt.Println(req2)
	}

}

type ReqT struct {
	Val string
	I   int
}
