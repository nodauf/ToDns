package server

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/RichardKnop/machinery/v1/log"
	"golang.org/x/net/dns/dnsmessage"
)

func readAndSplitFile(file string, size int) [][]byte {
	var data [][]byte
	fileOutput, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("File reading error", err)
		return data
	}
	fileOutput = []byte(base64.StdEncoding.EncodeToString(fileOutput))

	return split(fileOutput, size)

}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf)
	}
	return chunks
}

func handleRequest(buf []byte) []byte {
	var err error
	var dnsParser dnsmessage.Parser
	header, err := dnsParser.Start(buf)
	check(err)
	q, err := dnsParser.Question()
	check(err)
	domainRequested := q.Name.String()
	log.INFO.Println("DNS Request: " + domainRequested)

	if strings.Split(domainRequested, ".")[1] == "d" {
		return download(buf, header, q, domainRequested)
	}
	log.ERROR.Println("Unknow action")
	return []byte{}

}

func download(buf []byte, header dnsmessage.Header, q dnsmessage.Question, domainRequested string) []byte {
	if q.Type.String() != "TypeTXT" {
		log.ERROR.Println("Only answering for TXT DNS query")
		return []byte{}
	}

	idRequested, err := strconv.Atoi(strings.Split(domainRequested, ".")[0])
	check(err)
	log.INFO.Println("Requesting chunck " + strconv.Itoa(idRequested) + "/" + strconv.Itoa(len(DataToSend)-1))

	if idRequested > len(DataToSend)-1 {
		log.ERROR.Println("ID requested is above the size of the array to send")
		return []byte{}
	}
	dataToReturned := DataToSend[idRequested]

	// Response

	headerAnswer := header
	headerAnswer.Response = true
	buildAnswer := dnsmessage.NewBuilder([]byte{}, headerAnswer)
	err = buildAnswer.StartQuestions()
	check(err)
	err = buildAnswer.Question(q)
	check(err)
	err = buildAnswer.StartAnswers()
	check(err)
	resourceHeader := dnsmessage.ResourceHeader{}
	resourceHeader.Name = q.Name
	resourceHeader.Class = q.Class

	TXTAnswer := dnsmessage.TXTResource{}
	TXTAnswer.TXT = append(TXTAnswer.TXT, string(dataToReturned))
	log.INFO.Println("Payload: " + string(dataToReturned))
	err = buildAnswer.TXTResource(resourceHeader, TXTAnswer)
	check(err)
	bytesToSend, err := buildAnswer.Finish()
	check(err)
	return bytesToSend
}

func check(err error) {
	if err != nil {
		log.ERROR.Println(err)
	}
}
