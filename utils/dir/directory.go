package dir

import (
	"io/fs"
	"io/ioutil"
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

func CreateDir(dirs ...string) (err error) {
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
	return err
}

func GetFiles(folder string) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(folder)
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
