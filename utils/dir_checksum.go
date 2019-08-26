package utils

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func calcDirChecksum(dir string) (dirHash string) {

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		Fatalf("couldn't list directory %s:\n%v\n", dir, err)
	}

	for _, fi := range fileInfos {
		if (!fi.IsDir()) && strings.HasSuffix(fi.Name(), ".conf") {
			dirHash += fmt.Sprintf("%s %d\n", fi.Name(), fi.ModTime().Unix())
		}
	}
	return dirHash
}
