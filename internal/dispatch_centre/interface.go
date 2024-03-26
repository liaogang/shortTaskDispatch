package dispatch_centre

import (
	"context"
	"root/extend/model/Task"
	"time"
)

type DispatchImpl interface {
	DispatchAndWaitFinish(ctx context.Context, item *Task.Item, timeout time.Duration, pinCode string) ([]byte, error)

	ClaimAndWait(ctx context.Context, pinCode string) (*Task.Item, error)
	Finish(id string, payload []byte, err error) error
}
