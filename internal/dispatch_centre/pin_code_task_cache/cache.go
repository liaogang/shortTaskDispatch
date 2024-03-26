package pin_code_task_cache

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"root/extend/model/Task"
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

type Wrap2 struct {
	id      string
	err     error
	payload []byte
}

type Wrap struct {
	err     error
	payload []byte
}

func (slf *Cache) DispatchAndWaitFinish(ctx context.Context, item *Task.Item, timeout time.Duration, pinCode string) ([]byte, error) {

	log.Info().Str("pinCode", pinCode).Interface("发布任务", item).Send()

	val, ok := slf.claimBufChannel.LoadAndDelete(pinCode)
	if !ok {
		log.Info().Str("pinCode", pinCode).Msg("no client for this pinCode waiting")
		return nil, fmt.Errorf("no client for this pinCode waiting")
	}

	var channel = val.(chan *Task.Item)

	log.Info().Str("pinCode", pinCode).Msg("等待任务被认领")
	channel <- item

	log.Info().Str("pinCode", pinCode).Msg("任务已被认领, 等待任务完成")
	var finishChannel = make(chan Wrap)
	slf.mFinishChannel.Store(item.Id, finishChannel)

	var t = time.NewTimer(timeout)

	select {
	case <-ctx.Done():
		log.Info().Str("pinCode", pinCode).Msg("用户取消")
		slf.mFinishChannel.Delete(item.Id)
		return nil, fmt.Errorf("timeout")
	case <-t.C:
		log.Info().Str("pinCode", pinCode).Msg("超时")
		slf.mFinishChannel.Delete(item.Id)
		return nil, fmt.Errorf("timeout")
	case wrap := <-finishChannel:
		log.Info().Str("pinCode", pinCode).Msg("已完成")
		return wrap.payload, wrap.err
	}

}

func (slf *Cache) ClaimAndWait(ctx context.Context, pinCode string) (*Task.Item, error) {

	log.Info().Str("pinCode", pinCode).Msg("认领任务")

	var channel = make(chan *Task.Item)
	slf.claimBufChannel.Store(pinCode, channel)

	//wait for a task come
	var t = time.NewTimer(time.Hour)
	select {
	case <-ctx.Done():
		log.Info().Str("pinCode", pinCode).Msg("用户取消")
		slf.claimBufChannel.Delete(pinCode)
		return nil, fmt.Errorf("user cancel")
	case <-t.C:
		log.Info().Str("pinCode", pinCode).Msg("认领任务超时")
		slf.claimBufChannel.Delete(pinCode)
		return nil, fmt.Errorf("timeout")
	case taskItem := <-channel:
		log.Info().Str("pinCode", pinCode).Msg("已认领到一个任务")
		return taskItem, nil
	}

}

func (slf *Cache) Finish(id string, payload []byte, err error) error {

	log.Info().Str("taskId", id).Msg("完成任务")

	//send taskId and payload ?

	if val, ok := slf.mFinishChannel.LoadAndDelete(id); ok {

		var channel = val.(chan Wrap)

		log.Info().Str("taskId", id).Msg("发送给发布方")
		channel <- Wrap{
			err:     err,
			payload: payload,
		}

		log.Info().Str("taskId", id).Msg("对方已接收")

		return nil
	} else {
		log.Error().Str("taskId", id).Msg("can no find this task id in dispatching")
		return fmt.Errorf("can no find this task id in dispatching")
	}

}
