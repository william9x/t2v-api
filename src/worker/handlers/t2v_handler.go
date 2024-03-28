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
	"strings"
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
	taskID := task.ResultWriter().TaskID()
	log.Infoc(ctx, "task %s is processing", taskID)
	defer log.Infoc(ctx, "task %s is done", taskID)

	var payload entities.InferPayload
	if err := msgpack.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unpack task failed: %v", err)
	}
	log.Debugc(ctx, "task payload: %+v", payload)

	if strings.Contains(payload.Prompt, "error") {
		go r.sendNoti(payload.UserID, payload.Agent, taskID, "", false)
		return fmt.Errorf("debu only: task %s is failed", taskID)
	}

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
	}

	result, err := r.inferencePort.Infer(ctx, cmd)
	if err != nil {
		go r.sendNoti(payload.UserID, payload.Agent, taskID, "", false)
		return err
	}
	log.Infoc(ctx, "task %s inference completed, start uploading file at %s", taskID, payload.TargetFileName)

	if err := r.objectStoragePort.UploadFilePath(ctx, result.VideoPath, payload.TargetFileName); err != nil {
		log.Errorf("upload video file error: %v", err)
		go r.sendNoti(payload.UserID, payload.Agent, taskID, "", false)
		return err
	}

	thumbnail := strings.ReplaceAll(payload.TargetFileName, ".mp4", "_thumbnail.jpg")
	go r.sendNoti(payload.UserID, payload.Agent, taskID, thumbnail, true)
	return nil
}

func (r *T2VHandler) sendNoti(userID, agent, taskID, imageURL string, success bool) {
	ctx := context.Background()
	subs, err := r.notiSubscriptionPort.FindByUserID(ctx, userID)
	if err != nil {
		log.Errorf("find user subscription error: %v", err)
	}
	log.Debugc(ctx, "found %d subscriptions for user %s", len(subs), userID)

	title := "🎬 Production Complete!"
	body := "Your video is ready for the spotlight. Watch your creation!"
	image, _ := r.objectStoragePort.GetPreSignedObject(ctx, imageURL)
	if !success {
		title = "🔄 Ooops! Take Another Shot!"
		body = "We're unable to create your video. Go ahead and give it another go!"
		image = ""
	}

	for _, sub := range subs {
		sentId, err := r.notificationPort.SendNoti(
			ctx,
			agent,
			taskID,
			title,
			body,
			image,
			sub.UserToken,
		)
		if err == nil {
			log.Debugc(ctx, "notification sent: %s", sentId)
			return
		}

		log.Errorf("send notification error: %v", err)
		if !strings.Contains(fmt.Sprintf("%v", err), "registration-token-not-registered") {
			return
		}

		log.Debugc(ctx, "not-registered, removing subscription for user %s token %s", userID, sub.UserToken)
		if err := r.notiSubscriptionPort.Delete(ctx, userID, sub.UserToken); err != nil {
			log.Errorf("delete subscription error: %v", err)
		}
	}
}
