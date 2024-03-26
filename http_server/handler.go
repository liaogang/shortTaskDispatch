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
	"time"
)

var NameToImpl = make(map[string]dispatch_centre.DispatchImpl)
var taskIdToImpl = make(map[string]dispatch_centre.DispatchImpl)

//const (
//	TaskTypeQQRegisterCaptcha = "qqRegisterCaptcha"
//	TaskTypeQQLoginCaptcha    = "qqLoginCaptcha"
//	TaskTypeQQRegisterSms     = "qqRegisterCheckOrSendSms"
//)
//
//var (
//	QQRegisterCaptcha = task_type_cache.NewCache()
//	QQLoginCaptcha    = task_type_cache.NewCache()
//	QQRegisterSms     = pin_code_task_cache.NewCache()
//)
//
//const (
//	PrefixQQRegisterCaptcha = "cap0"
//	PrefixQQLoginCaptcha    = "cap1"
//	PrefixQQRegisterSms     = "sms0"
//)

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

	var taskId = shortid.New()

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

	impl, ok := taskIdToImpl[taskId]
	if !ok {
		return fmt.Errorf("no this task type")
	}

	if errParam == "1" {
		return impl.Finish(taskId, nil, errors.New(string(body)))
	} else {
		return impl.Finish(taskId, body, nil)
	}

}
