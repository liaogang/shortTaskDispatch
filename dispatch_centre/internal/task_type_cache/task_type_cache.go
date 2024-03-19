package task_type_cache

import (
	"fmt"
	"root/model/Task"
	"sync"
	"time"
)

type Cache struct {
	TaskType        string
	claimBufChannel chan *Task.Item

	mFinishChannel sync.Map //string, chan *Task.Item
}

func (slf *Cache) DispatchAndWaitFinish(item *Task.Item, timeout time.Duration) error {

	//send to claim wait
	slf.claimBufChannel <- item

	//wait for finish
	var finishChannel = make(chan *CallbackItem)
	slf.mFinishChannel.Store(item.Id, finishChannel)

	var t = time.NewTimer(timeout)

	select {
	case <-t.C:
		print("")
	case callbackItem := <-finishChannel:
		fmt.Println(callbackItem)
	}

	return nil
}

type CallbackItem struct {
	Err     error
	Payload []byte
}

func (slf *Cache) ClaimAndWait() (*Task.Item, error) {

	//wait for a task come
	var taskItem = <-slf.claimBufChannel

	return taskItem, nil
}

func (slf *Cache) Finish(id string, payload []byte) error {

	if val, ok := slf.mFinishChannel.Load(id); ok {

		var channel = val.(chan *CallbackItem)

		var cb = &CallbackItem{
			Err:     nil,
			Payload: payload,
		}

		channel <- cb

		return nil
	} else {
		return fmt.Errorf("can no find this task id in dispatching")
	}

}
