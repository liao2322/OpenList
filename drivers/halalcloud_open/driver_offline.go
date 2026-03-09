package halalcloudopen

import (
	"context"
	"fmt"
	"strings"

	sdkModel "github.com/halalcloud/golang-sdk-lite/halalcloud/model"
	sdkOffline "github.com/halalcloud/golang-sdk-lite/halalcloud/services/offline"
)

type OfflineTaskInfo struct {
	Status         int32
	Code           int32
	Message        string
	Progress       float64
	BytesTotal     int64
	BytesProcessed int64
}

func (d *HalalCloudOpen) OfflineDownload(ctx context.Context, url string, dirPath string) (string, error) {
	if d.sdkOfflineService == nil {
		return "", fmt.Errorf("offline service is not initialized")
	}
	resp, err := d.sdkOfflineService.Add(ctx, &sdkOffline.UserTask{
		Url:      strings.TrimSpace(url),
		SavePath: dirPath,
	})
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", fmt.Errorf("empty offline task response")
	}
	if resp.TaskIdentity != "" {
		return resp.TaskIdentity, nil
	}
	if resp.Identity != "" {
		return resp.Identity, nil
	}
	return "", fmt.Errorf("empty task id in offline task response")
}

func (d *HalalCloudOpen) OfflineDownloadProcess(ctx context.Context, taskID string) (*OfflineTaskInfo, error) {
	task, err := d.getOfflineTask(ctx, taskID)
	if err != nil {
		return nil, err
	}
	progress := float64(task.Progress)
	if progress <= 0 && task.BytesTotal > 0 {
		progress = float64(task.BytesProcessed) / float64(task.BytesTotal) * 100
	}
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}
	return &OfflineTaskInfo{
		Status:         task.Status,
		Code:           task.Code,
		Message:        task.Message,
		Progress:       progress / 100,
		BytesTotal:     task.BytesTotal,
		BytesProcessed: task.BytesProcessed,
	}, nil
}

func (d *HalalCloudOpen) DeleteOfflineTask(ctx context.Context, taskID string, deleteFiles bool) error {
	if d.sdkOfflineService == nil {
		return fmt.Errorf("offline service is not initialized")
	}
	task, err := d.getOfflineTask(ctx, taskID)
	if err != nil {
		return err
	}
	identity := task.Identity
	if identity == "" {
		identity = taskID
	}
	_, err = d.sdkOfflineService.Delete(ctx, &sdkOffline.OfflineTaskDeleteRequest{
		Identity:    []string{identity},
		DeleteFiles: deleteFiles,
	})
	return err
}

func (d *HalalCloudOpen) getOfflineTask(ctx context.Context, taskID string) (*sdkOffline.UserTask, error) {
	if d.sdkOfflineService == nil {
		return nil, fmt.Errorf("offline service is not initialized")
	}
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, fmt.Errorf("task id is empty")
	}
	token := ""
	for i := 0; i < 20; i++ {
		resp, err := d.sdkOfflineService.List(ctx, &sdkOffline.OfflineTaskListRequest{
			ListInfo: &sdkModel.ScanListRequest{Token: token, Limit: 200},
		})
		if err != nil {
			return nil, err
		}
		if resp == nil || len(resp.Tasks) == 0 {
			break
		}
		for _, item := range resp.Tasks {
			if item == nil {
				continue
			}
			if item.TaskIdentity == taskID || item.Identity == taskID {
				return item, nil
			}
		}
		if resp.ListInfo == nil || resp.ListInfo.Token == "" || resp.ListInfo.Token == token {
			break
		}
		token = resp.ListInfo.Token
	}
	return nil, fmt.Errorf("offline task not found: %s", taskID)
}
