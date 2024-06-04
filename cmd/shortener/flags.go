package main

import "flag"

func ParseFlags(baseHost *string, shortHost *string) {
	flag.StringVar(baseHost, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(shortHost, "b", "localhost:7777", "base address and port of the resulting shortened URL")

	flag.Parse()
}
