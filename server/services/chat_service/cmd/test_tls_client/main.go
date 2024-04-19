package main

import (
	"chat_service/internal/client/auth"
	"chat_service/internal/config"
	"context"
	"log"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()

	cl, err := auth.Connect(slog.Default(),
		cfg.OtherServices.AuthService,
		cfg.OtherServices.Cert.CaPath,
		cfg.OtherServices.Cert.CertPath,
		cfg.OtherServices.Cert.KeyPath,
	)
	if err != nil {
		panic(err)
	}

	user, err := cl.AuthenticateUser(context.Background(), []byte("token"))
	if err != nil {
		panic(err)
	}

	log.Println(user)

}
