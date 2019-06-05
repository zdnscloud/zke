package cache

import (
	"fmt"
	"testing"
	"time"

	ut "github.com/zdnscloud/cement/unittest"
)

const defaultTtl = 2 * time.Second

type student struct {
	id  string
	age int
}

func hashStudent(s Value) Key {
	return HashString(s.(*student).id)
}

func TestCache(t *testing.T) {
	cache := New(3, hashStudent, false)
	ut.Equal(t, cache.Len(), 0)

	s := &student{"s1", 10}
	cache.Add(s, defaultTtl)
	ut.Equal(t, cache.Len(), 1)

	same_s, found := cache.Get(hashStudent(s))
	ut.Assert(t, found == true, "student should be fetched")
	ut.Equal(t, same_s.(*student).age, 10)

	cache.Add(s, defaultTtl)
	ut.Equal(t, cache.Len(), 1)

	cache.Add(&student{"s2", 20}, 2*defaultTtl)
	cache.Add(&student{"s3", 30}, 2*defaultTtl)
	ut.Equal(t, cache.Len(), 3)

	cache.Add(&student{"s3", 30}, 2*defaultTtl)
	ut.Equal(t, cache.Len(), 3)

	<-time.After(defaultTtl)
	_, found = cache.Get(hashStudent(s))
	ut.Assert(t, found == false, "s1 should expired")
	ut.Equal(t, cache.Len(), 3)

	s2, found := cache.Get(HashString("s2"))
	ut.Assert(t, found == true, "s2 should exists")
	ut.Equal(t, s2.(*student).age, 20)
	cache.Remove(HashString("s2"))
	_, found = cache.Get(HashString("s2"))
	ut.Assert(t, found == false, "s2 is removed")
	ut.Equal(t, cache.Len(), 2)

	for i := 10; i < 100; i++ {
		name := fmt.Sprintf("s%d", i)
		cache.Add(&student{name, i}, 2*defaultTtl)
	}
	ut.Equal(t, cache.Len(), 3)
	for i := 97; i < 100; i++ {
		name := fmt.Sprintf("s%d", i)
		s, found := cache.Get(HashString(name))
		ut.Assert(t, found == true, "last add is in the front")
		ut.Equal(t, s.(*student).age, i)
	}
}
