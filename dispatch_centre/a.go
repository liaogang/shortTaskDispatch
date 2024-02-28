package dispatch_centre

import (
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
	Begin    time.Time
	Timeout  time.Duration
	Body     []byte
}

var taskIdToType = make(map[string]string)

// Dispatch return taskId, 生成任务
func Dispatch(taskType string, body []byte, timeout time.Duration) string {

	var id = shortid.New()

	var item = &TaskItem{
		TaskId:   id,
		TaskType: taskType,
		Begin:    time.Now(),
		Timeout:  timeout,
		Body:     body,
	}

	taskIdToType[id] = taskType

	task_type_cache.Add(item)

	return id
}

// WaitDone 等任务完成
func WaitDone(taskId string) error {

	var c = make(chan struct{})
	<-c

	taskType, ok := taskIdToType[taskId]
	if ok {
		task_type_cache.Delete(taskType, taskId)
	}

	return nil
}
