package set

import (
	ut "github.com/zdnscloud/cement/unittest"
	"testing"
)

func TestSet(t *testing.T) {
	s := NewSet()
	s.Add("good")
	s.Add("boy")

	ut.Assert(t, s.Contains("good"), "set should contains good")
	ut.Assert(t, s.Contains("boy"), "set should contains boy")
	ut.Assert(t, s.Contains("goood") == false, "set shouldn't contains goood")
	s.Remove("good")
	ut.Assert(t, s.Contains("good") == false, "set shouldn't contains good after remove")
}
