package http_server

import (
	"fmt"
	"io"
	"net/http"
	"root/extend/model/Task"
	"root/extend/shortid"
	"root/manager"
	"strconv"
	"time"
)

func publishTask(w http.ResponseWriter, r *http.Request) error {

	query := r.URL.Query()

	group := query.Get("group")
	timeout := query.Get("timeout")

	channel, ok := manager.GetGroup(group)
	if !ok {
		return fmt.Errorf("manger get group fail, %s", group)
	}

	iTimeout, err := strconv.Atoi(timeout)
	if err != nil {
		return fmt.Errorf("parse timeout err: %v", err)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read request body err: %v", err)
	}

	_ = r.Body.Close()

	var reqCtx = r.Context()

	var taskId = shortid.New()

	var item = &Task.Item{
		Id:   taskId,
		Body: body,
	}

	durTimeout := time.Second * time.Duration(iTimeout)

	pinCode := query.Get("pinCode")
	resp, err2 := channel.DispatchAndWaitFinish(reqCtx, item, durTimeout, pinCode)
	if err2 != nil {
		return fmt.Errorf("dispatch fail: %w", err2)
	}

	_, _ = w.Write(resp)

	return nil
}
