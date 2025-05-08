/*
 *  * Copyright (c) 2025 @umfaka
 *  *
 *  * This program is free software: you can redistribute it and/or modify
 *  * it under the terms of the GNU Affero General Public License as
 *  * published by the Free Software Foundation, either version 3 of the
 *  * License, or (at your option) any later version.
 *  *
 *  * 本程序为自由软件：你可以在遵守 GNU Affero 通用公共许可证（AGPL）第3版或（由你选择）任何更高版本的前提下重新发布和/或修改本程序。
 *  *
 *  * This program is distributed in the hope that it will be useful,
 *  * but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 *  * GNU Affero General Public License for more details.
 *  *
 *  * 本程序以“按现状”提供，不提供任何明示或暗示的保证，包括但不限于适销性或特定用途适用性的保证。
 *  * 有关详细信息，请参阅 GNU Affero 通用公共许可证。
 *  *
 *  * You should have received a copy of the GNU Affero General Public License
 *  * along with this program. If not, see <https://www.gnu.org/licenses/>.
 *  *
 *  * 您应该已经收到一份 GNU Affero 通用公共许可证的副本；如果没有，请访问：<https://www.gnu.org/licenses/>
 *  *
 *  * Team:   @umfaka
 *  * Author: @uncledawen
 *  * Project: cmpdir
 *  * File: main.go
 *  * LastModified: 2025/5/8 15:34
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// 常用国家区号列表（ITU 国际标准）
var countryCodes = []string{
	"1", "7", "20", "27", "30", "31", "32", "33", "34", "36", "39",
	"40", "41", "43", "44", "45", "46", "47", "48", "49",
	"51", "52", "53", "54", "55", "56", "57", "58",
	"60", "61", "62", "63", "64", "65", "66",
	"81", "82", "84", "86",
	"90", "91", "92", "93", "94", "95", "98",
	"211", "212", "213", "216", "218",
	"220", "221", "222", "223", "224", "225", "226", "227", "228", "229",
	"230", "231", "232", "233", "234", "235", "236", "237", "238", "239",
	"240", "241", "242", "243", "244", "245", "246", "247", "248", "249",
	"250", "251", "252", "253", "254", "255", "256", "257", "258", "260",
	"261", "262", "263", "264", "265", "266", "267", "268", "269",
	"290", "291", "297", "298", "299",
	"350", "351", "352", "353", "354", "355", "356", "357", "358", "359",
	"370", "371", "372", "373", "374", "375", "376", "377", "378", "379",
	"380", "381", "382", "385", "386", "387", "389",
	"420", "421", "423",
	"500", "501", "502", "503", "504", "505", "506", "507", "508", "509",
	"590", "591", "592", "593", "594", "595", "596", "597", "598", "599",
	"670", "672", "673", "674", "675", "676", "677", "678", "679", "680",
	"681", "682", "683", "685", "686", "687", "688", "689", "690", "691", "692",
	"850", "852", "853", "855", "856", "870", "871", "872", "873", "874", "878",
	"880", "881", "882", "883", "886", "888",
	"960", "961", "962", "963", "964", "965", "966", "967", "968", "970", "971", "972", "973", "974", "975", "976", "977", "979",
	"992", "993", "994", "995", "996", "998",
}

func main() {
	baseDir, err := os.Getwd()
	if err != nil {
		fmt.Println("获取工作目录失败：", err)
		return
	}

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		fmt.Println("读取目录失败：", err)
		return
	}

	// 匹配区号+手机号格式：数字，前面可以有 "+"
	re := regexp.MustCompile(`^\+?(\d{1,4})_(\d+)$`)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		originalName := entry.Name()

		// 先尝试匹配“区号_手机号”的格式
		matches := re.FindStringSubmatch(originalName)
		if matches != nil {
			// 匹配到了区号和手机号
			areaCode := matches[1] // 区号部分
			//phoneNumber := matches[2] // 手机号部分

			// 如果区号在列表中
			matched := false
			for _, code := range countryCodes {
				if areaCode == code {
					// 匹配到区号
					srcPath := filepath.Join(baseDir, originalName)
					targetDir := filepath.Join(baseDir, code)
					dstPath := filepath.Join(targetDir, originalName)

					// 创建目标文件夹
					err := os.MkdirAll(targetDir, os.ModePerm)
					if err != nil {
						fmt.Println("创建目录失败：", err)
						continue
					}

					// 移动文件夹
					if _, err := os.Stat(dstPath); err == nil {
						fmt.Printf("目标已存在，跳过：%s\n", dstPath)
						matched = true
						break
					}

					err = os.Rename(srcPath, dstPath)
					if err != nil {
						fmt.Printf("移动失败 %s -> %s：%v\n", srcPath, dstPath, err)
					} else {
						fmt.Printf("已移动 %s -> %s\\%s\n", originalName, code, originalName)
					}
					matched = true
					break
				}
			}

			if !matched {
				fmt.Println("未识别区号：", originalName)
			}
		} else {
			// 没有匹配到区号+手机号格式，检查文件夹名是否是纯数字
			// 匹配纯数字文件夹名
			areaCode := strings.TrimLeft(originalName, "+") // 去掉可能的"+"前缀

			// 检查文件夹名是否纯数字，并且是一个有效的区号
			if isNumeric(areaCode) {
				// 将纯数字拆分成区号和手机号
				matched := false
				for _, code := range countryCodes {
					// 如果区号部分在列表中
					if strings.HasPrefix(areaCode, code) {
						// 匹配到区号，后面的部分为手机号
						phoneNumber := areaCode[len(code):]
						if phoneNumber != "" {
							// 匹配到区号和手机号
							srcPath := filepath.Join(baseDir, originalName)
							targetDir := filepath.Join(baseDir, code)
							dstPath := filepath.Join(targetDir, originalName)

							// 创建目标文件夹
							err := os.MkdirAll(targetDir, os.ModePerm)
							if err != nil {
								fmt.Println("创建目录失败：", err)
								continue
							}

							// 移动文件夹
							if _, err := os.Stat(dstPath); err == nil {
								fmt.Printf("目标已存在，跳过：%s\n", dstPath)
								matched = true
								break
							}

							err = os.Rename(srcPath, dstPath)
							if err != nil {
								fmt.Printf("移动失败 %s -> %s：%v\n", srcPath, dstPath, err)
							} else {
								fmt.Printf("已移动 %s -> %s\\%s\n", originalName, code, originalName)
							}
							matched = true
							break
						}
					}
				}

				if !matched {
					fmt.Println("未识别区号：", originalName)
				}
			} else {
				// 不是纯数字，跳过
				fmt.Println("不匹配格式：", originalName)
			}
		}
	}

	fmt.Println("完成！按 Enter 键退出...")
	fmt.Scanln()
}

// 判断字符串是否是纯数字
func isNumeric(str string) bool {
	for _, c := range str {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
