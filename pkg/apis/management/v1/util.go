package v1

import (
	"os"
	"path/filepath"

	"google.golang.org/grpc"
)

func DefaultManagementSocket() string {
	// check if we're root
	if os.Geteuid() == 0 {
		return "unix:///run/monty/management.sock"
	}
	// check if $XDG_RUNTIME_DIR is set
	if runUser, ok := os.LookupEnv("XDG_RUNTIME_DIR"); ok {
		return "unix://" + filepath.Join(runUser, "monty/management.sock")
	}
	return "unix:///tmp/monty/management.sock"
}

func UnderlyingConn(client ManagementClient) grpc.ClientConnInterface {
	return client.(*managementClient).cc
}
