package main

import (
	"flag"
	"fmt"
)

func main() {
	address := flag.String("address", "localhost:8080", "--address=:8080 or --address=192.168.1.10:80")
	flag.Parse()
	fmt.Printf("http://%s\r\n", *address)
	RunPage(*address)
}
