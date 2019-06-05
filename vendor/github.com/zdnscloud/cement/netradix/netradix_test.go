package netradix

import (
	"net"
	"testing"

	"github.com/zdnscloud/cement/slice"
	ut "github.com/zdnscloud/cement/unittest"
)

type TestResult struct {
	ip     string
	udata  string
	result bool
}

type TestData struct {
	network string
	udata   string
}

func TestSearchBest(t *testing.T) {
	rtree := NewNetRadixTree()

	udata, found := rtree.SearchBest(net.ParseIP("::1"))
	ut.Assert(t, found == false, "ipv6 addr isn't in radix tree")
	ut.Equal(t, udata, nil)

	udata, found = rtree.SearchBest(net.ParseIP("1.0.0.0"))
	ut.Assert(t, found == false, "1.0.0.0 isn't in radix tree")
	ut.Equal(t, udata, nil)

	initial := []TestData{
		{"217.72.192.0/20", "UDATA1"},
		{"217.72.195.0/24", "UDATA2"},
		{"195.161.113.74/32", "UDATA3"},
		{"172.16.2.2", "UDATA4"},
		{"10.42.0.0/16", "UDATA5"},
	}

	expected := []TestResult{
		{"217.72.192.1", "UDATA1", true},
		{"217.72.195.42", "UDATA2", true},
		{"195.161.113.74", "UDATA3", true},
		{"172.16.2.2", "UDATA4", true},
		{"15.161.13.75", "", false},
		{"10.42.1.0", "UDATA5", true},
		{"10.42.1.8", "UDATA5", true},
	}

	for _, value := range initial {
		err := rtree.Add(value.network, value.udata)
		ut.Assert(t, err == nil, "add error")
	}

	for _, value := range expected {
		udata, found := rtree.SearchBest(net.ParseIP(value.ip))
		ut.Equal(t, found, value.result)
		if found {
			ut.Equal(t, udata.(string), value.udata)
		}
	}

	rtree.Add("2.2.2.2/24", 10)
	udata, found = rtree.SearchBest(net.ParseIP("2.2.2.3"))
	ut.Assert(t, found, "2.2.2.3 is in subnet 2.2.2.0")
	ut.Equal(t, udata.(int), 10)

	rtree.Delete("2.2.2.2/24")
	udata, found = rtree.SearchBest(net.ParseIP("2.2.2.3"))
	ut.Assert(t, found == false, "2.2.2.3 is not in subnet 2.2.2.0")

	udata, found = rtree.SearchBest(net.ParseIP("172.16.2.2"))
	ut.Assert(t, found, "172.16.2.2 exists")
	rtree.Delete("172.16.2.2")
	udata, found = rtree.SearchBest(net.ParseIP("172.16.2.2"))
	ut.Assert(t, found == false, "172.16.2.2 had been deleted")
}

func TestSearchBestRealView(t *testing.T) {
	viewConf := map[string][]string{
		"BGW":  []string{"100.68.254.0/24"},
		"NXLT": []string{"100.65.0.0/16", "100.67.0.0/16", "42.63.0.0/16", "221.199.64.0/18", "100.72.0.0/16", "100.71.0.0/16", "100.70.254.61/32", "100.70.254.35/32", "100.74.0.0/16", "100.73.0.0/16", "100.75.0.0/16", "100.76.0.0/16", "100.69.0.0/16", "100.68.0.0/16", "100.70.0.0/16"},
		"GDKD": []string{"100.68.0.0/14", "100.72.0.0/12"},
	}

	viewOrders := []string{"GDKD", "BGW", "NXLT"}
	for i := 0; i < 1000; i++ {
		rtree := NewNetRadixTree()
		for _, view := range viewOrders {
			for _, ip := range viewConf[view] {
				err := rtree.Add(ip, view)
				ut.Assert(t, err == nil, "add error")
			}
		}
		udata, found := rtree.SearchBest(net.ParseIP("100.68.254.0"))
		ut.Equal(t, found, true)
		ut.Equal(t, udata, "BGW")
		slice.Shuffle(viewOrders)
	}
}
