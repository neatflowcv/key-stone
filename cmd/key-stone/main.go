package main

import (
	"log"
	"net/http"

	tokenserver "github.com/neatflowcv/key-stone/gen/http/token/server"
	userserver "github.com/neatflowcv/key-stone/gen/http/user/server"
	"github.com/neatflowcv/key-stone/gen/token"
	"github.com/neatflowcv/key-stone/gen/user"
	goahttp "goa.design/goa/v3/http"
)

var version = "dev"

func main() {
	log.Println("version", version)

	service := NewHandler()
	userEndpoints := user.NewEndpoints(service)
	tokenEndpoints := token.NewEndpoints(service)

	mux := goahttp.NewMuxer()
	requestDecoder := goahttp.RequestDecoder
	responseEncoder := goahttp.ResponseEncoder

	userHandler := userserver.New(userEndpoints, mux, requestDecoder, responseEncoder, nil, nil)
	tokenHandler := tokenserver.New(tokenEndpoints, mux, requestDecoder, responseEncoder, nil, nil)

	userserver.Mount(mux, userHandler)
	tokenserver.Mount(mux, tokenHandler)

	port := "8080"
	server := &http.Server{Addr: ":" + port, Handler: mux} //nolint:exhaustruct,gosec

	log.Printf("Starting service on :%s", port)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
