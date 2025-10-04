package main

import (
	"log"
	"net/http"

	tokenserver "github.com/neatflowcv/key-stone/gen/http/token/server"
	userserver "github.com/neatflowcv/key-stone/gen/http/user/server"
	"github.com/neatflowcv/key-stone/gen/token"
	"github.com/neatflowcv/key-stone/gen/user"
	goahttp "goa.design/goa/v3/http"

	_ "goa.design/goa/v3/codegen"
	_ "goa.design/goa/v3/codegen/generator"

	"github.com/neatflowcv/key-stone/internal/app/flow"
	"github.com/neatflowcv/key-stone/internal/pkg/credentialrepository/fake"
	vaultgenerator "github.com/neatflowcv/key-stone/internal/pkg/tokengenerator/vault"
)

var version = "dev"

func main() {
	log.Println("version", version)

	pubVault := vaultgenerator.NewGenerator("key-stone", []byte("public-key"))
	priVault := vaultgenerator.NewGenerator("key-stone", []byte("private-key"))
	repository := fake.NewRepository()

	service := flow.NewService(repository, pubVault, priVault)

	handler := NewHandler(service)
	userEndpoints := user.NewEndpoints(handler)
	tokenEndpoints := token.NewEndpoints(handler)

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
