package service

import (
	"picture-oss-proxy/pkg/e"
	"picture-oss-proxy/serializer"
)

type PictureService struct {
}

func (service PictureService) Get() serializer.Response {
	code := e.SUCCESS
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}

}
