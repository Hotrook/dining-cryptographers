package main

import (
	"log"
	"crypto/tls"
	"io/ioutil"
	"crypto/x509"
)

func main(){

	certPool := x509.NewCertPool()
	data, err := ioutil.ReadFile("resources/server/server.crt")

	if err != nil {
		log.Println( err )
	}
	clientCert, err := tls.LoadX509KeyPair("resources/clients/crts/client1.crt", "resources/clients/keys/client1.key")

	certPool.AppendCertsFromPEM(data)

	conf := &tls.Config{
		//InsecureSkipVerify: true,
		RootCAs:certPool,
		Certificates: []tls.Certificate{clientCert},
	}

	conn, err := tls.Dial("tcp", "localhost:443", conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	n, err := conn.Write([]byte("hello\n"))
	if err != nil {
		log.Println(n, err)
		return
	}

	buf := make([]byte, 100)
	n, err = conn.Read(buf)
	if err != nil {
		log.Println(n, err)
		return
	}

	println(string(buf[:n]))
}


