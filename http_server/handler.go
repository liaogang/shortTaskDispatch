package http_server

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"root/extend/model/Task"
	"root/extend/shortid"
	"root/internal/dispatch_centre"
	"strconv"
	"sync"
	"time"
)

var NameToImpl = make(map[string]dispatch_centre.DispatchImpl)
var taskIdToImpl sync.Map //=  make(map[string]dispatch_centre.DispatchImpl)

func genUniqueTaskId() string {

	for {
		var id = shortid.New()
		if _, has := taskIdToImpl.Load(id); has {
			//already have, skip
		} else {
			return id
		}
	}

}

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

	var reqCtx = ctx.Request.Context()

	impl, ok := NameToImpl[taskType]
	if !ok {
		return fmt.Errorf("no this task type")
	}

	pinCode := ctx.Query("pinCode")

	var taskId = genUniqueTaskId()

	taskIdToImpl.Store(taskId, impl)

	var item = &Task.Item{
		Id:   taskId,
		Body: body,
	}

	resp, err2 := impl.DispatchAndWaitFinish(reqCtx, item, time.Second*time.Duration(iTimeout), pinCode)
	if err2 != nil {
		return err2
	}

	ctx.Writer.Write(resp)

	return nil
}

func claimTask(ctx *gin.Context) error {

	taskType := ctx.Query("taskType")
	pinCode := ctx.Query("pinCode")

	impl, ok := NameToImpl[taskType]
	if !ok {
		return fmt.Errorf("no this task type")
	}

	task, err := impl.ClaimAndWait(ctx.Request.Context(), pinCode)
	if err != nil {
		return err
	}

	ctx.Header("TaskId", task.Id)
	ctx.Writer.Write(task.Body)

	return nil
}

func finishTask(ctx *gin.Context) error {

	taskId := ctx.Query("taskId")
	errParam := ctx.Query("bodyIsError")

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	ctx.Request.Body.Close()

	val, ok := taskIdToImpl.LoadAndDelete(taskId)
	if !ok {
		return fmt.Errorf("no this task type")
	}

	impl := val.(dispatch_centre.DispatchImpl)

	if errParam == "1" {
		return impl.Finish(taskId, nil, errors.New(string(body)))
	} else {
		return impl.Finish(taskId, body, nil)
	}

}
