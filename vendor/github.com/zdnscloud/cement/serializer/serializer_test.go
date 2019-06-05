package serializer

import (
	ut "github.com/zdnscloud/cement/unittest"
	"testing"
)

type student struct {
	Name         string
	Age          int
	LectureNames []string
	Scores       []int
}

func TestEncodeAndDecode(t *testing.T) {
	s := student{"ben", 30, []string{"math", "physics"}, []int{89, 90}}
	serializer := NewSerializer()
	err := serializer.Register(s)
	ut.Assert(t, err != nil, "register student will failed")

	err = serializer.Register(&student{})
	ut.Assert(t, err == nil, "register student pointer should ok")

	j, err := serializer.Encode(&s)
	ut.Assert(t, err == nil, "encode student should succeed but get %v", err)

	var s2 student
	err = serializer.Fill(j, &s2)
	ut.Assert(t, err == nil, "decode student should succeed but get %v", err)
	ut.Equal(t, s2, s)

	i3, err := serializer.Decode(j)
	ut.Assert(t, err == nil, "decode student should succeed but get %v", err)
	s3, _ := i3.(*student)
	ut.Equal(t, *s3, s)

	i4, err := serializer.Decode(j)
	ut.Assert(t, err == nil, "decode student should succeed but get %v", err)
	s4, _ := i4.(*student)
	ut.Equal(t, *s4, s)

	s3.Name = "xxxx"
	ut.Equal(t, *s4, s)
}
