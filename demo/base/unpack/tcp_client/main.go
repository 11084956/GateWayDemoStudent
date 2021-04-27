package main

import (
	"GateWayDemoStudent/demo/base/unpack/unpack"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		fmt.Printf("connect failed, err: %v\n", err.Error())
		return
	}

	defer conn.Close()

	_ = unpack.Encode(conn, "hello world 0!!!")
}
