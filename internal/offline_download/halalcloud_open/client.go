package halalcloud_open

import (
	"context"
	"fmt"
	"strings"

	halalcloudopen "github.com/OpenListTeam/OpenList/v4/drivers/halalcloud_open"
	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/errs"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/offline_download/tool"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/OpenListTeam/OpenList/v4/internal/setting"
)

type HalalCloudOpen struct{}

func (*HalalCloudOpen) Name() string {
	return "HalalCloudOpen"
}

func (*HalalCloudOpen) Items() []model.SettingItem {
	return nil
}

func (*HalalCloudOpen) Run(_ *tool.DownloadTask) error {
	return errs.NotSupport
}

func (*HalalCloudOpen) Init() (string, error) {
	return "ok", nil
}

func (*HalalCloudOpen) IsReady() bool {
	tempDir := setting.GetStr(conf.HalalCloudOpenTempDir)
	if tempDir != "" {
		storage, _, err := op.GetStorageAndActualPath(tempDir)
		if err == nil {
			if _, ok := storage.(*halalcloudopen.HalalCloudOpen); ok {
				return true
			}
		}
	}
	for _, storage := range op.GetAllStorages() {
		if _, ok := storage.(*halalcloudopen.HalalCloudOpen); ok {
			return true
		}
	}
	return false
}

func (*HalalCloudOpen) AddURL(args *tool.AddUrlArgs) (string, error) {
	storage, actualPath, err := op.GetStorageAndActualPath(args.TempDir)
	if err != nil {
		return "", err
	}
	driverHalalCloudOpen, ok := storage.(*halalcloudopen.HalalCloudOpen)
	if !ok {
		return "", fmt.Errorf("unsupported storage driver for offline download, only HalalCloudOpen is supported")
	}

	ctx := context.Background()
	if err := op.MakeDir(ctx, storage, actualPath); err != nil {
		return "", err
	}
	parentDir, err := op.GetUnwrap(ctx, storage, actualPath)
	if err != nil {
		return "", err
	}
	return driverHalalCloudOpen.OfflineDownload(ctx, args.Url, parentDir.GetPath())
}

func (*HalalCloudOpen) Remove(task *tool.DownloadTask) error {
	storage, _, err := op.GetStorageAndActualPath(task.TempDir)
	if err != nil {
		return err
	}
	driverHalalCloudOpen, ok := storage.(*halalcloudopen.HalalCloudOpen)
	if !ok {
		return fmt.Errorf("unsupported storage driver for offline download, only HalalCloudOpen is supported")
	}
	return driverHalalCloudOpen.DeleteOfflineTask(context.Background(), task.GID, false)
}

func (*HalalCloudOpen) Status(task *tool.DownloadTask) (*tool.Status, error) {
	storage, _, err := op.GetStorageAndActualPath(task.TempDir)
	if err != nil {
		return nil, err
	}
	driverHalalCloudOpen, ok := storage.(*halalcloudopen.HalalCloudOpen)
	if !ok {
		return nil, fmt.Errorf("unsupported storage driver for offline download, only HalalCloudOpen is supported")
	}

	info, err := driverHalalCloudOpen.OfflineDownloadProcess(context.Background(), task.GID)
	if err != nil {
		return nil, err
	}

	statusStr := fmt.Sprintf("status=%d", info.Status)
	completed := false
	var taskErr error
	switch info.Status {
	case 0:
		statusStr = "queued"
	case 1:
		statusStr = "downloading"
	case 2:
		statusStr = "succeed"
		completed = true
	case 3:
		statusStr = "failed"
		taskErr = fmt.Errorf("offline download failed")
	case 4:
		statusStr = "canceled"
		taskErr = fmt.Errorf("offline download canceled")
	}
	if info.Code != 0 {
		taskErr = fmt.Errorf("offline task error code: %d", info.Code)
	}
	if info.Message != "" {
		msg := strings.TrimSpace(info.Message)
		if taskErr != nil {
			taskErr = fmt.Errorf("%w: %s", taskErr, msg)
		} else {
			statusStr = fmt.Sprintf("%s (%s)", statusStr, msg)
		}
	}
	if info.Progress >= 1 {
		completed = taskErr == nil
	}

	return &tool.Status{
		TotalBytes: info.BytesTotal,
		Progress:   info.Progress,
		Completed:  completed,
		Status:     statusStr,
		Err:        taskErr,
	}, nil
}

var _ tool.Tool = (*HalalCloudOpen)(nil)

func init() {
	tool.Tools.Add(&HalalCloudOpen{})
}
