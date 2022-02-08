package source

import (
	"errors"
	"github.com/cbot/targets"
)

var endError error = errors.New("Source Read Over!")

type Source interface {

	Put(target targets.Target) error

	OpenReader(name string,rtypes []string, capacity int) (*SourceReader,error)

	CloseReader(r *SourceReader)

	Start() error

	Stop()

	AtEnd()
}

type SourceReader struct {

	name 		string
	isEnd 		bool
	rch 		chan targets.Target
	capacity 	int
	rtypes      []string
}

func NewSourceReader(name string,rtypes []string, capacity int) *SourceReader {

	return &SourceReader{
		name: name,
		isEnd:    false,
		rch:      make(chan targets.Target),
		capacity: capacity,
		rtypes:   rtypes,
	}
}

func (r *SourceReader) Read() (targets.Target,error) {

	if r.isEnd {
		return nil,endError
	}

	select {
	case entry := <-r.rch:
		return entry,nil

	default:
		return nil,nil
	}


}

func (r *SourceReader) Push( entry targets.Target) {

	r.rch<- entry
}
