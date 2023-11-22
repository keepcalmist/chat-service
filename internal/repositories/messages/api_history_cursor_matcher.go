package messagesrepo

import (
	"fmt"

	"github.com/golang/mock/gomock"
)

var _ gomock.Matcher = CursorMatcher{}

// CursorMatcher is intended to be used only in tests.
type CursorMatcher struct {
	c Cursor
}

func NewCursorMatcher(c Cursor) CursorMatcher {
	return CursorMatcher{c: c}
}

func (cm CursorMatcher) Matches(x interface{}) bool {
	v, ok := x.(*Cursor)
	if !ok {
		return false
	}

	if v == nil {
		return false
	}

	return cm.c.LastCreatedAt.Equal(v.LastCreatedAt) && cm.c.PageSize == v.PageSize
}

func (cm CursorMatcher) String() string {
	return fmt.Sprintf("{ps=%d, last_created_at=%d}", cm.c.PageSize, cm.c.LastCreatedAt.UnixNano())
}
