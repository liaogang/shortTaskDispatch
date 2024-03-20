package pin_code_task_cache

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"root/model/Task"
	"sync"
	"time"
)

/*
需要接头密码
*/

type Cache struct {
	claimBufChannel sync.Map //pinCode, channel

	mFinishChannel sync.Map //string, channel []byte
}

func NewCache() *Cache {
	var slf = new(Cache)

	return slf
}

func (slf *Cache) DispatchAndWaitFinish(pinCode string, item *Task.Item, timeout time.Duration) ([]byte, error) {

	log.Trace().Interface("dispatch", item).Send()

	val, ok := slf.claimBufChannel.Load(pinCode)
	if !ok {
		return nil, fmt.Errorf("no client for this pincode waiting")
	}

	var channel = val.(chan *Task.Item)

	//send to claim wait
	channel <- item

	//wait for finish
	var finishChannel = make(chan []byte)
	slf.mFinishChannel.Store(item.Id, finishChannel)

	var t = time.NewTimer(timeout)

	select {
	case <-t.C:
		return nil, fmt.Errorf("timeout")
	case payload := <-finishChannel:
		return payload, nil
	}

}

func (slf *Cache) ClaimAndWait(pinCode string) (*Task.Item, error) {

	log.Trace().Msg("ClaimAndWait")

	var channel = make(chan *Task.Item)
	slf.claimBufChannel.Store(pinCode, channel)

	//wait for a task come
	var taskItem = <-channel

	return taskItem, nil
}

func (slf *Cache) Finish(id string, payload []byte) error {

	log.Trace().Str("finish", id).Send()

	if val, ok := slf.mFinishChannel.Load(id); ok {

		var channel = val.(chan []byte)

		channel <- payload

		return nil
	} else {
		return fmt.Errorf("can no find this task id in dispatching")
	}

}
