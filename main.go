package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	n := flag.Int("n", 1, "number of gigabytes to allocate")
	flag.Parse()
	bs := make([]byte, (*n)*1024*1024*1024)
	// Set the memory to non-zero so it will actually be allocated.
	// Apparently it's lazily allocated.
	for i := range bs {
		bs[i] = byte(1)
	}
	fmt.Printf("Allocated %dGiB (%d bytes)\n", *n, len(bs))
	fmt.Println("Wallowing in the memory. Press ctrl-C to quit.")
	for {
		time.Sleep(time.Second)
	}
}
