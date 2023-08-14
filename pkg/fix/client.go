package fix

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FixClient struct {
	mu             sync.Mutex
	hostName       string
	ServerAddr     string
	Timeout        time.Duration
	readBufferSize int
	conn           net.Conn
	reader         *bufio.Reader
}

type FixResponse struct {
	BeginString  string // 8
	BodyLength   string // 9
	SenderCompID string // 49
	TargetCompID string // 56
	TargetSubID  string // 57
	SenderSubID  string // 50
	MsgSeqNum    string // 34
	SendingTime  string // 52
	Checksum     string // 10
	MsgType      string // 35
	Body         map[int]string
}

func NewTCPClient(hostName string, port int, timeout time.Duration, readBufferSize int) *FixClient {
	return &FixClient{
		ServerAddr:     fmt.Sprintf("%s:%d", hostName, port),
		hostName:       hostName,
		Timeout:        timeout,
		readBufferSize: readBufferSize,
		mu:             sync.Mutex{},
	}
}

func (c *FixClient) RoundTrip(message string) ([]*FixResponse, error) {
	message = strings.ReplaceAll(message, "|", "\u0001")
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.Send([]byte(message))
	if err != nil {
		log.Fatalf("send: %+v", err)
		return nil, err
	}
	messages, err := c.Receive()
	if err != nil {
		log.Fatalf("receive: %+v", err)

		return nil, err
	}

	var responses []*FixResponse
	for _, msg := range messages {

		respString := string(msg)
		if !validateChecksum(respString) {
			return nil, errors.New("invalid checksum")
		}
		tmpResponse := newParseFixResponse(msg)
		responses = append(responses, tmpResponse)
	}

	return responses, nil

}

func (c *FixClient) Dial() error {
	tcpAddress, err := net.ResolveTCPAddr("tcp", c.ServerAddr)
	if err != nil {
		return err
	}
	tlsConfig := getTlsConfig(c.hostName)
	dialCtx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	dialer := tls.Dialer{
		Config: tlsConfig,
	}
	conn, err := dialer.DialContext(dialCtx, "tcp", tcpAddress.String())
	cancel()
	if err != nil {
		return err
	}
	c.conn = conn
	c.reader = bufio.NewReader(c.conn)
	return nil
}

func (c *FixClient) Send(data []byte) (int, error) {
	if c.conn == nil {
		return 0, errors.New("connection is nil")
	}
	return c.conn.Write(data)
}

func (c *FixClient) Receive() ([][]byte, error) {
	if c.conn == nil {
		return nil, errors.New("connection is nil")
	}
	var buffer = make([]byte, c.readBufferSize)
	var data []byte
	var messages [][]byte
	// var gotChecksum bool
	for {
		n, err := c.reader.Read(buffer)
		// n, err := c.conn.Read(buffer)
		if err != nil {
			log.Fatalf("read: %+v", err)

			return nil, err
		}
		data = append(data, buffer[:n]...)
		//start of new stuff
		startIndex := 0
		for {
			if len(data[startIndex:]) <= 7 {
				//10=xxx| is trailer, if less than 7 then have hit end of message
				break
			}
			index := bytes.Index(data[startIndex:], []byte("\u000110="))
			if index >= 0 {
				endIndex := index + startIndex + 7 // + 8 to account for actual trailer
				if endIndex < len(data) && string(data[endIndex]) == "\u0001" {
					messages = append(messages, data[startIndex:endIndex+1])
					startIndex = endIndex + 1

				} else {
					break
				}
			} else {
				break
			}

		}

		if n == c.readBufferSize {
			data = data[startIndex:] //so we don't lose incomplete message in the even that the entire buffer is full and there is still data to be read
			continue
		}
		return messages, nil

		//end of new stuff

		// 	if len(data) > 7 && bytes.Contains(data[len(data)-8:], []byte("10=")) {
		// 		if bytes.HasSuffix(data, []byte("|")) {
		// 			break
		// 		}
		// 	}
		// }
		// return data, nil
	}
}

func (c *FixClient) Close() error {
	if c.conn == nil {
		return errors.New("connection is nil")
	}
	return c.conn.Close()
}

func getTlsConfig(hostName string) *tls.Config {
	config := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         hostName,
		//RootCAs:            x509.NewCertPool(),
	}
	return config
}

func validateChecksum(message string) bool {
	messageHeaderAndBody := message[:len(message)-7]
	messageHeaderAndBody = strings.ReplaceAll(messageHeaderAndBody, "|", "\u0001")
	lastTag := message[len(message)-7:]
	parts := strings.Split(lastTag, "=")
	if len(parts) != 2 {
		log.Fatal("Invalid checksum format")
		return false
	}
	resChecksumStr := strings.TrimSuffix(parts[1], "\u0001")
	resChecksum, err := strconv.Atoi(resChecksumStr)
	if err != nil {
		log.Fatal(err)
		return false
	}
	calcChecksum := calculateChecksum(messageHeaderAndBody)
	return calcChecksum == resChecksum
}

func newParseFixResponse(data []byte) *FixResponse {
	dataString := string(data)

	pairs := strings.Split(dataString, "\u0001")
	response := &FixResponse{
		Body: make(map[int]string),
	}
	for _, pair := range pairs {
		tagValue := strings.SplitN(pair, "=", 2)
		if len(tagValue) != 2 {
			continue
		}
		tag, _ := strconv.Atoi(tagValue[0])
		value := tagValue[1]
		switch tag {
		case 8:
			response.BeginString = value
		case 9:
			response.BodyLength = value
		case 49:
			response.SenderCompID = value
		case 56:
			response.TargetCompID = value
		case 57:
			response.TargetSubID = value
		case 50:
			response.SenderSubID = value
		case 34:
			response.MsgSeqNum = value
		case 52:
			response.SendingTime = value
		case 10:
			response.Checksum = value
		case 35:
			response.MsgType = value
		default:
			response.Body[tag] = value
		}
	}
	return response
}
