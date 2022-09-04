package core

import "os"

func Start() error {
	downloadPath, err := GetDownloadPath("")
	if err != nil {
		return err
	}

	_, err = os.Stat(downloadPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(downloadPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	tempPath, err := GetTempPath()
	if err != nil {
		return err
	}

	_, err = os.Stat(tempPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(tempPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
