package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

func addFileToIPFS(f io.Reader) (string, error) {
	tmpFile, err := ioutil.TempFile("", "_ipfs_tmp_")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	w := bufio.NewWriter(tmpFile)
	io.Copy(w, f)
	w.Flush()

	cid, err := exec.Command("ipfs", "add", tmpFile.Name()).Output()
	reg := regexp.MustCompile(`added\s+([[:alnum:]]+)\s+`)
	cidStr := string(reg.FindSubmatch(cid)[1])

	if err != nil {
		return "", err
	}
	return cidStr, nil
}
