package gui

import (
	"fmt"
	"simple-gui/core"
	"simple-gui/helper"

	g "github.com/AllenDang/giu"
	"github.com/google/uuid"
)

func startDownloadClick() {
	downloadId := uuid.NewString()

	fmt.Println("____________STARTING DOWNLOAD________________", downloadId)
	downloadLink := *GetCurrentDownloadLink()

	if !helper.IsValidUrl(downloadLink) {
		SetEnterUrlError("Entered Link is Not a valid url!")
		return
	}
	fmt.Println("here1", downloadId)
	hideInputWindow()

	fmt.Println("here2")
	fileName, size, err := core.GetFileDetails(downloadLink)
	if err != nil {
		SetBoxError(err.Error())
		return
	}

	fmt.Println(fileName, "here3", downloadId)

	addDownload(GuiDownload{
		Id:       downloadId,
		FileName: fileName,
		Size:     helper.IntToFloatString(size),
		State:    DownloadStateDownloading,
	})

	go func(currentLink string, downloadId string) {
		fmt.Println("here4", downloadId)
		err := core.StartDownload(currentLink)
		if err != nil {
			SetBoxError(err.Error())
		}
		fmt.Println("here5", downloadId)
		updateDownloadState(downloadId, DownloadStateDone)
	}(downloadLink, downloadId)

	fmt.Println("here6", downloadId)

	SetEnterUrlError("")
	SetCurrentDownloadLink("")
	fmt.Println("____________ENDING DOWNLOAD________________", downloadId)
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
