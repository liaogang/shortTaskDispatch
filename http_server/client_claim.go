package http_server

import (
	"fmt"
	"net/http"
	"root/manager"
)

func claimTask(w http.ResponseWriter, r *http.Request) error {

	query := r.URL.Query()

	group := query.Get("group")
	pinCode := query.Get("pinCode")
	clientName := query.Get("client_name")

	clientName = clientName + "_" + r.RemoteAddr

	channel, err := manager.GetGroup(group)
	if err != nil {
		return fmt.Errorf("no this task type, %s", group)
	}

	task, err := channel.ClaimAndWait(clientName, r.Context(), pinCode)
	if err != nil {
		return err
	}

	r.Header.Set("task-id", task.Id)
	_, _ = w.Write(task.Body)

	return nil
}
