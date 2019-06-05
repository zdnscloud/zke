package errgroup

import (
	"fmt"
	"sort"
	"testing"
	"time"

	ut "github.com/zdnscloud/cement/unittest"
)

const SleepTime = 1

func dumpAppend(param interface{}) (interface{}, error) {
	<-time.After(time.Duration(SleepTime) * time.Second)

	return param.(string) + "good", nil
}

func TestBatch(t *testing.T) {
	var urls = []string{"s1", "s3", "s2"}
	start := time.Now()
	resultCh, err := Batch(urls, dumpAppend)
	ut.Assert(t, time.Since(start).Seconds() < float64(len(urls)*SleepTime), "task should execute concurrently")
	ut.Assert(t, err == nil, "no error should occur")

	var results []string
	for result := range resultCh {
		results = append(results, result.(string))
	}
	sort.Sort(sort.StringSlice(results))
	ut.Equal(t, []string{"s1good", "s2good", "s3good"}, results)
}

func dumpAppendError(param interface{}) (interface{}, error) {
	str := param.(string)
	if str == "s1" {
		return nil, fmt.Errorf("bad")
	}
	return str + "good", nil
}

func TestBatchWithError(t *testing.T) {
	var urls = []string{"s1", "s3", "s2"}
	resultCh, err := Batch(urls, dumpAppendError)
	ut.Assert(t, err != nil, "error should occur")
	ut.Equal(t, err.Error(), "bad")

	var results []string
	for result := range resultCh {
		results = append(results, result.(string))
	}
	sort.Sort(sort.StringSlice(results))
	ut.Equal(t, []string{"s2good", "s3good"}, results)
}

func TestBatchWithNoneSlice(t *testing.T) {
	_, err := Batch(10, dumpAppendError)
	ut.Equal(t, err, ErrInvalidParameter)
}
