package domaintree

import (
	ut "github.com/zdnscloud/cement/unittest"
	"github.com/zdnscloud/g53"
	"testing"
)

func stringToName(n string) *g53.Name {
	name, _ := g53.NameFromString(n)
	return name
}

func TestDomainTree(t *testing.T) {
	tree := NewDomainTree()
	tree.Insert(stringToName("a.b.cn."), 1)
	tree.Insert(stringToName("com."), 2)
	tree.Insert(stringToName("b.cn."), 3)

	sum := 0
	tree.ForEach(func(d interface{}) {
		if d != nil {
			sum += d.(int)
		}
	})
	ut.Equal(t, sum, 6)

	name, data, match := tree.Search(stringToName("."))
	ut.Assert(t, match == NotFound, ". has no data")
	ut.Assert(t, name == nil, "root has no exists")

	name, data, match = tree.Search(stringToName("c.b.cn."))
	ut.Assert(t, match == ClosestEncloser, "no c.b.cn exists")
	ut.Equal(t, data.(int), 3)
	ut.Equal(t, name.String(false), "b.cn.")

	name, data, match = tree.Search(stringToName("b.com."))
	ut.Assert(t, match == ClosestEncloser, "no b.com exists")
	ut.Equal(t, name.String(false), "com.")
	ut.Equal(t, data.(int), 2)

	name, data, match = tree.Search(stringToName("a.b.cn."))
	ut.Assert(t, match == ExactMatch, "a.b.cn exists")
	ut.Equal(t, name.String(false), "a.b.cn.")
	ut.Equal(t, data.(int), 1)

	name, data, match = tree.Search(stringToName("a"))
	ut.Assert(t, match == NotFound, "a. not exists")
	ut.Assert(t, data == nil, "no parent encloser exists")
	ut.Assert(t, name == nil, "no parent encloser exists")

	name, data, match = tree.Search(stringToName("a.a.b.cn."))
	ut.Assert(t, match == ClosestEncloser, "a.a.b.cn. not exists")
	ut.Equal(t, name.String(false), "a.b.cn.")
	ut.Equal(t, data.(int), 1)

	tree.Delete(stringToName("a.b.cn"))
	name, data, match = tree.Search(stringToName("a.a.b.cn."))
	ut.Assert(t, match == ClosestEncloser, "a.a.b.cn. not exists")
	ut.Equal(t, name.String(false), "b.cn.")
	ut.Equal(t, data.(int), 3)
}
