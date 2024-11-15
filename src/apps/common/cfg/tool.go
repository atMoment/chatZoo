package cfg

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func WriteAnything(fileName string, data interface{}) {
	fp, err := os.Create(fileName)
	if err != nil {
		log.Panicf("WriteAnything  failed, os.create err: %s", err)
	}

	d, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Panicf("WriteAnything  failed, json.marshal err: %s", err)
	}
	fp.Write(d)
	log.Printf("WriteAnything over fileName: %v\n", fileName)
}

func WriteAndAppendAnything(fileName string, data interface{}) {
	fp := getFileAboutAppendOrCreate(fileName)

	d, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Panicf("WriteAnything  failed, json.marshal err: %s", err)
	}
	fp.Write(d)
	log.Printf("WriteAnything over fileName: %v\n", fileName)
}

func ReadAnything(fileName string, data interface{}) {
	buffer, err := os.ReadFile(fileName)
	if err != nil {
		log.Panicf("ReadAnything  failed, os.create err: %s", err)
	}
	if len(buffer) == 0 {
		fmt.Println("file data is nil ", fileName)
		return
	}
	err = json.Unmarshal(buffer, data)
	if err != nil {
		log.Panicf("ReadAnything  failed, json.marshal err: %s", err)
	}
	log.Printf("ReadAnything over")
}

// getFileAboutAppendOrCreate 按照无则创建有则追加, 得到文件描述符
func getFileAboutAppendOrCreate(filePath string) *os.File {
	var fp *os.File
	var err error
	if checkFileIsExit(filePath) {
		fp, err = os.OpenFile(filePath, os.O_APPEND, 0777)
		if err != nil {
			fmt.Println("getFileAboutAppendOrCreate failed ", err)
		}
	} else {
		fp, err = os.Create(filePath)
		if err != nil {
			fmt.Println("getFileAboutAppendOrCreate failed ", err)
		}
	}
	return fp
}

func checkFileIsExit(filename string) bool {
	exist := true
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			exist = false
		}
	}
	return exist
}
