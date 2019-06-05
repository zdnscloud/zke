package set

import (
	"testing"

	ut "github.com/zdnscloud/cement/unittest"
)

func TestStringSetAddDelete(t *testing.T) {
	ss := StringSetFromSlice([]string{"c", "b", "ab", "a"})
	for i := 0; i < 10; i++ {
		ut.Equal(t, ss.ToSortedSlice(), []string{"a", "ab", "b", "c"})
	}

	ut.Assert(t, ss.Member("a"), "")
	ss.Remove("a")
	ut.Assert(t, ss.Member("a") == false, "")
	ss.Add("a")
	ut.Assert(t, ss.Member("a"), "")

	for _, s := range []string{"ab", "c"} {
		ss.Remove(s)
		ut.Assert(t, ss.Member(s) == false, "")
	}
}

func TestStringSetMember(t *testing.T) {
	ss := StringSetFromSlice([]string{"a", "b", "ab", "c"})
	cases := []struct {
		s        string
		isMember bool
	}{
		{"a", true},
		{"b", true},
		{"ab", true},
		{"c", true},
		{"d", false},
		{"cc", false},
		{"aB", false},
	}

	for _, tc := range cases {
		ut.Equal(t, ss.Member(tc.s), tc.isMember)
	}
}

func TestStringSetDifference(t *testing.T) {
	cases := []struct {
		ss1        []string
		ss2        []string
		difference []string
	}{
		{
			[]string{"a", "b", "c"},
			[]string{"a", "b"},
			[]string{"c"},
		},
		{
			[]string{"a"},
			[]string{"a", "b"},
			nil,
		},
		{
			[]string{"a"},
			[]string{"a"},
			nil,
		},

		{
			[]string{"a", "c"},
			[]string{"a", "b"},
			[]string{"c"},
		},
	}

	for _, tc := range cases {
		ss1 := StringSetFromSlice(tc.ss1)
		ss2 := StringSetFromSlice(tc.ss2)
		difference := StringSetFromSlice(tc.difference)
		ut.Assert(t, difference.Equal(ss1.Difference(ss2)), "")
	}
}

func TestStringSetUnion(t *testing.T) {
	cases := []struct {
		ss1   []string
		ss2   []string
		union []string
	}{
		{
			[]string{"a", "b", "c"},
			[]string{"a", "b"},
			[]string{"a", "b"},
		},
		{
			[]string{"a"},
			[]string{"a", "b"},
			[]string{"a"},
		},
		{
			[]string{"a"},
			[]string{"a"},
			[]string{"a"},
		},

		{
			[]string{"c"},
			[]string{"b"},
			nil,
		},
	}

	for _, tc := range cases {
		ss1 := StringSetFromSlice(tc.ss1)
		ss2 := StringSetFromSlice(tc.ss2)
		difference := StringSetFromSlice(tc.union)
		ut.Assert(t, difference.Equal(ss1.Union(ss2)), "")
	}
}
