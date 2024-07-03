package main

import "flag"

func ParseFlags(p *params) {
	flag.StringVar(&p.ServerAddress, "a", p.ServerAddress, "address to run server")
	flag.StringVar(&p.BaseAddress, "b", p.BaseAddress, "base address of the resulting shortened URL")
	flag.StringVar(&p.FileStoragePath, "f", p.FileStoragePath, "file storage path")
	flag.StringVar(&p.Database_DSN, "d", p.Database_DSN, "data connection Database")
	flag.Parse()
}
