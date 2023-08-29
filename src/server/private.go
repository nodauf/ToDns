package server

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/RichardKnop/machinery/v1/log"
	"golang.org/x/net/dns/dnsmessage"
)

func readAndSplitFile(file string, size int, base64Encoded bool) [][]byte {
	var data [][]byte
	fileOutput, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("File reading error", err)
		return data
	}
	// For TXT query we send text, otherwise it is encoded in the IP
	if base64Encoded {
		fileOutput = []byte(base64.StdEncoding.EncodeToString(fileOutput))
	}

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

func (options *Options) handleRequest(buf []byte) []byte {
	var err error
	var dnsParser dnsmessage.Parser
	header, err := dnsParser.Start(buf)
	check(err)
	q, err := dnsParser.Question()
	check(err)
	domainRequested := q.Name.String()
	log.INFO.Println("DNS Request: " + domainRequested)

	// if the format is <id>.d.<domain> it is for download
	if strings.Split(domainRequested, ".")[1] == "d" {
		return options.download(buf, header, q, domainRequested)
		// if the format is <id>.u.<domain> it is for download
	} else if strings.Split(domainRequested, ".")[1] == "u" {
		// Upload not implemented
		return []byte{}
		//return upload(buf, header, q, domainRequested)
	}
	log.ERROR.Println("Unknow action")
	return []byte{}

}

func (options *Options) download(buf []byte, header dnsmessage.Header, q dnsmessage.Question, domainRequested string) []byte {
	if options.QueryType == "TXT" {
		if q.Type.String() != "TypeTXT" {
			log.ERROR.Println("Only answering for TXT DNS query to download")
			return []byte{}
		}

		idRequested, err := strconv.Atoi(strings.Split(domainRequested, ".")[0])
		check(err)
		log.INFO.Println("Requesting chunck " + strconv.Itoa(idRequested+1) + "/" + strconv.Itoa(len(DataToSend)))

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
	} else if options.QueryType == "A" {
		if q.Type.String() != "TypeA" {
			log.ERROR.Println("Only answering for A DNS query to download")
			return []byte{}
		}

		idRequested, err := strconv.Atoi(strings.Split(domainRequested, ".")[0])
		check(err)
		log.INFO.Println("Requesting chunck " + strconv.Itoa(idRequested+1) + "/" + strconv.Itoa(len(DataToSend)))

		if idRequested > len(DataToSend)-1 {
			log.ERROR.Println("ID requested is above the size of the array to send")
			return []byte{}
		}

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

		var dataToReturned [4]byte
		var numberOfIPsReturned int
		var AAnswer dnsmessage.AResource
		if len(DataToSend[idRequested])%4 != 0 {
			numberOfIPsReturned = len(DataToSend[idRequested])/4 + 1
		} else {
			numberOfIPsReturned = len(DataToSend[idRequested]) / 4
		}
		for i := 0; i < numberOfIPsReturned; i++ {
			copy(dataToReturned[:], DataToSend[idRequested][i*4:i*4+4])
			AAnswer.A = dataToReturned
			log.INFO.Printf("Payload: %+v", dataToReturned)
			err = buildAnswer.AResource(resourceHeader, AAnswer)
			check(err)
		}
		bytesToSend, err := buildAnswer.Finish()
		check(err)
		return bytesToSend
	} else {
		log.ERROR.Printf("Query type %s not implemented\n", options.QueryType)
	}
	return []byte{}
}

// Not implemented yet
/*
func upload(buf []byte, header dnsmessage.Header, q dnsmessage.Question, domainRequested string) []byte {
	if q.Type.String() != "TypeCNAME" {
		log.ERROR.Println("Only answering for CNAME DNS query to upload")
		return []byte{}
	}

	idSent, err := strconv.Atoi(strings.Split(domainRequested, ".")[0])
	check(err)
	log.INFO.Println("Chunck sent" + strconv.Itoa(idSent+1))

	/*dataToReturned := DataToSend[idRequested]

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
	return bytesToSend*/
//return []byte{}
//}

func check(err error) {
	if err != nil {
		log.ERROR.Println(err)
	}
}
