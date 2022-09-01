package gui

import (
	"fmt"
	"os/exec"
	"simple-gui/core"

	g "github.com/AllenDang/giu"
)

func buildTableRows() []*g.TableRowWidget {
	downloads := getDownloads()

	rows := make([]*g.TableRowWidget, len(downloads))

	for i, download := range downloads {
		currentId := download.Id

		rows[i] = g.TableRow(
			g.Label(download.FileName),
			g.Label(download.Size),
			g.Condition(download.State == DownloadStateDownloading, g.Layout{
				g.Label("Downloading"),
			}, g.Layout{
				g.Row(
					g.Label("Done"),
					g.Button("Remove").OnClick(func() {
						removeDownload(currentId)
					}),
					g.Button("Open folder").OnClick(func() {
						downloadPath, err := core.GetDownloadPath("")
						if err != nil {
							fmt.Println(err)
						}

						cmd := exec.Command("open", downloadPath)
						cmd.Run()
						_, err = cmd.Output()
						if err != nil {
							fmt.Println(err)
						}
					}),
				),
			}),
		)
	}

	return rows
}

func showErrors() {
	if GetBoxError() != "" {
		g.Msgbox("Error", GetBoxError()).ResultCallback(func(_ g.DialogResult) {
			SetBoxError("")
		})
	}
}

func showMainWindow(mainWindow *g.WindowWidget) {
	builtTableRows := buildTableRows()

	tableRows := []*g.TableRowWidget{g.TableRow(g.Label("Name"), g.Label("Size"), g.Label("State"))}
	tableRows = append(tableRows, builtTableRows...)

	mainWindow.Layout(
		g.Button("Set Link").OnClick(func() {
			showInputWindow()
		}),
		g.Table().Rows(tableRows...),
		g.PrepareMsgbox(),
	)
}
