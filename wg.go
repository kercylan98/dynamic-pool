package dynamic_pool

import uuid "github.com/satori/go.uuid"

func newWorkGroup(workerNum int, workChannelMax int) *workGroup {
	slf := new(workGroup)
	slf.id = uuid.NewV4().String()
	slf.workers = make([]*worker, workerNum)
	slf.work = make(chan func() error, workChannelMax)
	for i := 0; i < workerNum; i++ {
		slf.workers[i] = newWorker().run(slf.work)
	}
	return slf
}

// 工作组
type workGroup struct {
	id      string
	workers []*worker
	work chan func() error
}

func (slf *workGroup) do(work func() error) {
	slf.work <- work
}

func (slf *workGroup) close(finish func(id string)) {
	for len(slf.work) == 0 {
		for i := 0; i < len(slf.workers); i++ {
			slf.workers[i].close()
		}
		finish(slf.id)
		close(slf.work)
	}
}