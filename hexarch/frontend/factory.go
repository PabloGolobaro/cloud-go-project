package frontend

import (
	"cloud-go-project/hexarch/core"
	"fmt"
)

type FrontEnd interface {
	Start(kv *core.KeyValueStore) error
}

func NewFrontend(frontend string) (FrontEnd, error) {
	switch frontend {
	case "rest":
		return &restFrontEnd{}, nil
	case "grpc":
		return &grpcFrontEnd{}, nil
	case "":
		return nil, fmt.Errorf("Frontend type not defined")
	default:
		return nil, fmt.Errorf("no such frontend %s", frontend)
	}
}
