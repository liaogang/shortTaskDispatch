package task_type_cache

import (
	"fmt"
	"root/dispatch_centre"
)

//var itemMap = make(map[string]*dispatch_centre.TaskItem)

var tableTypeToIdMap = make(map[string]map[string]*dispatch_centre.TaskItem)

func Add(item *dispatch_centre.TaskItem) {

	itemMap, ok := tableTypeToIdMap[item.TaskType]
	if !ok {
		//todo
		itemMap = make(map[string]*dispatch_centre.TaskItem)
		tableTypeToIdMap[item.TaskType] = itemMap
	}

	itemMap[item.TaskId] = item
}

func Delete(taskType, taskId string) {

}

func Find(taskType string) (*dispatch_centre.TaskItem, error) {

	itemMap, ok := tableTypeToIdMap[taskType]
	if !ok {
		return nil, fmt.Errorf("no item in tableTypeToIdMap")
	}

	for _, item := range itemMap {
		return item, nil
	}

	return nil, fmt.Errorf("no item in itemMap")
}
