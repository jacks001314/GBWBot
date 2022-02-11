package source

import (
	"github.com/cbot/targets"
	"sync"
)

type SourcePool struct {

	lock sync.Mutex

	readers map[string]*sourceReader

	sources map[string] Source

}

type sourceReader struct {

	name 		string

	rch 		chan targets.Target

	capacity 	int

	accept  func (targets.Target) bool
}

func NewSourcePool() *SourcePool {


	return &SourcePool{
		lock:    sync.Mutex{},
		readers: make(map[string]*sourceReader),
		sources: make(map[string]Source),
	}

}

func (p *SourcePool) StartSource(s Source)  {

	p.lock.Lock()

	name := s.Name()

	if _,ok := p.readers[name]; !ok {

		//no existed

		p.sources[name] = s
	}

	p.lock.Unlock()

	go s.Start()

}

func (p *SourcePool) StopSource(s Source) {

	p.lock.Lock()
	defer p.lock.Unlock()

	name := s.Name()

	if s,ok := p.sources[name]; ok {

		s.Stop()
		delete(p.sources,name)
	}
}


func (p *SourcePool) SubTarget(name string,capacity int,accept func (targets.Target) bool) chan targets.Target {

	p.lock.Lock()
	defer p.lock.Unlock()

	if v,ok := p.readers[name]; ok {

		//existed
		return v.rch
	}

	nr := &sourceReader{
		name:    name,
		rch:     make(chan targets.Target,capacity),
		capacity: capacity,
		accept:   accept,
	}

	p.readers[name] = nr

	return nr.rch
}

func (p *SourcePool) put(s Source,target targets.Target) {

	p.lock.Lock()
	defer p.lock.Unlock()

	for _,r:= range p.readers {

		if r.accept(target) {

			r.rch<-target
		}
	}

}