package dispatch_centre

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"root/dispatch_centre/internal/shortid"
	"root/dispatch_centre/internal/task_type_cache"
	"root/model/Task"
	"time"
)

/*
发布任务的两个接口
*/
var taskIdToType = make(map[string]string)

// Dispatch return taskId, 生成任务
func Dispatch(taskType string, body []byte) *Task.Item {

	log.Trace().Str("dispatch task", taskType).Send()

	var id = shortid.New()

	var item = &Task.Item{
		Id:   id,
		Type: taskType,
		Body: body,
	}

	taskIdToType[id] = taskType

	task_type_cache.Dispatch(item)

	return item
}

var chanCache = make(map[string]chan []byte)

// WaitDone 等任务完成
func WaitDone(item *Task.Item, timeout time.Duration) ([]byte, error) {

	log.Trace().Str("wait done", item.Id).Send()

	var id = item.Id

	var c = make(chan []byte)
	chanCache[id] = c

	var t = time.NewTimer(timeout)

	select {
	case respData := <-c:

		//clean
		taskType, ok := taskIdToType[id]
		if ok {
			task_type_cache.Delete(taskType, id)
		}

		delete(chanCache, id)

		return respData, nil
	case <-t.C:

		//clean
		taskType, ok := taskIdToType[id]
		if ok {
			task_type_cache.Delete(taskType, id)
		}

		delete(chanCache, id)

		return nil, fmt.Errorf("timeout")
	}

}
