package main

import (
	"archive/zip"
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
	// 目标文件夹
	rootDir := "."
	// 输出的ZIP文件存放目录
	outputDir := "zip_output"
	// 创建输出目录
	err := os.MkdirAll(filepath.Join(".", outputDir), os.ModePerm)
	if err != nil {
		fmt.Printf("创建输出目录失败: %v\n", err)
		return
	}
	// 获取当前运行的可执行文件名
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("运行目录,获取失败: %v\n", err)
		return
	}
	currentExecutable := strings.TrimSuffix(filepath.Base(execPath), filepath.Ext(execPath))

	if err := compressFilesInRoot(rootDir, outputDir, currentExecutable); err != nil {
		fmt.Printf("压缩出错: %v\n", err)
	} else {
		fmt.Println("压缩完成")
	}
}

func compressFilesInRoot(rootDir string, outputDir string, excludeFileName string) error {
	// 存储 ZIP Writer，按基础文件名区分
	processedFiles := make(map[string]*zip.Writer)
	// 读取根目录下的文件和文件夹
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return fmt.Errorf("读取目录(%s)失败: %v", rootDir, err)
	}

	for _, entry := range entries {
		// 获取文件名的基础部分（去除扩展名）
		baseName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		if baseName == excludeFileName {
			// 排除当前运行的可执行文件
			continue
		}
		if entry.Name() == outputDir {
			// 排除输出目录
			continue
		}

		zipFilePath := filepath.Join(outputDir, baseName+".zip")
		// 获取或创建 ZIP Writer
		zipWriter, exists := processedFiles[baseName]
		if !exists {
			zipFile, err := os.Create(zipFilePath)
			if err != nil {
				return fmt.Errorf("创建zip文件失败: %v", err)
			}
			defer zipFile.Close()

			zipWriter = zip.NewWriter(zipFile)
			processedFiles[baseName] = zipWriter
			defer zipWriter.Close() // 确保关闭 ZIP Writer
		}

		fullPath := filepath.Join(rootDir, entry.Name())
		if entry.IsDir() {
			if err := compressDir(zipWriter, fullPath, zipFilePath); err != nil {
				return fmt.Errorf("失败压缩 %s 进 %s : %v\n", fullPath, zipFilePath, err)
			}
			fmt.Printf("成功压缩 %s 进 %s\n", fullPath, zipFilePath)

		} else {
			if err := compressFile(zipWriter, fullPath, zipFilePath); err != nil {
				return fmt.Errorf("失败压缩 %s 进 %s : %v\n", fullPath, zipFilePath, err)
			}
			fmt.Printf("成功压缩 %s 进 %s\n", fullPath, zipFilePath)
		}
	}
	return nil
}

// compressDir 压缩指定目录
func compressDir(zipWriter *zip.Writer, dirPath string, zipFilePath string) error {
	// 遍历文件夹并添加到zip中
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 跳过当前文件夹
		if path == dirPath {
			return nil
		}

		// 计算相对路径
		relPath, err := filepath.Rel(filepath.Dir(dirPath), path)
		if err != nil {
			return fmt.Errorf("计算相对路径(%s,%s): %v", dirPath, path, err)
		}

		if info.IsDir() {
			return nil
		}
		return addFileToZip(zipWriter, path, relPath)

	})
	if err != nil {
		return err
	}

	return nil
}

func compressFile(zipWriter *zip.Writer, filePath string, zipFilePath string) error {
	relPath := filepath.Base(filePath)

	return addFileToZip(zipWriter, filePath, relPath)
}

// 降单个文件添加到zip文件中
func addFileToZip(zipWriter *zip.Writer, filePath string, relPath string) error {
	srcFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开文件(%s)失败: %v", filePath, err)
	}
	defer srcFile.Close()

	zipFileWriter, err := zipWriter.Create(relPath)
	if err != nil {
		return fmt.Errorf("创建zip条目(%s)失败: %v", relPath, err)
	}

	// 复制文件内容到zip条目
	_, err = io.Copy(zipFileWriter, srcFile)
	if err != nil {
		return fmt.Errorf("复制文件(%s)内容到zip失败: %v", filePath, err)
	}

	return nil
}
