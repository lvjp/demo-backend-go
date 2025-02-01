package misc

import (
	"strconv"

	"go.lvjp.me/demo-backend-go/pkg/buildinfo"
)

type Service interface {
	Version() Version
}

func NewService() Service {
	return &impl{}
}

type impl struct{}

func (*impl) Version() Version {
	info := buildinfo.Get()

	return Version{
		Revision:     info.Revision,
		RevisionTime: info.RevisionTime,
		GoVersion:    info.GoVersion,
		Platform:     info.GoOS + "/" + info.GoArch,
		Modified:     strconv.FormatBool(info.Modified),
	}
}
