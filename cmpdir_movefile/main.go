package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	defer func() {
		fmt.Scanln()
	}()
	runDir := "./" // 当前运行目录
	execName, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if err := OrganizeFilesByBaseName(runDir, execName); err != nil {
		fmt.Println("出错了:", err)
	} else {
		fmt.Println("同名文件整理完成!")
	}
}

// OrganizeFilesByBaseName 按基础文件名分组并将文件放到相应的目录中
func OrganizeFilesByBaseName(runDir, execPath string) error {
	files := make(map[string][]string)
	execName := filepath.Base(execPath)

	// 遍历运行目录中的所有文件
	err := filepath.Walk(runDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过子目录
		if info.IsDir() || info.Name() == execName {
			return nil
		}

		// 获取文件的基础名称（去掉扩展名）
		baseName := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
		files[baseName] = append(files[baseName], path)
		return nil
	})

	if err != nil {
		return err
	}

	// 为每个基础文件名创建目录并移动文件
	for baseName, filePaths := range files {
		destDir := filepath.Join(runDir, baseName)
		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			return err
		}

		for _, filePath := range filePaths {
			fileName := filepath.Base(filePath)
			destPath := filepath.Join(destDir, fileName)

			if err := MoveFile(filePath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// MoveFile 将文件从源路径移动到目标路径
func MoveFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 将文件内容复制到目标文件
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// 删除源文件
	return os.Remove(src)
}
