package main

import "testing"

func TestMatchDN(t *testing.T) {
	cases := []struct {
		dn   string
		want bool
	}{
		{"www.google.com", true},
		{"google.com", true},
		{"mkyong123.com", true},
		{"mkyong-info.com", true},
		{"sub.mkyong.com", true},
		{"sub.mkyong-info.com", true},
		{"mkyong.com.au", true},
		{"g.co", true},
		{"mkyong.t.t.co", true},
		{"mkyong.t.t.c", false},
		{"mkyong,com", false},
		{"mkyong", false},
		{"mkyong.123", false},
		{".com", false},
		{"mkyong.com/users", false},
		{"-mkyong.com", false},
		{"mkyong-.com", false},
		{"sub.-mkyong.com", false},
		{"sub.mkyong-.com", false},
		{"", false},
	}
	for _, c := range cases {
		got := matchDN(c.dn)
		if got != c.want {
			t.Errorf("matchDN(%q) == %v, want %v", c.dn, got, c.want)
		}
	}
}
