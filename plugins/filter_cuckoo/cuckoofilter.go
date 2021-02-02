package impl

import (
	f "github.com/seiflotfy/cuckoofilter"
	"sync"
)

type CuckooFilterImpl struct {
	l  sync.Mutex
	cf *f.CuckooFilter
}

func (filter *CuckooFilterImpl) Open() error {

	filter.cf = f.NewCuckooFilter(10000000)
	return nil
}

func (filter *CuckooFilterImpl) Close() error {
	filter.l.Lock()
	defer filter.l.Unlock()
	return nil
}

func (filter *CuckooFilterImpl) Exists(key []byte) bool {
	filter.l.Lock()
	defer filter.l.Unlock()
	return filter.cf.Lookup(key)
}

func (filter *CuckooFilterImpl) Add(key []byte) error {
	filter.l.Lock()
	defer filter.l.Unlock()
	filter.cf.Insert(key)
	return nil
}

func (filter *CuckooFilterImpl) Delete(key []byte) error {
	filter.l.Lock()
	defer filter.l.Unlock()
	filter.cf.Delete(key)
	return nil
}
