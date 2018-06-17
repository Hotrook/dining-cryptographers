package cryptographer

import (
	"crypto/tls"
	"log"
	"crypto/x509"
	"io/ioutil"
	"strconv"
	"net"
	"bufio"
	"math/rand"
	"time"
)

type Cryptographer struct {
	Id 				int
	RootCAPath 		string
	Config 			*tls.Config
	Payed 			bool
}

func (c * Cryptographer) Run() {
	ln, _ := c.init()
	defer ln.Close()

	rand.Seed(time.Now().UTC().UnixNano())
	firstNumber := rand.Intn(2)

	if c.Id == 1{
		err := c.sendNumber(firstNumber)
		if err != nil {
			log.Println(err)
		}
	}

	secondNumber, err := c.receiveNumber(ln)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Received number: %d\n", secondNumber)

	if c.Id != 1 {
		err := c.sendNumber(firstNumber)
		if err != nil {
			log.Println(err)
		}
	}

	result := firstNumber ^ secondNumber
	if c.Payed {
		result = 1 - result
	}

	c.sendResultToServer(result)
	finalResult, err := c.receiveNumber(ln)

	if finalResult == 0 {
		log.Println("NSA payed for us!")
	} else {
		log.Println("One of us payed")
	}
}

func (c * Cryptographer) sendNumber(firstNumber int) error {
	nextCryptographer := c.Id%3 + 1
	log.Printf("Sending number %d to cryptographer %d\n", firstNumber, nextCryptographer)

	conn, err := tls.Dial("tcp", "localhost:808"+strconv.Itoa(nextCryptographer), c.Config)
	if err != nil {
		return err
	}

	defer conn.Close()

	message := strconv.Itoa(firstNumber) + "\n"
	_, err = conn.Write([]byte(message))
	if err != nil {
		return err
	}

	return nil
}

func (c * Cryptographer) receiveNumber(ln net.Listener) (int, error) {
	log.Println("Waiting for number...")

	conn, err := ln.Accept()
	if err != nil {
		log.Println(err)
		return -1, nil
	}

	num, err := c.HandleConnection(conn)
	if err != nil {
		return -1, err
	}

	return num, nil
}

func (c * Cryptographer) init() (net.Listener, error) {
	certificatePath := "resources/clients/crts/client" + strconv.Itoa(c.Id) + ".crt"
	keyPath := "resources/clients/keys/client" + strconv.Itoa(c.Id) + ".key"
	cer, err := tls.LoadX509KeyPair(certificatePath, keyPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	certPool := x509.NewCertPool()
	data, err := ioutil.ReadFile("resources/CA/server.crt")

	if err != nil {
		log.Println(err)
	}
	certPool.AppendCertsFromPEM(data)
	c.Config = &tls.Config{
		RootCAs:      certPool,
		ClientCAs:    certPool,
		Certificates: []tls.Certificate{cer},
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	ln, err := tls.Listen("tcp", ":808" + strconv.Itoa(c.Id), c.Config)
	if err != nil {
		log.Println(err)
	}

	return ln, err
}

func (c * Cryptographer) HandleConnection(conn net.Conn) (int, error) {
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

	conn.Close()

	return result, nil
}

func (c * Cryptographer) sendResultToServer(result int) {
	log.Printf("Sending result[%d] to server.", result)

	conn, err := tls.Dial("tcp", "localhost:443", c.Config)
	if err != nil {
		log.Println(err)
	}

	defer conn.Close()

	message := strconv.Itoa(result) + "\n"
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Println(err)
	}
}
