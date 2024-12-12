package dir

import (
	"fmt"
	"io"
	"os"
)

//文件目录是否存在

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//批量创建文件夹

func CreateDir(dirs ...string) error {
	for _, v := range dirs {
		exist, err := PathExists(v)
		if err != nil {
			return err
		}
		if !exist {
			err = os.MkdirAll(v, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetFiles(folder string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			GetFiles(folder + "/" + file.Name())
		}
	}
	return files, nil
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

// 文件是否存在
func FileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	// 其他错误
	return false, err
}
