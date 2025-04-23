package manager

import (
	"root/group"
	"sync"
)

var cache sync.Map // groupName -> *Channel

func SetupGroupChannels(items []string) {

	for _, item := range items {
		c := group.NewChannel(item)

		cache.Store(item, c)
	}

}

func GetGroup(groupName string) (*group.Channel, bool) {

	val, ok := cache.Load(groupName)
	if !ok {
		return nil, false
	}

	return val.(*group.Channel), true
}
