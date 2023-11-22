package canreceiveproblems_test

import (
	"testing"

	"github.com/keepcalmist/chat-service/internal/types"
	canreceiveproblems "github.com/keepcalmist/chat-service/internal/usecases/manager/can-receive-problems"
)

func TestRequest_Validate(t *testing.T) {
	tCase := []struct {
		name    string
		request canreceiveproblems.Request
		wantErr bool
	}{
		{
			name: "valid",
			request: canreceiveproblems.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.NewUserID(),
			},
			wantErr: false,
		},
		{
			name: "empty request id",
			request: canreceiveproblems.Request{
				ID:        types.RequestIDNil,
				ManagerID: types.NewUserID(),
			},
			wantErr: true,
		},
		{
			name: "empty manager id",
			request: canreceiveproblems.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.UserIDNil,
			},
			wantErr: true,
		},
		{
			name: "empty struct",
			request: canreceiveproblems.Request{
				ID:        types.RequestIDNil,
				ManagerID: types.UserIDNil,
			},
			wantErr: true,
		},
	}

	for _, tc := range tCase {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.request.Validate()
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
