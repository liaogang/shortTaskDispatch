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

	channel, err := manager.GetGroup(group)
	if err != nil {
		return fmt.Errorf("no this task of group: %s", group)
	}

	taskId := query.Get("taskId")
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

	channel, err := manager.GetGroup(group)
	if err != nil {
		return fmt.Errorf("no this task of group: %s", group)
	}

	taskId := query.Get("taskId")
	clientName := query.Get("client_name")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	_ = r.Body.Close()

	return channel.Finish(clientName, taskId, nil, errors.New(string(body)))
}
