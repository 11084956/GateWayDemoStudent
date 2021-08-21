package main

import (
	"GateWayDemoStudent/proxy/zookeeper"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//获取zk节点列表
	zkManager := zookeeper.NewZkManager([]string{"127.0.0.1:2181"})
	err := zkManager.GetConnect()
	if err != nil {
		return
	}

	defer zkManager.Close()
	zList, err := zkManager.GetServerListByPath("/real_server")
	fmt.Println("server node:", zList)
	if err != nil {
		log.Println(err)
	}

	//动态监听节点变化
	chanList, chanErr := zkManager.WatchServerListByPath("/real_server")
	go func() {
		for true {
			select {
			case changeErr := <-chanErr:
				fmt.Println("changErr", changeErr)
			case changedList := <-chanList:
				fmt.Println("watch node changed", changedList)
			}
		}
	}()

	//关闭信号监听
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
}
