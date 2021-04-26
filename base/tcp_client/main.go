package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	doSend()
	fmt.Print("doSend over")

	doSend()
	fmt.Print("doSend over")
}

func doSend() {
	//1.连接服务器
	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		fmt.Printf("connect failed, err: %v\n", err.Error())
		return
	}
	defer conn.Close()

	//2.读取命令行输入
	inputReader := bufio.NewReader(os.Stdin)
	for true {
		//3.一直读取知道读取\n(换行符)
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Printf("read from console, err: %v\n", err)
			break
		}

		//4.读取到 Q 时停止
		trimmedInput := strings.TrimSpace(input)
		if trimmedInput == "Q" {
			break
		}

		//回复服务器信息
		_, err = conn.Write([]byte(trimmedInput))
		if err != nil {
			fmt.Printf("write failed, err : %v\n", err)
			break
		}
	}

}
