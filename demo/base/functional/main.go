package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func (f HandlerFunc) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

func main() {
	hf := HandlerFunc(HelloHandler)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", bytes.NewBuffer([]byte("test")))

	hf.ServerHTTP(resp, req)

	bts, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bts))
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello world"))
}
