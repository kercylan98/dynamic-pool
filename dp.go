package dynamic_pool

import (
	"errors"
	"sync"
)

type DynamicPool struct {
	mutex 			*sync.Mutex
	initSync 		*sync.Once
	closeLocking 	bool

	joinChannel		chan func() error
	workGroups 		map[string]*workGroup
	runningWorkGroupId string
	isClose			bool
}

func (slf *DynamicPool) init() {
	if slf.initSync == nil {
		slf.initSync = new(sync.Once)
	}
	slf.initSync.Do(func() {
		wg := newWorkGroup(5, 1000)
		slf.mutex = new(sync.Mutex)
		slf.joinChannel = make(chan func() error, 1000)
		slf.workGroups = map[string]*workGroup{}
		slf.workGroups[wg.id] = wg
		slf.runningWorkGroupId = wg.id
	})
}

func (slf *DynamicPool) Close() {
	slf.init()
	slf.mutex.Lock()
	slf.isClose = true
	slf.mutex.Unlock()

	go func() {
		for len(slf.workGroups[slf.runningWorkGroupId].work) == 0 && len(slf.joinChannel) == 0 {
			close(slf.joinChannel)
			go slf.workGroups[slf.runningWorkGroupId].close(func(id string) {
				delete(slf.workGroups, id)
			})
			return
		}
	}()

}

func (slf *DynamicPool) Do(work func() error) error {
	slf.init()
	if slf.isClose {
		return errors.New("dp closed")
	}

	slf.workGroups[slf.runningWorkGroupId].do(work)
	return nil
}

func (slf *DynamicPool) NewWorkerNum(num int) {
	slf.init()
	slf.mutex.Lock()
	slf.closeLocking = true

	wg := newWorkGroup(num, 1000)
	slf.workGroups[wg.id] = wg
	go slf.workGroups[slf.runningWorkGroupId].close(func(id string) {
		if slf.closeLocking {
			delete(slf.workGroups, id)
		}else {
			slf.mutex.Lock()
			delete(slf.workGroups, id)
			slf.mutex.Unlock()
		}
	})
	slf.runningWorkGroupId = wg.id

	slf.mutex.Unlock()
	slf.closeLocking = false
}