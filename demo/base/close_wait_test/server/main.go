package main

import (
	"fmt"
	"net"
)

func main() {
	//1.监听端口
	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		fmt.Printf("listen fail ,err: %v\n", err)
		return
	}

	//2.建立套接字链接
	for true {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept fail ,err: %v\n", err)
			continue
		}

		//3.创建协程
		go func(conn net.Conn) {
			defer conn.Close() //协程退出一定要关闭资源,否则会一直占用

			for true {
				var buf [128]byte
				n, err := conn.Read(buf[:])
				if err != nil {
					fmt.Printf("read from connect failed, err: %v\n", err)
					break
				}

				str := string(buf[:n])
				fmt.Printf("receive from client, data: %v\n", str)
			}
		}(conn)
	}
}
