package http_server

import (
	"github.com/gin-gonic/gin"
	"io"
	"root/dispatch_centre"
	"time"
)

func publishTask(ctx *gin.Context) error {
	taskType := ctx.Query("taskType")
	//timeout := ctx.Query("timeout")//todo

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	ctx.Request.Body.Close()

	item := dispatch_centre.Dispatch(taskType, body)

	resp, err := dispatch_centre.WaitDone(item, time.Second*180)

	ctx.Writer.Write(resp)

	return err
}

func claimTask(ctx *gin.Context) error {

	taskType := ctx.Query("taskType")

	taskItem, err := dispatch_centre.ClaimTask(taskType)
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
