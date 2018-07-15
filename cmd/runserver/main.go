package main

import (
	"github.com/Hotrook/dining_cryptographers/server"
	"flag"
	"github.com/Hotrook/dining_cryptographers/logutils"
)

var(
	caCert = flag.String("cert", "resources/server/server.crt", "A PEM encoded certificate file.")
	keyPath = flag.String("key", "resources/server/server.key", "A PEM encoded certificate file.")
)

func main(){

	flag.Parse()
	logutils.InitLogger()

	server := server.Server{CertificatePath: *caCert, KeyPath: *keyPath}
	server.Run()
}