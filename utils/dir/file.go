package dir

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func SaveFile(fileName string, UTXOInfo any) {
	// 保存到文件
	// 将结构体编码为JSON
	jsonData, err := json.Marshal(UTXOInfo)
	if err != nil {
		log.Fatalf("Error marshaling JSON:%v", err)
		return
	}
	// 将JSON数据写入文件
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Error creating file:%v", err)
		return
	}
	defer file.Close()
	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing file:%v", err)
		return
	}
}

// 读取文件内容
func GetFileContent(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 读取文件内容
	return io.ReadAll(file)
}

func ReadFiles(folder string) ([][]byte, error) {
	dir, err := os.Open(folder)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	// 读取目录下的文件和子目录
	files, err := dir.Readdir(-1) // -1 表示读取所有条目
	if err != nil {
		return nil, err
	}
	var fileList [][]byte
	for _, file := range files {
		if file.IsDir() {
			continue // 忽略子目录
		}
		filePath := fmt.Sprintf("%s/%s", folder, file.Name())
		content, err := GetFileContent(filePath)
		if err != nil {
			continue
		}
		fileList = append(fileList, content)
	}
	return fileList, nil
}
