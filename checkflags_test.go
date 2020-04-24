package main

import "testing"

func TestChecFlags(t *testing.T) {
	cases := []struct {
		args []string
		want int
	}{
		{[]string{"-c=alterconf.toml"}, cStart},
		{[]string{"-q"}, cClean},
		{[]string{""}, cStart},
		{[]string{"abra", "kadabra"}, cStart},
		{[]string{"--s", "x", "y", "z"}, cStop},
		{[]string{"abrakadabra", "-r"}, cStart},
		{[]string{"-s", "-q"}, cError},
	}
	var flags tFlags
	for _, c := range cases {
		got := checkFlags(&flags, c.args)
		if got != c.want {
			t.Errorf("checkFlags(%v) == %v, want %v", c.args, got, c.want)
		}
	}
}
