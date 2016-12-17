package main

import (
	"net"
	"strings"

	"encoding/hex"
	"sync"

	"fmt"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infoln("tubes")
	socket, err := net.Listen("tcp", ":4196")
	if err != nil {
		panic(err)
	}

	onoff := make(chan string)

	wg := sync.WaitGroup{}
	go func() {
		for {
			conn, err := socket.Accept()
			log.Infoln("incoming connection from", conn.RemoteAddr)
			if err != nil {
				log.Warnln("Error accepting connection:", err)
				continue
			}

			go handleSocket(conn, &wg, onoff)
		}
	}()

	for {
		var txt string
		fmt.Scanln(&txt)

		switch strings.TrimSpace(txt) {
		case "bye":
			return
		case "on":
			onoff <- "on"
		case "off":
			onoff <- "off"
		}
	}
}

func handleSocket(c net.Conn, wg *sync.WaitGroup, onoff chan string) {
	defer c.Close()

	buffer := make([]byte, 128, 128)

	init1 := []byte{0x5a, 0xa5, 0x00, 0x07, 0x02, 0x05, 0x0d, 0x07, 0x05, 0x07, 0x12, 0xc6, 0x5b, 0xb5}
	init2 := []byte{0x5a, 0xa5, 0x00, 0x01, 0x02, 0xfd, 0x5b, 0xb5}
	init3 := []byte{0x5a, 0xa5, 0x00, 0x02, 0x05, 0x01, 0xf9, 0x5b, 0xb5}
	turnOn := []byte{0x5a, 0xa5, 0x00, 0x17, 0x10, 0x01, 0x01, 0x0a, 0xe0, 0x32, 0x23,
		0x4c, 0x90, 0x52, 0xff, 0xfe, 0x00, 0x00, 0x10, 0x11, 0x00, 0x00, 0x01, 0x00, 0x00,
		0x00, 0xff, 0x62, 0x5b, 0xb5}
	turnOff := []byte{0x5a, 0xa5, 0x00, 0x17, 0x10, 0x01, 0x01, 0x0a, 0xe0, 0x32, 0x23, 0x4c,
		0x90, 0x52, 0xff, 0xfe, 0x00, 0x00, 0x10, 0x11, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
		0x00, 0x61, 0x5b, 0xb5}
	wg.Add(1)
	defer wg.Done()

	log.Infoln("Beginning handshake")
	c.Write(init1)
	c.Write(init2)
	readAndLog(c, buffer)

	mac := buffer[10:13]
	log.Infoln("Got mac", hex.EncodeToString(mac))

	c.Write(init3)
	log.Infoln("Completed handshake")
	ch := make(chan int)
	go readLoop(c, buffer, ch)
	for {
		select {
		case msg := <-onoff:
			log.Infoln("Turning it", msg)
			switch strings.TrimSpace(msg) {
			case "on":
				c.Write(turnOn)
			case "off":
				c.Write(turnOff)
			}

		case length := <-ch:
			if length == 16 {
				log.Info("responding")
				c.Write([]byte{0x5a, 0xa5, 0x00, 0x01, 0x06, 0xf9, 0x5b, 0xb5})
			}
		}
	}
}

func readLoop(c net.Conn, buffer []byte, ch chan int) {
	for {
		m, e := readAndLog(c, buffer)
		if e != nil {
			panic(e)
		}
		ch <- m
	}
}

func readAndLog(c net.Conn, buffer []byte) (int, error) {
	n, e := c.Read(buffer)

	if e != nil {
		log.Errorln("Got an error:", e, "and", n, "bytes")
		return 0, e
	}
	log.Infoln("Got", n, "bytes")
	return n, nil
}
