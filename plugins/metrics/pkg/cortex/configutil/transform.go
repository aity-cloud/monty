package configutil

import (
	"buf.build/go/protoyaml"
	"github.com/cortexproject/cortex/pkg/util/flagext"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	yamlv2 "gopkg.in/yaml.v2"
)

func loadDefaults[T flagext.Registerer](dest T) {
	flagext.DefaultValues(dest)
}

func LoadFromAPI[T any](dest T, src proto.Message) error {
	if src == nil {
		return nil
	}
	m, err := protoyaml.Marshal(src)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to marshal config for api type %s: %v", proto.MessageName(src), err)
	}
	if err := yamlv2.Unmarshal(m, dest); err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to load config for api type %s into %T: %v", proto.MessageName(src), dest, err)
	}
	return nil
}
