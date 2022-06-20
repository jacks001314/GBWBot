package store

import "fmt"

type Store interface {

	//open a database to store
	Open(cfg *Config) (Store,error)

	Close() error

	Put(key string,time uint64,value interface{}) error

	Get(key string,value interface{}) (bool,error)

	Del(key string) error


	Query(queryString string,timeRange [2]uint64,pageable *Pageable) (*QueryResult,error)

	FlushDB() error

	Count() uint64

	CountWithQuery(query string) uint64

	Facet(query string ,term string,num uint64,isDec bool) ([]*TermFacet,error)
}

type Pageable struct {

	Page uint64
	Size uint64
	ISDec  bool
}

type ResultEntry struct {

	Key string
	Value string
}

type QueryResult struct {

	Page uint64
	Size uint64
	TPage uint64
	TNum  uint64
	Results []*ResultEntry

}

type TermFacet struct {

	Key string  `json:"key"`
	Count uint64 `json:"count"`
}

func (t *TermFacet) String()string {

	return fmt.Sprintf(`{"key":%s,"count":%d}`,t.Key,t.Count)
}