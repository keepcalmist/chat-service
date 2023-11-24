package websocketstream

import (
	"io"
	"net/http"
	"time"

	gorillaws "github.com/gorilla/websocket"
)

type Websocket interface {
	SetWriteDeadline(t time.Time) error
	NextWriter(messageType int) (io.WriteCloser, error)
	WriteMessage(messageType int, data []byte) error
	WriteControl(messageType int, data []byte, deadline time.Time) error

	SetPongHandler(h func(appData string) error)
	PongHandler() func(appData string) error
	SetReadDeadline(t time.Time) error
	NextReader() (messageType int, r io.Reader, err error)

	Close() error
}

type Upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Websocket, error)
}

type upgraderImpl struct {
	upgrader *gorillaws.Upgrader
}

func NewUpgrader(allowOrigins []string, secWsProtocol string) Upgrader {
	originsSet := make(map[string]struct{}, len(allowOrigins))
	for _, o := range allowOrigins {
		originsSet[o] = struct{}{}
	}

	upgrader := &gorillaws.Upgrader{
		HandshakeTimeout: time.Millisecond * 500,
		ReadBufferSize:   1000 * 3,
		Subprotocols:     []string{secWsProtocol},
		CheckOrigin: func(r *http.Request) bool {
			_, ok := originsSet[r.Header.Get("Origin")]
			return ok
		},
		EnableCompression: false,
	}

	return &upgraderImpl{
		upgrader: upgrader,
	}
}

func (u *upgraderImpl) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Websocket, error) {
	return u.upgrader.Upgrade(w, r, responseHeader)
}
