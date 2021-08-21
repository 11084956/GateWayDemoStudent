package main

import (
	"GateWayDemoStudent/proxy/zookeeper"
	"fmt"
	"time"
)

func main() {
	zkManger := zookeeper.NewZkManager([]string{"127.0.0.1:2181"})
	err := zkManger.GetConnect()
	if err != nil {
		return
	}

	defer zkManger.Close()

	var i = 0
	for true {
		conf := fmt.Sprintf("{name:" + fmt.Sprint(i) + "}")
		err := zkManger.SetPathData("/rs_server_conf", []byte(conf))
		if err != nil {
			return
		}

		time.Sleep(time.Second * 5)
		i++
	}
}
