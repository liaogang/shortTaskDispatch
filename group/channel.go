package group

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"root/extend/model/Task"
	"sync"
	"time"
)

type Channel struct {
	claimBufChannel chan *Task.Item

	mFinishChannel sync.Map //string, channel []byte

	group string
}

func NewChannel(groupName string) *Channel {
	var slf = new(Channel)

	slf.claimBufChannel = make(chan *Task.Item)
	slf.group = groupName

	return slf
}

type ResultItem struct {
	err     error
	payload []byte
}

func (slf *Channel) DispatchAndWaitFinish(ctx context.Context, item *Task.Item, timeout time.Duration, pinCode string) ([]byte, error) {

	_ = pinCode //no need here

	taskId := item.Id

	var t = time.NewTimer(timeout)

	log.Info().
		Str("taskId", taskId).
		Str("group", slf.group).
		Msg("开始发布，等待client端接收")

	select {
	case slf.claimBufChannel <- item:
		log.Info().Str("taskId", taskId).Msg("发送给client")

	case <-t.C:
		var err = fmt.Errorf("dispatch task %s timeout, no client now", item.Id)

		log.Info().Str("taskId", taskId).Err(err).Send()
		return nil, err
	}

	//wait for finish
	var finishChannel = make(chan ResultItem)
	slf.mFinishChannel.Store(taskId, finishChannel)

	defer slf.mFinishChannel.Delete(taskId)

	select {
	case <-ctx.Done():
		log.Info().Str("tasId", taskId).Msg("发布者取消了任务")
		return nil, fmt.Errorf("user cancel")
	case <-t.C:
		log.Info().Str("tasId", taskId).Msg("等待结果超时")
		return nil, fmt.Errorf("timeout")
	case wrap := <-finishChannel:
		log.Info().Str("tasId", taskId).Msg("收到任务结果")
		return wrap.payload, wrap.err
	}

}

func (slf *Channel) ClaimAndWait(workerTag string, ctx context.Context, pinCode string) (*Task.Item, error) {

	_ = pinCode //no need here

	log.Info().
		Str("workerTag", workerTag).
		Str("group", slf.group).
		Msg("认领任务 ..")

	select {
	case taskItem := <-slf.claimBufChannel:

		log.Info().
			Str("workerTag", workerTag).
			Str("taskId", taskItem.Id).
			Msg("认领到了一个任务")

		return taskItem, nil
	case <-ctx.Done():

		log.Info().
			Str("workerTag", workerTag).
			Msg("认领任务 http request 已取消")

		return nil, fmt.Errorf("http request cancelled")
	}

}

func (slf *Channel) Finish(workerTag string, taskId string, payload []byte, err error) error {

	log.Info().
		Str("workerTag", workerTag).
		Str("taskId", taskId).
		Msg("提交任务结果")

	if payload != nil && err != nil {
		return fmt.Errorf("logic err, pass payload or err")
	}

	if val, ok := slf.mFinishChannel.LoadAndDelete(taskId); ok {

		var channel = val.(chan ResultItem)

		channel <- ResultItem{
			err:     err,
			payload: payload,
		}

		log.Info().
			Str("workerTag", workerTag).
			Str("taskId", taskId).
			Msg("已提交结果给任务发布方")

		return nil
	} else {

		log.Info().
			Str("workerTag", workerTag).
			Str("taskId", taskId).
			Msg("can no find this task taskId in dispatching, 可能已经超时? 发布都取消了?")

		return fmt.Errorf("can no find this task taskId in dispatching")
	}

}
