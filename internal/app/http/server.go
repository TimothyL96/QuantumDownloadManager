package http

import (
	"net/http"

	"github.com/ttimt/QuantumDownloadManager/internal/app/downloader/manager"
	"github.com/ttimt/QuantumDownloadManager/internal/app/downloader/user/setting"
)

var (
	downloader  manager.Download
	userSetting setting.Setting
)

func Get(w http.ResponseWriter, r *http.Request) {

}
