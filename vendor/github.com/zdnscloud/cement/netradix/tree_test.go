package netradix

import (
	"testing"

	ut "github.com/zdnscloud/cement/unittest"
)

func TestTree(t *testing.T) {
	tr := NewTree(0)
	if tr == nil || tr.root == nil {
		t.Error("Did not create tree properly")
	}
	err := tr.AddCIDR("1.2.3.0/25", 1)
	ut.Assert(t, err == nil, "add cidr failed")

	testAddrs := []string{"1.2.3.1/25", "1.2.3.60/32", "1.2.3.60", "1.2.3.160/32", "1.2.3.160", "1.2.3.128/25", "1.2.3.0/24"}
	expectResult := []interface{}{1, 1, 1, nil, nil, nil, nil}

	for i, addr := range testAddrs {
		inf, err := tr.FindCIDR(addr)
		ut.Assert(t, err == nil, "add cidr failed")
		ut.Equal(t, inf, expectResult[i])
	}

	// Covering defined
	err = tr.AddCIDR("1.2.3.0/24", 2)
	ut.Assert(t, err == nil, "add cidr failed")

	testAddrs = []string{"1.2.3.0/24", "1.2.3.0/25", "1.2.3.160/32"}
	expectResult = []interface{}{2, 1, 2}

	for i, addr := range testAddrs {
		inf, err := tr.FindCIDR(addr)
		ut.Assert(t, err == nil, "add cidr failed")
		ut.Equal(t, inf, expectResult[i])
	}

	// Delete internal
	err = tr.DeleteCIDR("1.2.3.0/25")
	ut.Assert(t, err == nil, "delete cidr failed")

	// Hit covering with old IP
	inf, err := tr.FindCIDR("1.2.3.0/32")
	ut.Assert(t, err == nil, "add cidr failed")
	ut.Equal(t, inf.(int), 2)

}

func TestTree6(t *testing.T) {
	tr := NewTree(0)
	ut.Assert(t, tr != nil && tr.root != nil, "new tree failed")

	err := tr.AddCIDR("dead::0/16", 3)
	ut.Assert(t, err == nil, "add cidr failed")

	// Matching defined cidr
	inf, err := tr.FindCIDR("dead::beef")
	ut.Assert(t, err == nil, "find cidr failed")
	ut.Equal(t, inf.(int), 3)

	// Outside
	inf, err = tr.FindCIDR("deed::beef/32")
	ut.Assert(t, err == nil, "add cidr failed")
	ut.Equal(t, inf, nil)

	err = tr.AddCIDR("dead:beef::0/48", 4)
	ut.Assert(t, err == nil, "add cidr failed")

	// Match defined subnet
	inf, err = tr.FindCIDR("dead:beef::0a5c:0/64")
	ut.Assert(t, err == nil, "find cidr failed")
	ut.Equal(t, inf.(int), 4)

	// Match outside defined subnet
	inf, err = tr.FindCIDR("dead:0::beef:0a5c:0/64")
	ut.Assert(t, err == nil, "find cidr failed")
	ut.Equal(t, inf.(int), 3)
}

func TestRegression6(t *testing.T) {
	tr := NewTree(0)
	if tr == nil || tr.root == nil {
		t.Error("Did not create tree properly")
	}
	// in one of the implementations /128 addresses were causing panic...
	tr.AddCIDR("2620:10f::/32", 54321)
	tr.AddCIDR("2620:10f:d000:100::5/128", 12345)

	inf, err := tr.FindCIDR("2620:10f:d000:100::5/128")
	if err != nil {
		t.Errorf("Could not get /128 address from the tree, error: %s", err)
	} else if inf.(int) != 12345 {
		t.Errorf("Wrong value from /128 test, got %d, expected 12345", inf)
	}
}

func TestV4UniversalSubnet(t *testing.T) {
	tr := NewTree(0)
	err := tr.AddCIDR("0.0.0.0/0", 1)
	ut.Assert(t, err == nil, "add cidr failed")
	err = tr.AddCIDR("1.2.3.4/8", 2)
	ut.Assert(t, err == nil, "add cidr failed")

	testAddrs := []string{"2.2.3.1", "3.2.3.60", "1.2.3.60", "1.2.3.160/16", "8.8.8.8"}
	expectResult := []interface{}{1, 1, 2, 2, 1}
	for i, addr := range testAddrs {
		inf, err := tr.FindCIDR(addr)
		ut.Assert(t, err == nil, "add cidr failed")
		ut.Equal(t, inf, expectResult[i])
	}
}
