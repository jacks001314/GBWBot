package logstream

import "sync"

type LogStream struct {

	lock sync.Mutex
	logSubs map[string]*LogSub
}

type LogSub struct {

	name string
	LogChan chan []byte
}


func NewLogStream() *LogStream {

	return &LogStream{
		lock:    sync.Mutex{},
		logSubs: make(map[string]*LogSub,0),
	}

}

func (log *LogStream) Sub(name string ) *LogSub {

	log.lock.Lock()
	defer log.lock.Unlock()

	ls := &LogSub{name: name,LogChan:make(chan []byte)}

	log.logSubs[name] = ls

	return ls
}

func (log *LogStream) UnSub(ls *LogSub)  {

	log.lock.Lock()
	defer log.lock.Unlock()
	delete(log.logSubs,ls.name)

}

func (log *LogStream) Log(message []byte) {

	log.lock.Lock()
	defer log.lock.Unlock()

	if len(log.logSubs)>0 {

		for _,ls:= range log.logSubs {

			ls.LogChan<-message
		}
	}
}