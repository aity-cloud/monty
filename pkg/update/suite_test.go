package update_test

import (
	"testing"

	_ "github.com/aity-cloud/monty/pkg/test/setup"
	"github.com/aity-cloud/monty/pkg/update"
	"github.com/aity-cloud/monty/pkg/update/noop"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
)

func TestPatch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Update Suite")
}

type mockHandler struct {
	streamHandlerCalls int
	update.UpdateTypeHandler
}

func newMockHandler() *mockHandler {
	return &mockHandler{
		UpdateTypeHandler: noop.NewSyncServer(),
	}
}

func (h *mockHandler) mockInterceptor() *mockInterceptor {
	return &mockInterceptor{
		mockHandler: h,
	}
}

type mockInterceptor struct {
	*mockHandler
}

func (i *mockInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		i.streamHandlerCalls++
		return handler(srv, stream)
	}
}
