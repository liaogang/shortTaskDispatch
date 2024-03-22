package http_server

import (
	"fmt"
	"net"
	"net/http"
	"testing"
)

type THandler struct {
}

func (slf THandler) ServeHTTP(_ http.ResponseWriter, req *http.Request) {

	fmt.Println("wait")

	<-req.Context().Done()
	fmt.Println("done")

}

func TestServerCheckHttpRequestCancel(t *testing.T) {

	l, _ := net.Listen("tcp", ":4040")
	http.Serve(l, &THandler{})

}
