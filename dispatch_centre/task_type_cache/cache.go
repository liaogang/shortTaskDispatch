package task_type_cache

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"root/model/Task"
	"sync"
	"time"
)

type Cache struct {
	claimBufChannel chan *Task.Item
	mFinishChannel  sync.Map //string, channel []byte
}

func NewCache() *Cache {
	var slf = new(Cache)

	slf.claimBufChannel = make(chan *Task.Item, 20)

	return slf
}

func (slf *Cache) DispatchAndWaitFinish(item *Task.Item, timeout time.Duration) ([]byte, error) {

	log.Trace().Interface("dispatch", item).Send()

	//send to claim wait
	slf.claimBufChannel <- item

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

func (slf *Cache) ClaimAndWait() (*Task.Item, error) {

	log.Trace().Msg("ClaimAndWait")

	//wait for a task come
	var taskItem = <-slf.claimBufChannel

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
