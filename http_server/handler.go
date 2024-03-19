package http_server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"root/dispatch_centre/shortid"
	"root/dispatch_centre/task_type_cache"
	"root/model/Task"
	"strconv"
	"strings"
	"time"
)

const (
	TaskTypeQQRegisterCaptcha = "qqRegisterCaptcha"
	TaskTypeQQRegisterSms     = "qqRegisterCheckOrSendSms"
)

var QQRegisterCaptcha task_type_cache.Cache
var QQRegisterSms task_type_cache.Cache

func publishTask(ctx *gin.Context) error {

	taskType := ctx.Query("taskType")
	timeout := ctx.Query("timeout")

	var taskId string

	var cache *task_type_cache.Cache
	switch taskType {
	case TaskTypeQQRegisterCaptcha:
		taskId = "cap_" + shortid.New()
		cache = &QQRegisterCaptcha
	case TaskTypeQQRegisterSms:
		taskId = "sms_" + shortid.New()
		cache = &QQRegisterCaptcha
	default:
		return fmt.Errorf("no this task type")
	}

	iTimeout, err := strconv.Atoi(timeout)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	ctx.Request.Body.Close()

	var item = &Task.Item{
		Id:   taskId,
		Body: body,
	}

	resp, err := cache.DispatchAndWaitFinish(item, time.Second*time.Duration(iTimeout))
	if err != nil {
		return err
	}

	ctx.Writer.Write(resp)

	return nil
}

func claimTask(ctx *gin.Context) error {

	taskType := ctx.Query("taskType")

	var cache *task_type_cache.Cache
	switch taskType {
	case TaskTypeQQRegisterCaptcha:
		cache = &QQRegisterCaptcha
	case TaskTypeQQRegisterSms:
		cache = &QQRegisterCaptcha
	default:
		return fmt.Errorf("no this task type")
	}

	task, err := cache.ClaimAndWait()
	if err != nil {
		return err
	}

	ctx.Header("TaskId", task.Id)
	ctx.Writer.Write(task.Body)

	return nil
}

func finishTask(ctx *gin.Context) error {

	taskId := ctx.Query("taskId")

	before, _, found := strings.Cut(taskId, "_")
	if !found {
		return fmt.Errorf("invalid task id format")
	}

	var cache *task_type_cache.Cache
	switch before {
	case "cap":
		cache = &QQRegisterCaptcha
	case "sms":
		cache = &QQRegisterCaptcha
	default:
		return fmt.Errorf("no this task type")
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	ctx.Request.Body.Close()

	return cache.Finish(taskId, body)
}
