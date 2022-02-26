package bruteforce

import "sync"

type DictEntryQueue struct {

	lock sync.Mutex

	entries []*DictEntry

	indx int

	n int
}

func NewDictEntryQueue(entries []*DictEntry) *DictEntryQueue {

	return &DictEntryQueue{
		lock:    sync.Mutex{},
		entries: entries,
		indx:    0,
		n: len(entries),
	}

}

func (dq *DictEntryQueue)Pop() *DictEntry {

	dq.lock.Lock()
	defer dq.lock.Unlock()

	if dq.indx>=dq.n {
		return nil
	}

	entry := dq.entries[dq.indx]
	dq.indx=dq.indx+1

	return entry
}


