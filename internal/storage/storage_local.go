// storage is an abstraction to s3 buckets

package storage

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/TurnipXenon/turnip/internal/models"
)

type storageLocalImpl struct{}

func NewStorageLocal() Storage {
	return &storageLocalImpl{}
}

func (s *storageLocalImpl) GetHostMap() map[string]models.Host {
	sourceFile, err := os.Open("./configs/host_local.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer func(sourceFile *os.File) {
		_ = sourceFile.Close()
	}(sourceFile) // ok to ignore error: file was opened read-only.

	byteValue, err := ioutil.ReadAll(sourceFile)
	if err != nil {
		log.Fatalln(err)
	}

	var hostList []models.HostImpl
	err = json.Unmarshal(byteValue, &hostList)
	if err != nil {
		log.Fatalln(err)
	}

	hostMap := map[string]models.Host{}
	for index, host := range hostList {
		for _, alias := range host.GetAliasList() {
			hostMap[alias] = &hostList[index]
		}
	}

	return hostMap
}
