package http_server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"io"
	"root/extend/model/Task"
	"root/internal/dispatch_centre/pin_code_task_cache"
	"root/internal/dispatch_centre/shortid"
	"root/internal/dispatch_centre/task_type_cache"
	"strconv"
	"strings"
	"time"
)

const (
	TaskTypeQQRegisterCaptcha = "qqRegisterCaptcha"
	TaskTypeQQRegisterSms     = "qqRegisterCheckOrSendSms"
)

var QQRegisterCaptcha = task_type_cache.NewCache()
var QQRegisterSms = pin_code_task_cache.NewCache()

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

	var taskId string

	switch taskType {
	case TaskTypeQQRegisterCaptcha:
		taskId = "cap_" + shortid.New()

		var item = &Task.Item{
			Id:   taskId,
			Body: body,
		}

		resp, err2 := QQRegisterCaptcha.DispatchAndWaitFinish(item, time.Second*time.Duration(iTimeout))
		if err2 != nil {
			return err2
		}

		ctx.Writer.Write(resp)

	case TaskTypeQQRegisterSms:
		taskId = "sms_" + shortid.New()

		var item = &Task.Item{
			Id:   taskId,
			Body: body,
		}

		pinCode := ctx.Query("pinCode")

		resp, err2 := QQRegisterSms.DispatchAndWaitFinish(pinCode, item, time.Second*time.Duration(iTimeout))
		if err2 != nil {
			return err2
		}

		ctx.Writer.Write(resp)
	default:
		return fmt.Errorf("no this task type")
	}

	return nil
}

func claimTask(ctx *gin.Context) error {

	taskType := ctx.Query("taskType")

	switch taskType {
	case TaskTypeQQRegisterCaptcha:

		task, err := QQRegisterCaptcha.ClaimAndWait(ctx.Request.Context())
		if err != nil {
			return err
		}

		ctx.Header("TaskId", task.Id)
		ctx.Writer.Write(task.Body)

	case TaskTypeQQRegisterSms:

		pinCode := ctx.Query("pinCode")

		task, err := QQRegisterSms.ClaimAndWait(pinCode)
		if err != nil {
			return err
		}

		ctx.Header("TaskId", task.Id)
		ctx.Writer.Write(task.Body)

	default:
		return fmt.Errorf("no this task type")
	}

	return nil
}

func finishTask(ctx *gin.Context) error {

	taskId := ctx.Query("taskId")

	errParam := ctx.Query("error")

	before, _, found := strings.Cut(taskId, "_")
	if !found {
		return fmt.Errorf("invalid task id format")
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	ctx.Request.Body.Close()

	switch before {
	case "cap":

		if errParam == "1" {
			return QQRegisterCaptcha.Finish(taskId, nil, errors.New(string(body)))
		} else {
			return QQRegisterCaptcha.Finish(taskId, body, nil)
		}

	case "sms":
		if errParam == "1" {
			return QQRegisterSms.Finish(taskId, nil, errors.New(string(body)))
		} else {
			return QQRegisterSms.Finish(taskId, body, nil)
		}
	default:
		return fmt.Errorf("no this task type")
	}

}
