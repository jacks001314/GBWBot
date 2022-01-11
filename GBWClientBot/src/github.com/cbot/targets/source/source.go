package source

import "errors"

var endError error = errors.New("Source Read Over!")

type SourceEntry interface {

	IP() 		string
	Host() 		string
	Port() 		int
	Proto() 	string
	App()		string
}

type Source interface {

	Put(entry *SourceEntry) error

	OpenReader(name string,rtypes []string, capacity int) (*SourceReader,error)

	CloseReader(r *SourceReader)

	Start() error

	Stop()

	AtEnd()
}

type SourceReader struct {

	name 		string
	isEnd 		bool
	rch 		chan SourceEntry
	capacity 	int
	rtypes      []string
}

func NewSourceReader(name string,rtypes []string, capacity int) *SourceReader {

	return &SourceReader{
		name: name,
		isEnd:    false,
		rch:      make(chan SourceEntry),
		capacity: capacity,
		rtypes:   rtypes,
	}
}

func (r *SourceReader) Read() (SourceEntry,error) {

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

func (r *SourceReader) Push( entry SourceEntry) {

	r.rch<- entry
}
