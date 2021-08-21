package real_server_register

import (
	"GateWayDemoStudent/proxy/zookeeper"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type RealServer struct {
	Addr string
}

func main() {
	rs1 := &RealServer{
		Addr: "127.0.0.1:2003",
	}
	rs1.Run()

	time.Sleep(time.Second * 2)

	//监听关闭信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func (r *RealServer) Run() {
	log.Println("Starting http server at: " + r.Addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", r.HelloHandler)
	mux.HandleFunc("/base/error", r.ErrorHandler)

	server := &http.Server{
		Addr:         r.Addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}

	go func() {
		zkManager := zookeeper.NewZkManager([]string{"127.0.0.1:2181"})
		err := zkManager.GetConnect()
		if err != nil {
			fmt.Printf("connect zk error: %s", err)
		}

		defer zkManager.Close()
		err = zkManager.RegisterServerPath("real_server", r.Addr)
		if err != nil {
			fmt.Printf("register node err: %s", err)
		}

		zList, err := zkManager.GetServerListByPath("real_server")
		fmt.Println(zList)

		log.Fatal(server.ListenAndServe())
	}()
}

func (r RealServer) HelloHandler(w http.ResponseWriter, req *http.Request) {
	uPath := fmt.Sprintf("http://%s%s\n", r.Addr, req.URL.Path)

	_, err := io.WriteString(w, uPath)
	if err != nil {
		return
	}
}

func (r RealServer) ErrorHandler(w http.ResponseWriter, req *http.Request) {
	uPath := "error handler"

	w.WriteHeader(500)

	w.WriteHeader(500)

	_, err := io.WriteString(w, uPath)
	if err != nil {
		return
	}
}
