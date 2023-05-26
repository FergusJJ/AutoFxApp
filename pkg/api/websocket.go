package api

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"pollo/config"

	"github.com/fasthttp/websocket"
)

func CreateApiConnection(cid string) (*websocket.Conn, error) {
	apiAddress, err := config.Config("API_ADDRESS")
	if err != nil {
		return nil, err
	}
	if apiAddress == "" {
		return nil, errors.New("api address not specified")
	}
	wsPath, err := config.Config("WS_PATH")
	if err != nil {
		return nil, err
	}
	if wsPath == "" {
		return nil, errors.New("websocket path not specified")
	}
	rawQuery := fmt.Sprintf("id=%s", cid)
	addr := flag.String("addr", apiAddress, "http service address")
	u := url.URL{
		Scheme:   "ws",
		Host:     *addr,
		Path:     wsPath,
		RawQuery: rawQuery,
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (session *ApiSession) ListenForMessages() {
	for {
		messageType, message, err := session.Client.Connection.ReadMessage()
		if err != nil {
			switch {
			case err == io.EOF:
				// The remote peer closed the connection
				log.Println("WebSocket connection closed by remote peer")
				return
			case websocket.IsCloseError(err, websocket.CloseNormalClosure):
				// The connection was closed normally
				log.Println("WebSocket connection closed normally")
				return
			case websocket.IsCloseError(err, websocket.CloseAbnormalClosure):
				// The connection was closed abnormally
				log.Printf("WebSocket connection closed abnormally: %v\n", err)
				// If the connection was closed abnormally, you may want to attempt to reconnect
				return
			case websocket.IsCloseError(err, websocket.CloseGoingAway):
				// The connection is going away
				log.Printf("WebSocket connection is going away: %v\n", err)
				// If the connection is going away, you may want to attempt to reconnect
				return
			case websocket.IsCloseError(err, websocket.CloseNoStatusReceived):
				// The connection was closed without receiving a close status code
				log.Printf("WebSocket connection closed without status code: %v\n", err)
				// If the connection was closed without a status code, you may want to attempt to reconnect
				return
			case websocket.IsCloseError(err, websocket.CloseUnsupportedData):
				// The connection received unsupported data
				log.Printf("WebSocket connection received unsupported data: %v\n", err)
				// If the connection received unsupported data, you may want to attempt to reconnect
				return
			case websocket.IsCloseError(err, websocket.ClosePolicyViolation):
				// The connection violated a policy
				log.Printf("WebSocket connection violated policy: %v\n", err)
				// If the connection violated a policy, you may want to attempt to reconnect
				return
			case websocket.IsCloseError(err, websocket.CloseMessageTooBig):
				// The message received was too big
				log.Printf("WebSocket message received was too big: %v\n", err)
				// If the message received was too big, you may want to attempt to reconnect
				return
			default:
				// Other errors occurred that can be recovered
				log.Printf("WebSocket error occurred: %v\n", err)
				// If the error can be recovered, you may want to attempt to recover and continue the loop
				if errType, ok := err.(*websocket.CloseError); ok {
					log.Printf("WebSocket error type: %s\n", fmt.Sprintf("%d", errType.Code))
				} else {
					log.Printf("WebSocket error type: unknown\n")
				}
				return
				// If the error can be recovered, you may want to attempt to recover and continue the loop
			}
		}
		// messageStr := string(message)
		log.Printf("got message of type: %d\n", messageType)
		session.Client.CurrentMessage <- message
	}

}
