package main

import (
	"github.com/Hotrook/dining_cryptographers/server"
)

func main(){


	server := server.Server{CertificatePath: "resources/server/server.crt", KeyPath: "resources/server/server.key"}
	server.Run()

}