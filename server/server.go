package server

import (
	"crypto/tls"
	"net"
	"bufio"
	"log"
)

type Server struct {
	CertificatePath string
	KeyPath         string
}

func (server Server) Run(){
	cer, err := tls.LoadX509KeyPair(server.CertificatePath, server.KeyPath)
	if err != nil {
		println(err)
		log.Fatal(err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", ":443", config)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer ln.Close()

	for {
		println("Waiting for tls connection")
		conn, err := ln.Accept()
		if err != nil{
			log.Fatal(err)
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
			return
		}

		println(msg)

		_, err = conn.Write([]byte("World\n"))

		if err != nil {
			return
		}
	}
}

