package main

import "flag"

func ParseFlags(serverHost *string, shortHost *string) {
	flag.StringVar(serverHost, "a", *serverHost, "address to run server")
	flag.StringVar(shortHost, "b", *shortHost, "base address of the resulting shortened URL")
	flag.Parse()
}
