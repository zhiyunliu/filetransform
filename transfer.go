package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func GetAllFiles(dirPth string, suffix string) (result []string, err error) {
	items, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return []string{}, err
	}

	PthSep := string(os.PathSeparator)
	for _, fi := range items {
		if fi.IsDir() { // 鐩�綍, 閫掑綊閬嶅巻
			tmplist, err := GetAllFiles(dirPth+PthSep+fi.Name(), suffix)
			if err != nil {
				continue
			}
			result = append(result, tmplist...)
		} else {
			// 杩囨护鎸囧畾鏍煎紡
			if strings.HasSuffix(fi.Name(), suffix) {
				result = append(result, dirPth+PthSep+fi.Name())
			}
		}
	}
	return result, nil
}

func transEncoding(file string) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("ReadFile:", err)
	}
	destBytes, err := GbkToUtf8(bytes)
	if err != nil {
		fmt.Println("GbkToUtf8:", err)
	}
	ioutil.WriteFile(file, destBytes, fs.ModePerm)
	time.Sleep(100 * time.Millisecond)
}

func GbkToUtf8(str []byte) (b []byte, err error) {
	r := transform.NewReader(bytes.NewReader(str), simplifiedchinese.GBK.NewDecoder())
	b, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}
	return
}

var suffix string

func init() {
	flag.StringVar(&suffix, "s", ".cpp", "后缀名")
}

func main() {

	flag.Parse()

	fmt.Println("后缀名:", suffix)

	xfiles, err := GetAllFiles(".", suffix)
	if err != nil {
		fmt.Println("error:", err)
	}
	total := len(xfiles)
	fmt.Println("文件数:", total)
	for i, file := range xfiles {
		fmt.Printf("Encoding:%d/%d : %s \n", i, total, file)
		transEncoding(file)
	}
}
