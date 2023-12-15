package websocketstream

import (
	"context"
	"fmt"
	"io"
	"time"

	gorillaws "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/keepcalmist/chat-service/internal/middlewares"
	eventstream "github.com/keepcalmist/chat-service/internal/services/event-stream"
	"github.com/keepcalmist/chat-service/internal/types"
)

const (
	writeTimeout = time.Second
	pongWait     = 10 * time.Second
)

type eventStream interface {
	Subscribe(ctx context.Context, userID types.UserID) (<-chan eventstream.Event, error)
}

//go:generate options-gen -out-filename=handler_options.gen.go -from-struct=Options
type Options struct {
	pingPeriod time.Duration `default:"3s" validate:"omitempty,min=100ms,max=30s"`

	logger       *zap.Logger     `option:"mandatory" validate:"required"`
	eventStream  eventStream     `option:"mandatory" validate:"required"`
	eventAdapter EventAdapter    `option:"mandatory" validate:"required"`
	eventWriter  EventWriter     `option:"mandatory" validate:"required"`
	upgrader     Upgrader        `option:"mandatory" validate:"required"`
	shutdownCh   <-chan struct{} `option:"mandatory" validate:"required"`
}

type HTTPHandler struct {
	Options
}

func NewHTTPHandler(opts Options) (*HTTPHandler, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return &HTTPHandler{Options: opts}, nil
}

func (h *HTTPHandler) Serve(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()
	h.logger.Info("start serving websocket connection")
	userID, ok := middlewares.GetUserID(eCtx)
	if !ok {
		return fmt.Errorf("failed to get user id from context")
	}

	ws, err := h.upgrader.Upgrade(eCtx.Response(), eCtx.Request(), nil)
	if err != nil {
		return fmt.Errorf("failed to upgrade connection: %w", err)
	}
	closer := newWsCloser(h.logger, ws)
	defer func() {
		closer.Close(gorillaws.CloseNormalClosure)
	}()

	eventsChan, err := h.eventStream.Subscribe(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to subscribe to events: %w", err)
	}

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return h.readLoop(egCtx, ws)
	})
	eg.Go(func() error {
		return h.writeLoop(egCtx, ws, eventsChan)
	})

	go func() {
		if err := eg.Wait(); err != nil {
			h.logger.Error("failed to serve websocket connection", zap.Error(err))
			closer.Close(gorillaws.CloseInternalServerErr)
		}
	}()

	<-h.shutdownCh
	err = ws.WriteControl(
		gorillaws.CloseMessage,
		gorillaws.FormatCloseMessage(gorillaws.CloseNormalClosure, ""),
		time.Now().Add(writeTimeout),
	)
	if err != nil {
		return fmt.Errorf("failed to write close message: %w", err)
	}

	err = ws.Close()
	if err != nil {
		return fmt.Errorf("failed to close websocket connection: %w", err)
	}
	return nil
}

// readLoop listen PONGs.
func (h *HTTPHandler) readLoop(_ context.Context, ws Websocket) error {
	ws.SetPongHandler(func(appData string) error {
		err := ws.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return fmt.Errorf("SetReadDeadline: %w", err)
		}
		h.logger.Debug("pong")
		return nil
	})

	for {
		select {
		case <-h.shutdownCh:
			return nil
		default:
			t, reader, err := ws.NextReader()
			if err != nil {
				return fmt.Errorf("failed to get next reader: %w", err)
			}
			data, err := io.ReadAll(reader)
			if err != nil {
				return fmt.Errorf("failed to read data: %w", err)
			}
			h.logger.Debug("new message received", zap.Any("type", t), zap.Any("fromReader", data))
			if t != gorillaws.PongMessage {
				return fmt.Errorf("unexpected message type: %v", t)
			}
			h.logger.Debug("pong")
		}
	}
}

// writeLoop listen events and writes them into Websocket.
func (h *HTTPHandler) writeLoop(_ context.Context, ws Websocket, events <-chan eventstream.Event) error {
	ticker := time.NewTicker(h.pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-h.shutdownCh:
			return nil
		case event := <-events:

			adaptedEvent, err := h.eventAdapter.Adapt(event)
			if err != nil {
				return fmt.Errorf("failed to adapt event: %w", err)
			}
			h.logger.Info("new event received", zap.Any("event", adaptedEvent))

			err = ws.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err != nil {
				return fmt.Errorf("failed to set write deadline: %w", err)
			}

			w, err := ws.NextWriter(gorillaws.TextMessage)
			if err != nil {
				return fmt.Errorf("failed to get next writer: %w", err)
			}

			err = h.eventWriter.Write(adaptedEvent, w)
			if err != nil {
				return fmt.Errorf("failed to write event: %w", err)
			}

			h.logger.Info("event written", zap.Any("event", adaptedEvent))

			ticker.Reset(h.pingPeriod)
		case <-ticker.C:
			h.logger.Debug("ping")
			err := ws.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err != nil {
				return fmt.Errorf("failed to set write deadline: %w", err)
			}
			err = ws.WriteMessage(gorillaws.PingMessage, nil)
			if err != nil {
				return fmt.Errorf("failed to write ping message: %w", err)
			}
			ticker.Reset(h.pingPeriod)
		}
	}
}
