package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Download struct {
	Url          string
	TargetPath   string
	Filename     string
	TotalSection int
}

func (d Download) getFileInfo() (*http.Response, error) {
	fmt.Println("Making connection")
	r, err := d.getNewRequest("HEAD")
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func (d Download) DownloadSingleThreaded() error {
	r, err := d.getNewRequest("GET")
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(d.TargetPath, bytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (d Download) Do(size int) error {
	if size == 0 {
		return d.DownloadSingleThreaded()
	}

	fmt.Printf("Downloading %d bytes\n", size)

	var sections = make([][2]int, d.TotalSection)
	eachSize := size / d.TotalSection
	fmt.Printf("Each section is %d bytes\n", eachSize)

	for i := range sections {
		if i == 0 {
			// starting byte of first section
			sections[i][0] = 0
		} else {
			// starting byte of next section
			sections[i][0] = sections[i-1][1] + 1
		}

		if i < d.TotalSection-1 {
			// ending byte of other sections
			sections[i][1] = sections[i][0] + eachSize
		} else {
			// ending byte of last section
			sections[i][1] = size - 1
		}
	}

	wg := sync.WaitGroup{}

	for i, section := range sections {
		wg.Add(1)
		go func(i int, section [2]int) {
			defer wg.Done()
			err := d.downloadSection(i, section[0], section[1])
			if err != nil {
				// todo handle gracefully
				panic(err)
			}
		}(i, section)
	}

	wg.Wait()

	err := d.mergeFiles(sections)
	if err != nil {
		return err
	}

	return nil
}

func (d Download) downloadSection(index int, startByte int, endByte int) error {
	fmt.Printf("Downloading section %d\n", index+1)
	r, err := d.getNewRequest("GET")
	if err != nil {
		return err
	}
	r.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", startByte, endByte))
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("Downloaded %v bytes for section %v: %v\n", resp.Header.Get("Content-Length"), index+1, []int{startByte, endByte})

	// todo stream data directly to file
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	tempPath, err := GetTempPath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/section-%v-%s.tmp", tempPath, index+1, d.Filename), bytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (d Download) getNewRequest(method string) (*http.Request, error) {
	r, err := http.NewRequest(method, d.Url, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("User-Agent", "Dl manager v001")

	return r, nil
}

func (d Download) mergeFiles(sections [][2]int) error {
	f, err := os.OpenFile(d.TargetPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	tempPath, err := GetTempPath()
	if err != nil {
		return err
	}

	for i := range sections {
		fmt.Printf("Merging section %v\n", i+1)
		// todo stream sections to write instead of loading to memory
		b, err := ioutil.ReadFile(fmt.Sprintf("%s/section-%v-%s.tmp", tempPath, i+1, d.Filename))
		if err != nil {
			return err
		}
		n, err := f.Write(b)
		if err != nil {
			return err
		}
		fmt.Printf("Merged %v bytes\n", n)
	}

	fmt.Println("Merging complete! Deleting tmp files...")
	go func() {
		for i := range sections {
			err := os.Remove(fmt.Sprintf("%s/section-%v-%s.tmp", tempPath, i+1, d.Filename))
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	return nil
}

func GetFileDetails(link string) (name string, size int, err error) {
	d := Download{
		Url: link,
	}

	resp, err := d.getFileInfo()
	if err != nil {
		return "", 0, err
	}

	if resp.Header.Get("Content-Length") != "" {
		size, err = strconv.Atoi(resp.Header.Get("Content-Length"))
		if err != nil {
			return "", 0, err
		}
	}

	fileType := strings.Split(resp.Header.Get("Content-Type"), "/")[1]

	if len(strings.Split(fileType, ";")) > 0 {
		fileType = strings.Split(fileType, ";")[0]
	}

	var fileName string

	if doesLinkIncludeFileName(link, fileType) {
		fileName = getLinkLastPart(link)
	} else {
		fileName = getLinkLastPart(link) + "." + fileType

		if len(getLinkLastPart(link)) > 15 {
			currentUnixTime := strconv.Itoa(int(time.Now().UnixMilli()))
			fileName = currentUnixTime + "." + fileType
		}
	}

	filePath, err := GetDownloadPath(fileName)
	if err != nil {
		return "", 0, err
	}

	if _, err := os.Stat(filePath); err == nil {
		currentUnixTime := strconv.Itoa(int(time.Now().UnixMilli()))

		fileName = currentUnixTime + "-" + fileName
	}

	return fileName, size, nil
}

func StartDownload(link string) error {
	startTime := time.Now()

	downloadsState.InProgressDownloads++

	d := Download{
		Url:          link,
		TotalSection: 20,
	}

	fileName, size, err := GetFileDetails(link)
	if err != nil {
		return err
	}

	downloadPath, err := GetDownloadPath(fileName)
	if err != nil {
		return err
	}

	d.TargetPath = downloadPath
	d.Filename = fileName

	err = d.Do(size)
	if err != nil {
		return err
	}

	fmt.Printf("Download took %s\n", time.Since(startTime))

	downloadsState.InProgressDownloads--

	return nil
}
