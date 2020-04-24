package main

import "flag"

const (
	cError = iota - 1
	cStart
	cStop
	cVerify
	cClean
	cRead
)

type tFlags struct {
	s bool
	v bool
	q bool
	r bool
	c string
}

func checkFlags(flags *tFlags, args []string) int {
	fs := flag.NewFlagSet("fs", flag.ExitOnError)
	fs.BoolVar(&flags.s, "s", false, "Stopping")
	fs.BoolVar(&flags.v, "v", false, "Verification")
	fs.BoolVar(&flags.q, "q", false, "Cleaning")
	fs.BoolVar(&flags.r, "r", false, "Refresh config")
	fs.StringVar(&flags.c, "c", "config.toml", "configuration file")
	fs.Parse(args)
	switch {
	case flags.s && !flags.v && !flags.q && !flags.r:
		{
			return cStop
		}
	case !flags.s && flags.v && !flags.q && !flags.r:
		{
			return cVerify
		}
	case !flags.s && !flags.v && flags.q && !flags.r:
		{
			return cClean
		}
	case !flags.s && !flags.v && !flags.q && flags.r:
		{
			return cRead
		}
	case !flags.s && !flags.v && !flags.q && !flags.r:
		{
			return cStart
		}
	default:
		{
			return cError
		}
	}

}
