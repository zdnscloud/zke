package slice

import (
	ut "github.com/zdnscloud/cement/unittest"
	"testing"
)

func TestSliceDifference(t *testing.T) {
	s1 := []string{"a", "b", "a", "c"}
	s2 := []string{"e", "b", "b"}
	s3 := SliceDifference(s1, s2)
	if len(s3) != 2 {
		t.Errorf("s1 has interleave with s2:%v", s3)
	}

	s1 = []string{"a", "b", "a", "c"}
	s2 = []string{"e"}
	s3 = SliceDifference(s1, s2)
	if len(s3) != 3 {
		t.Errorf("s1 has no interleave with s2:%v", s3)
	}

	s1 = []string{"a", "b", "a", "c"}
	s2 = []string{"a", "c", "c"}
	s3 = SliceDifference(s1, s2)
	if len(s3) != 1 {
		t.Errorf("s1 includes s2:%v", s3)
	}

	s1 = []string{"a", "b", "a", "c"}
	s2 = []string{"e"}
	s3 = SliceIntersection(s1, s2)
	if len(s3) != 0 {
		t.Errorf("s1 has no interleave with s2:%v", s3)
	}

	s1 = []string{"a", "b", "a", "c"}
	s2 = []string{"a", "c", "c"}
	s3 = SliceIntersection(s1, s2)
	if len(s3) != 2 {
		t.Errorf("s1 has two elements same with s2:%v", s3)
	}
}

func TestSliceIndexAndRemove(t *testing.T) {
	s1 := []string{"a", "b", "c", "d"}
	if SliceIndex(s1, "b") != 1 {
		t.Errorf("s1 has b but index doesn't find it")
	}

	if SliceIndex(s1, "e") != -1 {
		t.Errorf("s1 has no  e but index find it")
	}

	s1 = SliceRemoveAt(s1, 3)
	if len(s1) != 3 || s1[2] != "c" {
		t.Errorf("remove last one in slice failed")
	}

	s1 = SliceRemoveAt(s1, 0)
	if len(s1) != 2 || s1[0] != "b" {
		t.Errorf("remove first one in slice failed")
	}

	s1 = SliceRemoveAt(s1, 2)
	if len(s1) != 2 {
		t.Errorf("remove index out of range should return itself")
	}

	s1 = SliceRemove(s1, "b")
	if len(s1) != 1 {
		t.Errorf("remove last b should succeed")
	}
}

func TestRandElem(t *testing.T) {
	s1 := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	conflictTime := 0
	for i := 0; i < 1000; i++ {
		if RandElem(s1) == RandElem(s1) {
			conflictTime += 1
		}
	}
	ut.Assert(t, conflictTime < 120, "conflict time should smaller than 120 %d\n", conflictTime)
}

func TestShuffleSlice(t *testing.T) {
	s1 := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	for i := 0; i < 100; i++ {
		clone := make([]string, len(s1))
		copy(clone, s1)
		Shuffle(s1)
		ut.NotEqual(t, clone, s1)
	}
}

func TestSliceClone(t *testing.T) {
	s1 := []string{"a", "b", "c", "d"}
	s2 := Clone(s1)
	ut.Equal(t, s1, s2)
	ut.Equal(t, len(s1), len(s2))

	s2 = Clone(nil)
	ut.Equal(t, len(s2), 0)
}
