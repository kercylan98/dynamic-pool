package dynamic_pool

func newWorker() *worker {
	return &worker{
		off: make(chan bool, 1),
	}
}

type worker struct {
	off 	chan bool
}

func (slf *worker) close() {
	slf.off <- true
}

func (slf *worker) run(work <- chan func() error) *worker {
	go func() {
		for w := range work {
			w()
			select {
			case off := <- slf.off:
				if off {
					close(slf.off)
					return
				}
			default:
				continue
			}
		}
	}()
	return slf
}