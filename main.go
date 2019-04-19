package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	G := flag.Int("G", 1, "number of gigabytes to allocate")
	stride := flag.Int("stride", 1024, "how far apart to put the 1s to counteract lazy allocation")
	flag.Parse()
	bs := make([]byte, (*G)*1024*1024*1024)
	// Set the memory to non-zero so it will actually be allocated if the allocation is lazy.
	for i := 0; i < len(bs); i += *stride {
		bs[i] = byte(1)
		fmt.Printf("Allocated %d of %d bytes (%.2f%%)\r", i+1, len(bs), float32(i+1)/float32(len(bs))*100.0)
	}
	fmt.Printf("                                                              \r")
	fmt.Printf("Allocated %dGiB (%d bytes)\n", *G, len(bs))
	fmt.Println("Wallowing in the memory. Press ctrl-C to quit.")
	for {
		time.Sleep(time.Second)
	}
}
