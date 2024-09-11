package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	rootDir := "."

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != rootDir {
			if filepath.Dir(path) == rootDir {
				return compressDir(path)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", rootDir, err)
	}
}

// compressDir 压缩指定目录
func compressDir(dirPath string) error {
	zipFilePath := dirPath + ".zip"
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径
		relPath, _ := filepath.Rel(filepath.Dir(dirPath), path)
		if info.IsDir() {
			return nil
		}

		// 创建文件头
		zipHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		zipHeader.Name = filepath.Join(filepath.Base(dirPath), relPath)
		zipHeader.Method = zip.Deflate

		zipFileWriter, err := zipWriter.CreateHeader(zipHeader)
		if err != nil {
			return err
		}

		// 打开源文件
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// 复制文件内容到 zip
		_, err = io.Copy(zipFileWriter, srcFile)
		return err
	})
	if err != nil {
		return err
	}

	fmt.Printf("Compressed %s into %s\n", dirPath, zipFilePath)
	return nil
}
