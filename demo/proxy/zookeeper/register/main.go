package main

import (
	"GateWayDemoStudent/proxy/zookeeper"
	"fmt"
	"time"
)

func main() {
	zkManager := zookeeper.NewZkManager([]string{"127.0.0.1:2181"})
	err := zkManager.GetConnect()
	if err != nil {
		return
	}

	defer zkManager.Close()
	var i = 0
	for true {
		err := zkManager.RegisterServerPath("/real_server", fmt.Sprint(i))
		if err != nil {
			fmt.Println("zookKeeper Register", err)
			return
		}

		fmt.Println("Register", i)

		time.Sleep(time.Second * 5)
		i++
	}
}
