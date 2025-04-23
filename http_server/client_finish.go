package http_server

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"root/manager"
)

func finishTaskWithSuccess(w http.ResponseWriter, r *http.Request) error {

	query := r.URL.Query()

	group := query.Get("group")

	channel, ok := manager.GetGroup(group)
	if !ok {
		return fmt.Errorf("manger get group fail, %s", group)
	}

	taskId := query.Get("task_id")
	clientName := query.Get("client_name")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	_ = r.Body.Close()

	return channel.Finish(clientName, taskId, body, nil)
}

func finishTaskWithError(w http.ResponseWriter, r *http.Request) error {

	query := r.URL.Query()

	group := query.Get("group")

	channel, ok := manager.GetGroup(group)
	if !ok {
		return fmt.Errorf("manger get group fail, %s", group)
	}

	taskId := query.Get("task_id")
	clientName := query.Get("client_name")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	_ = r.Body.Close()

	return channel.Finish(clientName, taskId, nil, errors.New(string(body)))
}
