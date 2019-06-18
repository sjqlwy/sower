package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

var version, date string

func main() {
	buf := make([]byte, 64)
	copy(buf, "123")
	c1, c2 := net.Pipe()
	go c1.Write(buf)

	buf2 := make([]byte, 64)
	if _, err := io.ReadFull(c2, buf2); err != nil {
		fmt.Println(err)
	}

	buf2 = strings.TrimRight(buf2, string(0))
	fmt.Println(len(buf2), buf2, buf2)
}
