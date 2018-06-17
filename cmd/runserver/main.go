package main

import (
	"github.com/Hotrook/dining_cryptographers/server"
	"flag"
)

var(

)


func main(){

	var caCert = flag.String("cert", "someCertFile", "A PEM encoded certificate file.")

	flag.Parse()
	print("Read file: ")
	println(*caCert)
	server := server.Server{CertificatePath: *caCert, KeyPath: "resources/server/server.key"}
	server.Run()

}