package task_type_cache

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"root/model/Task"
)

var tableTypeToIdMap = make(map[string]map[string]*Task.Item)

func Add(item *Task.Item) {

	itemMap, ok := tableTypeToIdMap[item.Type]
	if !ok {
		//todo
		itemMap = make(map[string]*Task.Item)
		tableTypeToIdMap[item.Type] = itemMap
	}

	itemMap[item.Id] = item
}

func Delete(taskType, taskId string) {

	itemMap, ok := tableTypeToIdMap[taskType]
	if !ok {
		log.Trace().Msg("no this item map")
		return
	}

	delete(itemMap, taskId)
}

func Find(taskType string) (*Task.Item, error) {

	itemMap, ok := tableTypeToIdMap[taskType]
	if !ok {
		return nil, fmt.Errorf("no item in tableTypeToIdMap")
	}

	for _, item := range itemMap {
		return item, nil
	}

	return nil, fmt.Errorf("no item in itemMap")
}
