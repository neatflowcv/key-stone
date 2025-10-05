package main

import (
	"context"
	"log"
	"net/http"
	"os"

	_ "goa.design/goa/v3/codegen"
	_ "goa.design/goa/v3/codegen/generator"

	tokenserver "github.com/neatflowcv/key-stone/gen/http/token/server"
	userserver "github.com/neatflowcv/key-stone/gen/http/user/server"
	"github.com/neatflowcv/key-stone/gen/token"
	"github.com/neatflowcv/key-stone/gen/user"
	"github.com/neatflowcv/key-stone/internal/app/flow"
	"github.com/neatflowcv/key-stone/internal/pkg/credentialrepository/memory"
	vaultgenerator "github.com/neatflowcv/key-stone/internal/pkg/tokengenerator/vault"
	"github.com/urfave/cli/v3"
	goahttp "goa.design/goa/v3/http"
)

var version = "dev"

func main() {
	log.Println("version", version)

	const (
		flagPort       = "port"
		flagPublicKey  = "public-key"
		flagPrivateKey = "private-key"
	)

	app := &cli.Command{ //nolint:exhaustruct
		Name: "key-stone",
		Flags: []cli.Flag{
			&cli.StringFlag{ //nolint:exhaustruct
				Name:    flagPort,
				Value:   "8080",
				Usage:   "The port to listen on",
				Sources: cli.EnvVars("KS_PORT"),
			},
			&cli.StringFlag{ //nolint:exhaustruct
				Name:     flagPublicKey,
				Sources:  cli.EnvVars("KS_PUBLIC_KEY"),
				Usage:    "The public key to use for the token",
				Required: true,
			},
			&cli.StringFlag{ //nolint:exhaustruct
				Name:     flagPrivateKey,
				Sources:  cli.EnvVars("KS_PRIVATE_KEY"),
				Usage:    "The private key to use for the token",
				Required: true,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			port := c.String(flagPort)
			publicKey := c.String(flagPublicKey)
			privateKey := c.String(flagPrivateKey)

			return startServer(port, publicKey, privateKey)
		},
	}

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func startServer(port, publicKey, privateKey string) error {
	pubVault := vaultgenerator.NewGenerator("key-stone", []byte(publicKey))
	priVault := vaultgenerator.NewGenerator("key-stone", []byte(privateKey))
	repository := memory.NewRepository()

	service := flow.NewService(repository, pubVault, priVault)

	mux := goahttp.NewMuxer()
	requestDecoder := goahttp.RequestDecoder
	responseEncoder := goahttp.ResponseEncoder

	userHandler := NewUserHandler(service)
	userEndpoints := user.NewEndpoints(userHandler)
	userServer := userserver.New(userEndpoints, mux, requestDecoder, responseEncoder, nil, nil)
	userServer.Mount(mux)

	tokenHandler := NewTokenHandler(service)
	tokenEndpoints := token.NewEndpoints(tokenHandler)
	tokenServer := tokenserver.New(tokenEndpoints, mux, requestDecoder, responseEncoder, nil, nil)
	tokenServer.Mount(mux)

	server := &http.Server{Addr: ":" + port, Handler: mux} //nolint:exhaustruct,gosec

	log.Printf("Starting service on :%s", port)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
