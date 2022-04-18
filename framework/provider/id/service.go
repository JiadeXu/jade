package id

import "github.com/rs/xid"

type JadeIDService struct {
}

func NewJadeIDService(params ...interface{}) (interface{}, error) {
	return &JadeIDService{}, nil
}

func (j *JadeIDService) NewID() string {
	return xid.New().String()
}
