package gui

func removeDownloadFromList(list []GuiDownload, index int) []GuiDownload {
	newList := make([]GuiDownload, 0)

	newList = append(newList, list[:index]...)
	return append(newList, list[index+1:]...)
}
