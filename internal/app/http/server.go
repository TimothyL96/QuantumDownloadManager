package http

import (
	"net/http"

	"github.com/ttimt/QuantumDownloadManager/internal/app/downloader/manager"
	"github.com/ttimt/QuantumDownloadManager/internal/app/downloader/user"
)

var (
	downloader  manager.Download
	userSetting user.Setting
)

func Get(w http.ResponseWriter, r *http.Request) {

}
