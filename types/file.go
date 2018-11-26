package types

import (
	"io"
	"os"
)

type File struct {
}

// 判断文件是否存在
func (s *File) IsExist(name string) bool {
	_, err := os.Stat(name)

	return err == nil || os.IsExist(err)
}

// 文件大小
func (s *File) Size(name string) int64 {
	info, _ := os.Stat(name)

	return info.Size()
}

// 删除文件
func (s *File) Delete(name string) error {
	_, err := os.Stat(name)
	if err == nil || os.IsExist(err) {
		err = os.Remove(name)
		if err != nil {
			return err
		}
	}

	return nil
}

// 拷贝文件
func (s *File) Copy(source, dest string) (int64, error) {
	sourceFile, err := os.Open(source)
	if err != nil {
		return 0, err
	}
	defer sourceFile.Close()

	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		return 0, err
	}

	destFile, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, sourceFileInfo.Mode())
	if err != nil {
		return 0, err
	}
	defer destFile.Close()

	return io.Copy(destFile, sourceFile)
}
