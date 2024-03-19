package http_server

import (
	"github.com/gin-gonic/gin"
	"io"
	"root/dispatch_centre"
	"root/model/Task"
	"strconv"
	"time"
)

func publishTask(ctx *gin.Context) error {
	taskType := ctx.Query("taskType")
	timeout := ctx.Query("timeout")

	iTimeout, err := strconv.Atoi(timeout)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	ctx.Request.Body.Close()

	item := dispatch_centre.Dispatch(taskType, body)

	resp, err := dispatch_centre.WaitDone(item, time.Second*time.Duration(iTimeout))
	if err != nil {
		return err
	}

	ctx.Writer.Write(resp)
	return nil
}

func claimTask(ctx *gin.Context) error {

	taskType := ctx.Query("taskType")
	wait := ctx.Query("wait")

	var taskItem *Task.Item
	var err error
	if wait == "1" {
		taskItem, err = dispatch_centre.WaitClaimTask(taskType)
	} else {
		taskItem, err = dispatch_centre.TryClaimTask(taskType)
	}

	if err != nil {
		return err
	}

	ctx.Header("TaskId", taskItem.Id)
	ctx.Writer.Write(taskItem.Body)

	return nil
}

func finishTask(ctx *gin.Context) error {

	taskId := ctx.Query("taskId")

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	ctx.Request.Body.Close()

	return dispatch_centre.FinishTask(taskId, body)
}
