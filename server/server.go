package server

import (
	"crypto/tls"
	"net"
	"bufio"
	"log"
	"crypto/x509"
	"io/ioutil"
)

type Server struct {
	CertificatePath string
	KeyPath         string
}

func (server Server) Run(){
	cer, err := tls.LoadX509KeyPair(server.CertificatePath, server.KeyPath)
	if err != nil {
		log.Println(err)
		return
	}

	certPool := x509.NewCertPool()
	data, err := ioutil.ReadFile("resources/server/server.crt")

	if err != nil {
		log.Println( err )
	}
	certPool.AppendCertsFromPEM(data)

	config := &tls.Config{
		ClientCAs: certPool,
		Certificates: []tls.Certificate{cer},
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	ln, err := tls.Listen("tcp", ":443", config)

	if err != nil {
		log.Println(err)
		return
	}

	defer ln.Close()

	for {
		println("Waiting for tls connection")
		conn, err := ln.Accept()
		if err != nil{
			log.Println(err)
			continue
		}
		go server.HandleConnection(conn)
	}
}

func (s Server) HandleConnection(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)

	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		println(msg)

		_, err = conn.Write([]byte("World\n"))

		if err != nil {
			log.Println(err)
			return
		}
	}
}

