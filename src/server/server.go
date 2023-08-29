package server

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/RichardKnop/machinery/v1/log"
)

var DataToSend [][]byte

func (options *Options) Serve() {
	DataToSend = readAndSplitFile(options.File, options.Size, options.QueryType == "TXT") // if it's TXT we encode to base64
	if len(DataToSend) == 0 {
		log.FATAL.Fatalln("Something went wrong, no data to send. Exiting")
	}
	log.INFO.Println(strconv.Itoa(len(DataToSend)) + " parts to send")

	// Listen for incoming connections.
	addr := net.UDPAddr{
		Port: 53,
		IP:   net.ParseIP(options.ListenAddress),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer conn.Close()
	fmt.Printf("Listening on udp %s:53\n", options.ListenAddress)
	for {
		buf := make([]byte, 1024)
		rlen, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
		}
		log.INFO.Println("Incomming request from " + remoteAddr.String())
		// Handle connections in a new goroutine.
		go func() {
			bytesToSend := options.handleRequest(buf[0:rlen])
			time.Sleep(time.Duration(options.Wait) * time.Millisecond)
			_, _, err := conn.WriteMsgUDP(bytesToSend, []byte{}, remoteAddr)
			if err == nil {
				log.INFO.Println("Packet sent successfully")
			}

		}()
	}

}
