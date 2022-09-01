package gui

var (
	inputWindowShown = false
	enterUrlError    = ""
	boxError         = ""
	currentDownloads = []GuiDownload{}
)

type DownloadState string

const (
	DownloadStateDownloading DownloadState = "DOWNLOADING"
	DownloadStateDone        DownloadState = "DONE"
)

type GuiDownload struct {
	Id       string
	FileName string
	State    DownloadState
	Size     string
}

func removeDownload(id string) {
	for i, download := range currentDownloads {
		if download.Id == id {
			currentDownloads = removeDownloadFromList(currentDownloads, i)
		}
	}
}

func updateDownloadState(id string, newState DownloadState) {
	for i := range currentDownloads {
		if currentDownloads[i].Id == id {
			currentDownloads[i].State = DownloadStateDone
		}
	}
}

func addDownload(download GuiDownload) {
	currentDownloads = append(currentDownloads, download)
}

func getDownloads() []GuiDownload {
	return currentDownloads
}

func showInputWindow() {
	inputWindowShown = true
}

func hideInputWindow() {
	inputWindowShown = false
}

var (
	currentDownloadLink = ""
)

func GetCurrentDownloadLink() *string {
	return &currentDownloadLink
}

func SetCurrentDownloadLink(value string) {
	currentDownloadLink = value
}

func SetEnterUrlError(value string) {
	enterUrlError = value
}

func GetEnterUrlError() string {
	return enterUrlError
}

func SetBoxError(value string) {
	boxError = value
}

func GetBoxError() string {
	return boxError
}
