package manager

import (
	"errors"
	"root/group"
	"sync"
)

var cache sync.Map // groupName -> *Channel

//var allTaskId sync.Map// taskId -> struct{}

func SetupGroupChannels(items []string) {

	for _, item := range items {
		c := group.NewChannel(item)

		cache.Store(item, c)
	}

}

func GetGroup(groupName string) (*group.Channel, error) {

	val, ok := cache.Load(groupName)
	if !ok {
		return nil, errors.New("group not exist")
	}

	return val.(*group.Channel), nil
}
