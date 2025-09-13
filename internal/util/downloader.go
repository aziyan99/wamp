package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type WriteCounter struct {
	Total uint64
}

func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading... %s complete", Bytes(wc.Total))
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func DownloadFile(filepath string, url string) error {

	_, err := os.Stat(filepath)
	if err == nil {
		fmt.Printf("File %s already exists\n", filepath)

		// For now we assume the previous download is not corrupt just we just skip it
		return nil
	}

	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	response, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer response.Body.Close()

	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(response.Body, counter)); err != nil {
		out.Close()
		return err
	}

	fmt.Print("\n")

	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}

	return nil
}
