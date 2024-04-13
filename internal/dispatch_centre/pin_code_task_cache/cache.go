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

	var id = item.Id

	log.Info().
		Str("pinCode", pinCode).
		Str("task", id).
		Interface("a, 发布任务", item).Send()

	val, ok := slf.claimBufChannel.LoadAndDelete(pinCode)
	if !ok {
		log.Info().
			Str("pinCode", pinCode).
			Str("task", id).
			Msg("no client for this pinCode waiting")
		return nil, fmt.Errorf("no client for this pinCode waiting")
	}

	var channel = val.(chan *Task.Item)

	log.Info().
		Str("task", id).
		Str("pinCode", pinCode).
		Msg("a, 等待任务被认领")

	channel <- item

	log.Info().
		Str("pinCode", pinCode).
		Str("task", id).
		Msg("a, 任务已被认领, 等待任务完成")

	var finishChannel = make(chan Wrap)
	slf.mFinishChannel.Store(item.Id, finishChannel)

	var t = time.NewTimer(timeout)

	select {
	case <-ctx.Done():

		log.Info().
			Str("task", id).
			Str("pinCode", pinCode).
			Msg("a, 用户取消")

		slf.mFinishChannel.Delete(id)
		return nil, fmt.Errorf("timeout")
	case <-t.C:

		log.Info().
			Str("task", id).
			Str("pinCode", pinCode).
			Msg("a, 超时")

		slf.mFinishChannel.Delete(id)
		return nil, fmt.Errorf("timeout")
	case wrap := <-finishChannel:

		log.Info().
			Str("task", id).
			Str("pinCode", pinCode).
			Msg("a, 已完成")

		return wrap.payload, wrap.err
	}

}

func (slf *Cache) ClaimAndWait(workerTag string, ctx context.Context, pinCode string) (*Task.Item, error) {

	log.Info().Str("workerTag", workerTag).Str("pinCode", pinCode).Msg("b, 使用对接码来认领任务")

	var channel = make(chan *Task.Item)
	slf.claimBufChannel.Store(pinCode, channel)

	//wait for a task come
	var t = time.NewTimer(time.Hour)
	select {
	case <-ctx.Done():

		log.Info().Str("workerTag", workerTag).Str("pinCode", pinCode).Msg("b, 用户取消")

		slf.claimBufChannel.Delete(pinCode)
		return nil, fmt.Errorf("user cancel")
	case <-t.C:

		log.Info().Str("workerTag", workerTag).Str("pinCode", pinCode).Msg("b, 认领任务超时")

		slf.claimBufChannel.Delete(pinCode)
		return nil, fmt.Errorf("timeout")
	case taskItem := <-channel:

		log.Info().
			Str("task", taskItem.Id).
			Str("workerTag", workerTag).
			Str("pinCode", pinCode).
			Msg("b, 已认领到一个任务")

		return taskItem, nil
	}

}

func (slf *Cache) Finish(workerTag, id string, payload []byte, err error) error {

	log.Info().Str("workerTag", workerTag).Str("taskId", id).Msg("b, 完成任务")

	if val, ok := slf.mFinishChannel.LoadAndDelete(id); ok {

		var channel = val.(chan Wrap)

		log.Info().Str("workerTag", workerTag).Str("taskId", id).Msg("b, 发送给发布方")
		channel <- Wrap{
			err:     err,
			payload: payload,
		}

		log.Info().Str("workerTag", workerTag).Str("taskId", id).Msg("b, 对方已接收")

		return nil
	} else {

		log.Error().
			Str("workerTag", workerTag).
			Str("taskId", id).
			Msg("b, can no find this task id in dispatching")

		return fmt.Errorf("can no find this task id in dispatching")
	}

}
