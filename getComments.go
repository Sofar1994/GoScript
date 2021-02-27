package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	REG = "[\u4e00-\u9fa5]+[a-zA-Z\\s\\S]+[\u4e00-\u9fa5]+|[\u4e00-\u9fa5]+|[\u4e00-\u9fa5]+[\\s\\S]*?[\u4e00-\u9fa5]+"
	REG2 = "\\*[\\s\\S]*?|//.*[\u4e00-\u9fa5]+|<!-[\\s\\S]*?-->"
	TARGET = ""
)

func walkFunc(path string, info os.FileInfo, err error) error {
	if info == nil {
		// 文件名称超过限定长度等其他问题也会导致info == nil
		// 如果此时return err 就会显示找不到路径，并停止查找。
		fmt.Println("can't find:(" + path + ")")
		return nil
	}
	if info.IsDir() {
		//fmt.Println("This is folder:(" + path + ")")
		return nil
	}  else {
		pathArr := Explode(path, "/")
		if In_array(pathArr[len(pathArr)-2], []string{"assets", "fonts", "imgs"}) {
			return nil
		}
		fmt.Println("This is file:(" + path + ")")
		commentOut, output, needHandle, err := readFile(path)
		fmt.Printf("%s\n",commentOut)
		if err != nil {
			panic(err)
		}
		if needHandle {
			err = writeToFile(path, output)
			path2 := "/Users/lx/Documents/test/comments.csv"
			err2 := writeToCSV(path2, commentOut)
			if err != nil && err2 != nil{
				panic(err)
			}
			fmt.Println(path)
			fmt.Println(path2)
		}
		return nil
	}
}

func showFileList(root string) {
	err := filepath.Walk(root, walkFunc)
	if err != nil {
		fmt.Printf("filepath.Walk() error: %v\n", err)
	}
	return
}

func readFile(filePath string) ([][]string, []byte, bool, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, nil, false, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	needHandle := false
	output := make([]byte, 0)
	commentOut := make([][]string, 0)
	row := 0
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				//fmt.Printf("%s\n",commentOut)
				return commentOut, output, needHandle, nil
			}
			return nil, nil, needHandle, err
		}
		if !isPrefix {
			row++
		}

		ok2, _ := regexp.MatchString(REG2, string(line))
		if ok, _ := regexp.MatchString(REG, string(line)); ok&&!ok2 {
			reg := regexp.MustCompile(REG)
			text := reg.Find(line)
			//fmt.Printf("%s\n", text)

			//fmt.Printf("%s\n", text) //文本
			//fmt.Printf("%d\n", row) // row
			//fmt.Println(strings.Index(string(line), text[0])) // col
			col := strings.Index(string(line), string(text))
			newByte := reg.ReplaceAll(line, []byte(TARGET))
			output = append(output, newByte...)
			output = append(output, []byte("\n")...)

			commentOut = append(commentOut, []string{filePath, string(text), " row:"+strconv.Itoa(row)+" col:"+strconv.Itoa(col)})
			//commentOut = append(commentOut, []byte(filePath+": ")...)
			//commentOut = append(commentOut, text...)
			//commentOut = append(commentOut, []byte(" row:"+strconv.Itoa(row)+" col:")...)
			//commentOut = append(commentOut, []byte(strconv.Itoa(col))...)
			//commentOut = append(commentOut, []byte("\n")...)
			if !needHandle {
				needHandle = true
			}
		} else {
			output = append(output, line...)
			output = append(output, []byte("\n")...)
		}
	}
	return commentOut, output, needHandle, nil
}

func writeToFile(filePath string, outPut []byte) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_, _ = os.Create(filePath)
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	_, err = writer.Write(outPut)
	if err != nil {
		return err
	}
	_ = writer.Flush()
	return nil
}

func writeToCSV(filePath string, outPut [][]string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_, _ = os.Create(filePath)
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	_, _ = file.WriteString("\xEF\xBB\xBF")
	writer := csv.NewWriter(file)
	err = writer.WriteAll(outPut)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func Explode(delimiter, text string) []string {
	if len(delimiter) > len(text) {
		return strings.Split(delimiter, text)
	} else {
		return strings.Split(text, delimiter)
	}
}

func In_array(needle interface{}, hystack interface{}) bool {
	switch key := needle.(type) {
	case string:
		for _, item := range hystack.([]string) {
			if key == item {
				return true
			}
		}
	case int:
		for _, item := range hystack.([]int) {
			if key == item {
				return true
			}
		}
	case int64:
		for _, item := range hystack.([]int64) {
			if key == item {
				return true
			}
		}
	default:
		return false
	}
	return false
}

func main() {
	var path string
	//fmt.Println("输入要遍历的路径：")
	//_, _ = fmt.Scanf("%s", &path)
	//fmt.Println("=============")
	path = "/Users/lx/Documents/vue/fe/src"
	showFileList(path)
}

