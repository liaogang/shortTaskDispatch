package task_type_cache

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"root/extend/model/Task"
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

type Wrap struct {
	err     error
	payload []byte
}

func (slf *Cache) DispatchAndWaitFinish(item *Task.Item, timeout time.Duration) ([]byte, error) {

	log.Trace().Interface("dispatch", item).Send()

	//send to claim wait
	slf.claimBufChannel <- item

	//wait for finish
	var finishChannel = make(chan Wrap)
	slf.mFinishChannel.Store(item.Id, finishChannel)

	var t = time.NewTimer(timeout)

	select {
	case <-t.C:
		return nil, fmt.Errorf("timeout")
	case wrap := <-finishChannel:
		return wrap.payload, wrap.err
	}

}

func (slf *Cache) ClaimAndWait(ctx context.Context) (*Task.Item, error) {

	log.Trace().Msg("ClaimAndWait")

	select {
	case taskItem := <-slf.claimBufChannel:
		//wait for a task come
		return taskItem, nil
	case <-ctx.Done():
		log.Info().Msg("http request cancelled")
		return nil, fmt.Errorf("http request cancelled")
	}

}

func (slf *Cache) Finish(id string, payload []byte, err error) error {

	if payload != nil && err != nil {
		return fmt.Errorf("logic err, pass payload or err")
	}

	log.Trace().Str("finish", id).Send()

	if val, ok := slf.mFinishChannel.LoadAndDelete(id); ok {

		var channel = val.(chan Wrap)

		channel <- Wrap{
			err:     err,
			payload: payload,
		}

		return nil
	} else {
		return fmt.Errorf("can no find this task id in dispatching")
	}

}
