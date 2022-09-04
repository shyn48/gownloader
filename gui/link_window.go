package gui

import (
	"simple-gui/core"
	"simple-gui/helper"

	g "github.com/AllenDang/giu"
	"github.com/google/uuid"
)

func startDownloadClick() {
	downloadId := uuid.NewString()

	downloadLink := *GetCurrentDownloadLink()

	if !helper.IsValidUrl(downloadLink) {
		SetEnterUrlError("Entered Link is Not a valid url!")
		return
	}
	hideInputWindow()

	fileName, size, err := core.GetFileDetails(downloadLink)
	if err != nil {
		SetBoxError(err.Error())
		return
	}

	addDownload(GuiDownload{
		Id:       downloadId,
		FileName: fileName,
		Size:     helper.IntToFloatString(size),
		State:    DownloadStateDownloading,
	})

	go func(currentLink string, downloadId string) {
		err := core.StartDownload(currentLink)
		if err != nil {
			SetBoxError(err.Error())
		}
		updateDownloadState(downloadId, DownloadStateDone)
	}(downloadLink, downloadId)

	SetEnterUrlError("")
	SetCurrentDownloadLink("")
}

func showLinkWindow(linkWindow *g.WindowWidget, mainWindow *g.WindowWidget) {
	if inputWindowShown {
		if mainWindow.HasFocus() {
			linkWindow.BringToFront()
		}

		linkWindow.Pos(150, 200).IsOpen(&inputWindowShown).Size(500, 150).Layout(
			g.Label("Please enter download link"),
			g.InputText(GetCurrentDownloadLink()).Hint("download link"),
			g.Label(GetEnterUrlError()),
			g.Button("Start Download").OnClick(startDownloadClick),
		)
	}
}
