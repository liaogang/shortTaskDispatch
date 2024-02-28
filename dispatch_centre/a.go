package dispatch_centre

import (
	"fmt"
	"root/dispatch_centre/internal/shortid"
	"root/dispatch_centre/internal/task_type_cache"
	"time"
)

/*
发布任务的两个接口
*/

type TaskItem struct {
	TaskId   string
	TaskType string
	//Begin    time.Time
	//Timeout  time.Duration
	Body []byte
}

var taskIdToType = make(map[string]string)

// Dispatch return taskId, 生成任务
func Dispatch(taskType string, body []byte) *TaskItem {

	var id = shortid.New()

	var item = &TaskItem{
		TaskId:   id,
		TaskType: taskType,
		//Begin:    time.Now(),
		//Timeout:  timeout,
		Body: body,
	}

	taskIdToType[id] = taskType

	task_type_cache.Add(item)

	return item
}

var chanCache = make(map[string]chan []byte)

// WaitDone 等任务完成
func WaitDone(item *TaskItem, timeout time.Duration) ([]byte, error) {

	var c = make(chan []byte)
	chanCache[item.TaskId] = c

	var t = time.NewTimer(timeout)

	select {
	case respData := <-c:

		//clean
		taskType, ok := taskIdToType[item.TaskId]
		if ok {
			task_type_cache.Delete(taskType, item.TaskId)
		}

		return respData, nil
	case <-t.C:

		//clean
		taskType, ok := taskIdToType[item.TaskId]
		if ok {
			task_type_cache.Delete(taskType, item.TaskId)
		}

		return nil, fmt.Errorf("timeout")
	}

}
