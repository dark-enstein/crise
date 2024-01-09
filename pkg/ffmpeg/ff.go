package ffmpeg

import (
	"github.com/lijo-jose/gffmpeg/pkg/gffmpeg"
	"strconv"
)

type Service interface {
	ExtractFrames(inFile, outDir string, fps int) error
}

type service struct {
	ff gffmpeg.GFFmpeg
}

func New(ff gffmpeg.GFFmpeg) (Service, error) {
	return &service{ff: ff}, nil
}

func (svc *service) ExtractFrames(inFile, outDir string, fps int) error {
	bd := gffmpeg.NewBuilder()
	bd = bd.SrcPath(inFile).VideoFilters("fps=" + strconv.Itoa(fps)).DestPath(outDir)
	svc.ff = svc.ff.Set(bd)

	ret := svc.ff.Start(nil)
	if ret.Err != nil {
		return ret.Err
	}
	return nil
}
