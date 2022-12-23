// storage is an abstraction to s3 buckets

package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/TurnipXenon/Turnip/internal/models"
)

type storageDynamodDBImpl struct {
	svc *dynamodb.DynamoDB
}

func NewStorageDynamoDB(d *dynamodb.DynamoDB) Storage {
	s := storageDynamodDBImpl{}
	s.svc = d
	return &s
}

func (s *storageDynamodDBImpl) GetHostMap() map[string]models.Host {
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

	// todo: sample how to get object dynamodb
	//fmt.Println("Entering here")
	//input := &dynamodb.GetItemInput{
	//	TableName: aws.String(table),
	//	Key: map[string]*dynamodb.AttributeValue{
	//		"hostCode": {S: aws.String("turnip")},
	//	},
	//}
	//
	//fmt.Println("Entering here")
	//item, err := ddb.GetItem(input)
	//fmt.Println("Entering here!!!!")
	//if err != nil {
	//	print("Sad :(", err.Error())
	//	return nil
	//} else {
	//	fmt.Println("Wah wah wah")
	//	fmt.Println(item.Item)
	//}

	return hostMap
}
