package handles

import (
	"regexp"
	"strings"

	_115 "github.com/OpenListTeam/OpenList/v4/drivers/115"
	_115_open "github.com/OpenListTeam/OpenList/v4/drivers/115_open"
	_123 "github.com/OpenListTeam/OpenList/v4/drivers/123"
	_123_open "github.com/OpenListTeam/OpenList/v4/drivers/123_open"
	halalcloudopen "github.com/OpenListTeam/OpenList/v4/drivers/halalcloud_open"
	"github.com/OpenListTeam/OpenList/v4/drivers/pikpak"
	"github.com/OpenListTeam/OpenList/v4/drivers/thunder"
	"github.com/OpenListTeam/OpenList/v4/drivers/thunder_browser"
	"github.com/OpenListTeam/OpenList/v4/drivers/thunderx"
	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/offline_download/tool"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/OpenListTeam/OpenList/v4/internal/task"
	"github.com/OpenListTeam/OpenList/v4/server/common"
	"github.com/gin-gonic/gin"
)

type SetAria2Req struct {
	Uri    string `json:"uri" form:"uri"`
	Secret string `json:"secret" form:"secret"`
}

func SetAria2(c *gin.Context) {
	var req SetAria2Req
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	items := []model.SettingItem{
		{Key: conf.Aria2Uri, Value: req.Uri, Type: conf.TypeString, Group: model.OFFLINE_DOWNLOAD, Flag: model.PRIVATE},
		{Key: conf.Aria2Secret, Value: req.Secret, Type: conf.TypeString, Group: model.OFFLINE_DOWNLOAD, Flag: model.PRIVATE},
	}
	if err := op.SaveSettingItems(items); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	_tool, err := tool.Tools.Get("aria2")
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	version, err := _tool.Init()
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c, version)
}

type SetQbittorrentReq struct {
	Url      string `json:"url" form:"url"`
	Seedtime string `json:"seedtime" form:"seedtime"`
}

func SetQbittorrent(c *gin.Context) {
	var req SetQbittorrentReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	items := []model.SettingItem{
		{Key: conf.QbittorrentUrl, Value: req.Url, Type: conf.TypeString, Group: model.OFFLINE_DOWNLOAD, Flag: model.PRIVATE},
		{Key: conf.QbittorrentSeedtime, Value: req.Seedtime, Type: conf.TypeNumber, Group: model.OFFLINE_DOWNLOAD, Flag: model.PRIVATE},
	}
	if err := op.SaveSettingItems(items); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	_tool, err := tool.Tools.Get("qBittorrent")
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	if _, err := _tool.Init(); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c, "ok")
}

type SetTransmissionReq struct {
	Uri      string `json:"uri" form:"uri"`
	Seedtime string `json:"seedtime" form:"seedtime"`
}

func SetTransmission(c *gin.Context) {
	var req SetTransmissionReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	items := []model.SettingItem{
		{Key: conf.TransmissionUri, Value: req.Uri, Type: conf.TypeString, Group: model.OFFLINE_DOWNLOAD, Flag: model.PRIVATE},
		{Key: conf.TransmissionSeedtime, Value: req.Seedtime, Type: conf.TypeNumber, Group: model.OFFLINE_DOWNLOAD, Flag: model.PRIVATE},
	}
	if err := op.SaveSettingItems(items); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	_tool, err := tool.Tools.Get("Transmission")
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	if _, err := _tool.Init(); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c, "ok")
}

type Set115Req struct {
	TempDir string `json:"temp_dir" form:"temp_dir"`
}

func Set115(c *gin.Context) {
	var req Set115Req
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if req.TempDir != "" {
		storage, _, err := op.GetStorageAndActualPath(req.TempDir)
		if err != nil {
			common.ErrorStrResp(c, "storage does not exists", 400)
			return
		}
		if storage.Config().CheckStatus && storage.GetStorage().Status != op.WORK {
			common.ErrorStrResp(c, "storage not init: "+storage.GetStorage().Status, 400)
			return
		}
		if _, ok := storage.(*_115.Pan115); !ok {
			common.ErrorStrResp(c, "unsupported storage driver for offline download, only 115 Cloud is supported", 400)
			return
		}
	}
	items := []model.SettingItem{
		{Key: conf.Pan115TempDir, Value: req.TempDir, Type: conf.TypeString, Group: model.OFFLINE_DOWNLOAD, Flag: model.PRIVATE},
	}
	if err := op.SaveSettingItems(items); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	_tool, err := tool.Tools.Get("115 Cloud")
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	if _, err := _tool.Init(); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c, "ok")
}

type Set115OpenReq struct {
	TempDir string `json:"temp_dir" form:"temp_dir"`
}

func Set115Open(c *gin.Context) {
	var req Set115OpenReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if req.TempDir != "" {
		storage, _, err := op.GetStorageAndActualPath(req.TempDir)
		if err != nil {
			common.ErrorStrResp(c, "storage does not exists", 400)
			return
		}
		if storage.Config().CheckStatus && storage.GetStorage().Status != op.WORK {
			common.ErrorStrResp(c, "storage not init: "+storage.GetStorage().Status, 400)
			return
		}
		if _, ok := storage.(*_115_open.Open115); !ok {
			common.ErrorStrResp(c, "unsupported storage driver for offline download, only 115 Open is supported", 400)
			return
		}
	}
	items := []model.SettingItem{
		{Key: conf.Pan115OpenTempDir, Value: req.TempDir, Type: conf.TypeString, Group: model.OFFLINE_DOWNLOAD, Flag: model.PRIVATE},
	}
	if err := op.SaveSettingItems(items); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	_tool, err := tool.Tools.Get("115 Open")
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	if _, err := _tool.Init(); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c, "ok")
}

type Set123PanReq struct {
	TempDir string `json:"temp_dir" form:"temp_dir"`
}

func Set123Pan(c *gin.Context) {
	var req Set123PanReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if req.TempDir != "" {
		storage, _, err := op.GetStorageAndActualPath(req.TempDir)
		if err != nil {
			common.ErrorStrResp(c, "storage does not exists", 400)
			return
		}
		if storage.Config().CheckStatus && storage.GetStorage().Status != op.WORK {
			common.ErrorStrResp(c, "storage not init: "+storage.GetStorage().Status, 400)
			return
		}
		if _, ok := storage.(*_123.Pan123); !ok {
			common.ErrorStrResp(c, "unsupported storage driver for offline download, only 123Pan is supported", 400)
			return
		}
	}
	items := []model.SettingItem{
		{Key: conf.Pan123TempDir, Value: req.TempDir, Type: conf.TypeString, Group: model.OFFLINE_DOWNLOAD, Flag: model.PRIVATE},
	}
	if err := op.SaveSettingItems(items); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	_tool, err := tool.Tools.Get("123Pan")
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	if _, err := _tool.Init(); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c, "ok")
}

type Set123OpenReq struct {
	TempDir     string `json:"temp_dir" form:"temp_dir"`
	CallbackUrl string `json:"callback_url" form:"callback_url"`
}

func Set123Open(c *gin.Context) {
	var req Set123OpenReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if req.TempDir != "" {
		storage, _, err := op.GetStorageAndActualPath(req.TempDir)
		if err != nil {
			common.ErrorStrResp(c, "storage does not exists", 400)
			return
		}
		if storage.Config().CheckStatus && storage.GetStorage().Status != op.WORK {
			common.ErrorStrResp(c, "storage not init: "+storage.GetStorage().Status, 400)
			return
		}
		if _, ok := storage.(*_123_open.Open123); !ok {
			common.ErrorStrResp(c, "unsupported storage driver for offline download, only 123 Open is supported", 400)
			return
		}
	}
	items := []model.SettingItem{
		{Key: conf.Pan123OpenTempDir, Value: req.TempDir, Type: conf.TypeString, Group: model.OFFLINE_DOWNLOAD, Flag: model.PRIVATE},
		{Key: conf.Pan123OpenOfflineDownloadCallbackUrl, Value: req.CallbackUrl, Type: conf.TypeString, Group: model.OFFLINE_DOWNLOAD, Flag: model.PRIVATE},
	}
	if err := op.SaveSettingItems(items); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	_tool, err := tool.Tools.Get("123 Open")
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	if _, err := _tool.Init(); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c, "ok")
}


type SetHalalCloudOpenReq struct {
	TempDir string `json:"temp_dir" form:"temp_dir"`
}

func SetHalalCloudOpen(c *gin.Context) {
	var req SetHalalCloudOpenReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if req.TempDir != "" {
		storage, _, err := op.GetStorageAndActualPath(req.TempDir)
		if err != nil {
			common.ErrorStrResp(c, "storage does not exists", 400)
			return
		}
		if storage.Config().CheckStatus && storage.GetStorage().Status != op.WORK {
			common.ErrorStrResp(c, "storage not init: "+storage.GetStorage().Status, 400)
			return
		}
		if _, ok := storage.(*halalcloudopen.HalalCloudOpen); !ok {
			common.ErrorStrResp(c, "unsupported storage driver for offline download, only HalalCloudOpen is supported", 400)
			return
		}
	}
	items := []model.SettingItem{
		{Key: conf.HalalCloudOpenTempDir, Value: req.TempDir, Type: conf.TypeString, Group: model.OFFLINE_DOWNLOAD, Flag: model.PRIVATE},
	}
	if err := op.SaveSettingItems(items); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	_tool, err := tool.Tools.Get("HalalCloudOpen")
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	if _, err := _tool.Init(); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c, "ok")
}
type SetPikPakReq struct {
	TempDir string `json:"temp_dir" form:"temp_dir"`
}

func SetPikPak(c *gin.Context) {
	var req SetPikPakReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if req.TempDir != "" {
		storage, _, err := op.GetStorageAndActualPath(req.TempDir)
		if err != nil {
			common.ErrorStrResp(c, "storage does not exists", 400)
			return
		}
		if storage.Config().CheckStatus && storage.GetStorage().Status != op.WORK {
			common.ErrorStrResp(c, "storage not init: "+storage.GetStorage().Status, 400)
			return
		}
		if _, ok := storage.(*pikpak.PikPak); !ok {
			common.ErrorStrResp(c, "unsupported storage driver for offline download, only PikPak is supported", 400)
			return
		}
	}
	items := []model.SettingItem{
		{Key: conf.PikPakTempDir, Value: req.TempDir, Type: conf.TypeString, Group: model.OFFLINE_DOWNLOAD, Flag: model.PRIVATE},
	}
	if err := op.SaveSettingItems(items); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	_tool, err := tool.Tools.Get("PikPak")
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	if _, err := _tool.Init(); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c, "ok")
}

type SetThunderReq struct {
	TempDir string `json:"temp_dir" form:"temp_dir"`
}

func SetThunder(c *gin.Context) {
	var req SetThunderReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if req.TempDir != "" {
		storage, _, err := op.GetStorageAndActualPath(req.TempDir)
		if err != nil {
			common.ErrorStrResp(c, "storage does not exists", 400)
			return
		}
		if storage.Config().CheckStatus && storage.GetStorage().Status != op.WORK {
			common.ErrorStrResp(c, "storage not init: "+storage.GetStorage().Status, 400)
			return
		}
		if _, ok := storage.(*thunder.Thunder); !ok {
			common.ErrorStrResp(c, "unsupported storage driver for offline download, only Thunder is supported", 400)
			return
		}
	}
	items := []model.SettingItem{
		{Key: conf.ThunderTempDir, Value: req.TempDir, Type: conf.TypeString, Group: model.OFFLINE_DOWNLOAD, Flag: model.PRIVATE},
	}
	if err := op.SaveSettingItems(items); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	_tool, err := tool.Tools.Get("Thunder")
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	if _, err := _tool.Init(); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c, "ok")
}

type SetThunderXReq struct {
	TempDir string `json:"temp_dir" form:"temp_dir"`
}

func SetThunderX(c *gin.Context) {
	var req SetThunderXReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if req.TempDir != "" {
		storage, _, err := op.GetStorageAndActualPath(req.TempDir)
		if err != nil {
			common.ErrorStrResp(c, "storage does not exists", 400)
			return
		}
		if storage.Config().CheckStatus && storage.GetStorage().Status != op.WORK {
			common.ErrorStrResp(c, "storage not init: "+storage.GetStorage().Status, 400)
			return
		}
		if _, ok := storage.(*thunderx.ThunderX); !ok {
			common.ErrorStrResp(c, "unsupported storage driver for offline download, only ThunderX is supported", 400)
			return
		}
	}
	items := []model.SettingItem{
		{Key: conf.ThunderXTempDir, Value: req.TempDir, Type: conf.TypeString, Group: model.OFFLINE_DOWNLOAD, Flag: model.PRIVATE},
	}
	if err := op.SaveSettingItems(items); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	_tool, err := tool.Tools.Get("ThunderX")
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	if _, err := _tool.Init(); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c, "ok")
}

type SetThunderBrowserReq struct {
	TempDir string `json:"temp_dir" form:"temp_dir"`
}

func SetThunderBrowser(c *gin.Context) {
	var req SetThunderBrowserReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if req.TempDir != "" {
		storage, _, err := op.GetStorageAndActualPath(req.TempDir)
		if err != nil {
			common.ErrorStrResp(c, "storage does not exists", 400)
			return
		}
		if storage.Config().CheckStatus && storage.GetStorage().Status != op.WORK {
			common.ErrorStrResp(c, "storage not init: "+storage.GetStorage().Status, 400)
			return
		}
		switch storage.(type) {
		case *thunder_browser.ThunderBrowser, *thunder_browser.ThunderBrowserExpert:
		default:
			common.ErrorStrResp(c, "unsupported storage driver for offline download, only ThunderBrowser is supported", 400)
		}
	}
	items := []model.SettingItem{
		{Key: conf.ThunderBrowserTempDir, Value: req.TempDir, Type: conf.TypeString, Group: model.OFFLINE_DOWNLOAD, Flag: model.PRIVATE},
	}
	if err := op.SaveSettingItems(items); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	_tool, err := tool.Tools.Get("ThunderBrowser")
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	if _, err := _tool.Init(); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c, "ok")
}

func OfflineDownloadTools(c *gin.Context) {
	tools := tool.Tools.Names()
	common.SuccessResp(c, tools)
}

var magnetLinkRegex = regexp.MustCompile(`(?i)magnet:\?[^\s"'<>]+`)

type AddOfflineDownloadReq struct {
	Urls         []string `json:"urls"`
	RawText      string   `json:"raw_text"`
	Path         string   `json:"path"`
	Tool         string   `json:"tool"`
	DeletePolicy string   `json:"delete_policy"`
}

func collectOfflineURLs(req AddOfflineDownloadReq) []string {
	seen := make(map[string]struct{})
	urls := make([]string, 0)
	add := func(raw string) {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			return
		}
		matches := magnetLinkRegex.FindAllString(raw, -1)
		if len(matches) == 0 {
			cleaned := cleanOfflineURL(raw)
			if cleaned == "" {
				return
			}
			if _, ok := seen[cleaned]; ok {
				return
			}
			seen[cleaned] = struct{}{}
			urls = append(urls, cleaned)
			return
		}
		for _, match := range matches {
			cleaned := cleanOfflineURL(match)
			if cleaned == "" {
				continue
			}
			if _, ok := seen[cleaned]; ok {
				continue
			}
			seen[cleaned] = struct{}{}
			urls = append(urls, cleaned)
		}
	}
	for _, u := range req.Urls {
		add(u)
	}
	add(req.RawText)
	return urls
}

func cleanOfflineURL(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	trimmed = strings.Trim(trimmed, "\"'¡°¡±¡®¡¯()[]{}<>£¬¡£,.;£»!£¡?£¿")
	return strings.TrimSpace(trimmed)
}

func AddOfflineDownload(c *gin.Context) {
	user := c.Request.Context().Value(conf.UserKey).(*model.User)
	if !user.CanAddOfflineDownloadTasks() {
		common.ErrorStrResp(c, "permission denied", 403)
		return
	}

	var req AddOfflineDownloadReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	reqPath, err := user.JoinPath(req.Path)
	if err != nil {
		common.ErrorResp(c, err, 403)
		return
	}

	urls := collectOfflineURLs(req)
	if len(urls) == 0 {
		common.ErrorStrResp(c, "no valid download url found", 400)
		return
	}

	var tasks []task.TaskExtensionInfo
	for _, u := range urls {
		t, err := tool.AddURL(c, &tool.AddURLArgs{
			URL:          u,
			DstDirPath:   reqPath,
			Tool:         req.Tool,
			DeletePolicy: tool.DeletePolicy(req.DeletePolicy),
		})
		if err != nil {
			common.ErrorResp(c, err, 500)
			return
		}
		if t != nil {
			tasks = append(tasks, t)
		}
	}
	common.SuccessResp(c, gin.H{
		"tasks": getTaskInfos(tasks),
	})
}
