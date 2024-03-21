package handlers

import (
	"context"
	"fmt"
	"github.com/Braly-Ltd/t2v-api-core/constants"
	"github.com/Braly-Ltd/t2v-api-core/entities"
	"github.com/Braly-Ltd/t2v-api-core/ports"
	"github.com/Braly-Ltd/t2v-api-worker/properties"
	"github.com/golibs-starter/golib/log"
	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
)

type T2VHandler struct {
	objectStoragePort    ports.ObjectStoragePort
	inferencePort        ports.InferencePort
	notiSubscriptionPort ports.NotificationSubscriptionPort
	notificationPort     ports.NotificationPort
	fileProps            *properties.FileProperties
}

func NewT2VHandler(
	objectStoragePort ports.ObjectStoragePort,
	inferencePort ports.InferencePort,
	notiSubscriptionPort ports.NotificationSubscriptionPort,
	notificationPort ports.NotificationPort,
	fileProps *properties.FileProperties,
) *T2VHandler {
	return &T2VHandler{
		objectStoragePort:    objectStoragePort,
		inferencePort:        inferencePort,
		notiSubscriptionPort: notiSubscriptionPort,
		notificationPort:     notificationPort,
		fileProps:            fileProps,
	}
}

func (r *T2VHandler) Type() constants.TaskType {
	return constants.TaskTypeT2V
}

// Handle
// 1. Download file from MinIO
// 2. Process file
// 3. Upload processed file to MinIO
func (r *T2VHandler) Handle(ctx context.Context, task *asynq.Task) error {
	var payload entities.InferPayload
	if err := msgpack.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unpack task failed: %v", err)
	}
	taskID := task.ResultWriter().TaskID()
	log.Infoc(ctx, "task %s is processing", taskID)
	log.Debugc(ctx, "task payload: %+v", payload)

	cmd := entities.InferenceCommand{
		Prompt:            payload.Prompt,
		NegativePrompt:    payload.NegativePrompt,
		NumInferenceSteps: payload.NumInferenceSteps,
		NumFrames:         payload.NumFrames,
		Width:             payload.Width,
		Height:            payload.Height,
		GuidanceScale:     payload.GuidanceScale,
		OutputFilePath:    fmt.Sprintf("%s/%s", r.fileProps.BaseOutputPath, payload.TargetFileName),
		ModelID:           payload.Model,
		UserID:            payload.UserID,
	}
	if err := r.inferencePort.Infer(ctx, cmd); err != nil {
		return err
	}
	log.Infoc(ctx, "task %s inference completed, start uploading file at %s", taskID, payload.TargetFileName)

	if err := r.objectStoragePort.UploadFilePath(ctx, cmd.OutputFilePath, payload.TargetFileName); err != nil {
		log.Errorf("upload file error: %v", err)
		return err
	}

	subs, err := r.notiSubscriptionPort.FindByUserID(ctx, payload.UserID)
	if err != nil {
		log.Errorf("find user subscription error: %v", err)
		return err
	}

	for _, sub := range subs {
		if _, err := r.notificationPort.SendNoti(ctx, "t2v", "Your video is ready", "Your video is ready", "https://storage.bralyvn.com/sira/assets/Discover/DI01_Pixar_Trump_thumb.jpg", sub.UserToken); err != nil {
			log.Errorf("send notification error: %v", err)
		}
	}

	log.Infoc(ctx, "task %s is done", taskID)
	return nil
}
