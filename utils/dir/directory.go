package dir

import (
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
