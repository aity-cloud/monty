package mock_ident

import (
	"github.com/aity-cloud/monty/pkg/ident"
	"go.uber.org/mock/gomock"
)

func NewTestIdentProvider(ctrl *gomock.Controller, id string) ident.Provider {
	mockIdent := NewMockProvider(ctrl)
	mockIdent.EXPECT().
		UniqueIdentifier(gomock.Any()).
		Return(id, nil).
		AnyTimes()
	return mockIdent
}
