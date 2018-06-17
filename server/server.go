package server

import (
	"crypto/tls"
	"net"
	"bufio"
	"log"
	"crypto/x509"
	"io/ioutil"
	"strconv"
)

type Server struct {
	CertificatePath string
	KeyPath         string
	RootCAPath 		string
	Config 			*tls.Config

}

func (server* Server) Run(){
	err := server.init()
	if err != nil{
		log.Println(err)
	}

	ln, err := tls.Listen("tcp", ":443", server.Config)
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()

	nums := server.CollectXorNumbers(ln)

	log.Println("Calculating result...")
	result := nums[0] ^ nums[1] ^ nums[2]
	log.Printf("The result is %d\n", result)
	message := strconv.Itoa(result) + "\n"

	err = server.sendBackResult(message)
	if err != nil{
		log.Println(err)
	}
}

func (server* Server) sendBackResult(message string) (error) {
	for i := 1; i <= 3; i++ {
		log.Printf("Sending result to cryptographer %d", i)

		conn, err := tls.Dial("tcp", "localhost:808"+strconv.Itoa(i), server.Config)
		if err != nil {
			return err
		}

		defer conn.Close()

		_, err = conn.Write([]byte(message))
		if err != nil {
			return err
		}
	}

	return nil
}

func (server* Server) CollectXorNumbers(ln net.Listener) ([3]int){
	var nums [3]int

	for i := 0; i < 3; i++ {
		log.Println("Waiting for tls connection...")

		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		num, err := server.HandleConnection(conn)
		if err != nil {
			log.Println(err)
			continue
		}

		nums[ i ] = num
	}

	return nums
}

func (server* Server) init() (error) {
	cer, err := tls.LoadX509KeyPair(server.CertificatePath, server.KeyPath)
	if err != nil {
		log.Println(err)
		return err
	}
	certPool := x509.NewCertPool()
	data, err := ioutil.ReadFile("resources/CA/server.crt")
	if err != nil {
		log.Println(err)
		return err
	}
	certPool.AppendCertsFromPEM(data)
	server.Config = &tls.Config{
		RootCAs:      certPool,
		ClientCAs:    certPool,
		Certificates: []tls.Certificate{cer},
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	return nil
}

func (s* Server) HandleConnection(conn net.Conn) (int, error) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	result := 0 

	msg, err := r.ReadString('\n')
	if err != nil {
		return -1, err
	}

	result, err = strconv.Atoi(msg[0:1])
	if err != nil {
		return -1, err
	}

	log.Printf("Received message: %s\n", msg)

	conn.Close()

	return result, nil
}

