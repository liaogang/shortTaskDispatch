package dispatch_centre

import (
	"fmt"
	"root/dispatch_centre/internal/task_type_cache"
	"root/model/Task"
)

/*
领任务的两个接口
*/

func ClaimTask(taskType string) (*Task.Item, error) {

	item, err := task_type_cache.Find(taskType)
	if err != nil {
		return nil, fmt.Errorf("no content, %w", err)
	}

	return item, nil
}

func FinishTask(taskId string, payload []byte) error {

	var c, ok = chanCache[taskId]
	if !ok {
		return fmt.Errorf("no this taskId, or task released by timeout")
	}

	c <- payload

	return nil
}
