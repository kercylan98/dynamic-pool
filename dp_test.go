package dynamic_pool

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestDynamicPool_Do(t *testing.T) {
	dp := new(DynamicPool)
	var waitGroup = new(sync.WaitGroup)
	for i := 0; i < 150; i++ {
		count := i
		if i == 50 {
			dp.NewWorkerNum(100)
		}
		waitGroup.Add(1)
		if err := dp.Do(func() error {
			waitGroup.Done()
			fmt.Println(count)
			time.Sleep(1 * time.Second)
			return nil
		}); err != nil {
			waitGroup.Done()
			t.Log(err)
		}
	}

	time.Sleep(3 * time.Second)

	dp.Close()

	t.Log("failed,", dp.Do(func() error {
		return nil
	}))

	waitGroup.Wait()
}