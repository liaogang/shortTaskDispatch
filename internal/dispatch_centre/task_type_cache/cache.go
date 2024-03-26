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

func (slf *Cache) DispatchAndWaitFinish(ctx context.Context, item *Task.Item, timeout time.Duration, pinCode string) ([]byte, error) {

	_ = pinCode //no need here

	log.Info().Str("taskId", item.Id).Msg("发布任务并等待认领")

	//send to claim wait
	slf.claimBufChannel <- item

	log.Info().Str("tasId", item.Id).Msg("等待结束")
	//wait for finish
	var finishChannel = make(chan Wrap)
	slf.mFinishChannel.Store(item.Id, finishChannel)

	var t = time.NewTimer(timeout)

	select {
	case <-ctx.Done():
		log.Info().Str("tasId", item.Id).Msg("用户取消")
		slf.mFinishChannel.Delete(item.Id)
		return nil, fmt.Errorf("user cancel")
	case <-t.C:
		log.Info().Str("tasId", item.Id).Msg("等待超时")
		slf.mFinishChannel.Delete(item.Id)
		return nil, fmt.Errorf("timeout")
	case wrap := <-finishChannel:
		log.Info().Str("tasId", item.Id).Msg("收到任务结果")
		return wrap.payload, wrap.err
	}

}

func (slf *Cache) ClaimAndWait(ctx context.Context, pinCode string) (*Task.Item, error) {

	_ = pinCode //no need here

	log.Info().Msg("认领任务..")

	select {
	case taskItem := <-slf.claimBufChannel:
		log.Info().Msg("认领到了一个任务")
		return taskItem, nil
	case <-ctx.Done():
		log.Info().Msg("认领任务 http request 已取消")
		return nil, fmt.Errorf("http request cancelled")
	}

}

func (slf *Cache) Finish(id string, payload []byte, err error) error {

	log.Info().Str("taskId", id).Msg("提交任务结果")

	if payload != nil && err != nil {
		return fmt.Errorf("logic err, pass payload or err")
	}

	if val, ok := slf.mFinishChannel.LoadAndDelete(id); ok {

		var channel = val.(chan Wrap)

		channel <- Wrap{
			err:     err,
			payload: payload,
		}

		log.Info().Str("taskId", id).Msg("已提交结果给任务发布方")

		return nil
	} else {
		log.Info().Str("taskId", id).Msg("can no find this task id in dispatching, 可能已经超时? 发布都取消了?")
		return fmt.Errorf("can no find this task id in dispatching")
	}

}
