package messagesrepo

import (
	"context"
	"errors"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/keepcalmist/chat-service/internal/store/chat"
	"github.com/keepcalmist/chat-service/internal/store/message"
	"github.com/keepcalmist/chat-service/internal/types"
)

const (
	maxPageSize = 100
	minPageSize = 10
)

var (
	ErrInvalidPageSize = errors.New("invalid page size")
	ErrInvalidCursor   = errors.New("invalid cursor")
)

type Cursor struct {
	LastCreatedAt time.Time
	PageSize      int
}

// GetClientChatMessages returns Nth page of messages in the chat for client side.
func (r *Repo) GetClientChatMessages(
	ctx context.Context,
	clientID types.UserID,
	pageSize int,
	cursor *Cursor,
) ([]Message, *Cursor, error) {
	var (
		size      int
		createdAt *time.Time
	)

	switch {
	case cursor != nil:
		if cursor.PageSize < minPageSize || cursor.PageSize > maxPageSize {
			return nil, nil, ErrInvalidCursor
		}
		if cursor.LastCreatedAt.IsZero() {
			return nil, nil, ErrInvalidCursor
		}
		size = cursor.PageSize
		createdAt = &cursor.LastCreatedAt
	default:
		if pageSize < 10 || pageSize > maxPageSize {
			return nil, nil, ErrInvalidPageSize
		}
		size = pageSize
	}

	// 1) Если указан pageSize, то валидируем, что он в пределах [10, 100] и используем его для запроса.

	// 2) Если указан cursor, то используем данные из него, предварительно валидируя:
	// pageSize аналогично пункту выше и LastCreatedAt на заполненность в принципе.
	// FIXME: 4) Возвращаем очередную страницу сообщений в соответствии с параметрами запроса.
	// FIXME: Первое сообщение является последним по времени своего создания (наиболее свежее).
	qb := r.db.Message(ctx).
		Query().
		Unique(true).
		Where(
			message.HasChatWith(chat.ClientIDEQ(clientID)),
			message.IsVisibleForClientEQ(true),
		).
		Order(message.ByCreatedAt(func(options *sql.OrderTermOptions) {
			options.Desc = true
		})).Limit(size + 1)
	if createdAt != nil {
		qb = qb.Where(message.CreatedAtLT(*createdAt))
	}
	msgs, err := qb.All(ctx)
	if err != nil {
		return nil, nil, err
	}
	retVal := make([]Message, 0, len(msgs))

	for _, msg := range msgs {
		retVal = append(retVal, adaptStoreMessage(msg))
	}
	// FIXME: 5) Если впереди есть ещё страницы, то возвращаем курсор на следующую страницу, иначе nil.
	//
	// FIXME: 6) Пользуемся TEST_PSQL_DEBUG, чтобы понять, не превращает ли ent наш код в SQL-запрос, похожий на дичь.
	//
	// FIXME: 7) Отдельно обратите внимание на то, что
	// FIXME: - нужно доставать сообщения из клиентского чата (чужие чаты не должны попадать в выборку);
	// FIXME: - нужно доставать сообщения, видимые клиенту.
	if len(retVal) > size {
		return retVal[:size], &Cursor{
			PageSize:      size,
			LastCreatedAt: retVal[size-1].CreatedAt,
		}, nil
	}

	return retVal, nil, nil
}
