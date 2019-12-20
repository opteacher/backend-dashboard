package utils

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func PickPathsFromSwaggerJSON(fname string) ([]byte, error) {
	var buffer []byte
	length := 0
	if file, err := os.Open(fname); err != nil {
		return nil, err
	} else {
		defer file.Close()

		chunks := make([]byte, 1024, 1024)
		for {
			if n, err := file.Read(chunks); n == 0 {
				break
			} else if err != nil {
				return nil, err
			} else {
				length += n
				buffer = append(buffer, chunks...)
			}
		}
	}

	data := make(map[string]interface{})
	if err := json.Unmarshal(buffer[:length], &data); err != nil {
		return nil, err
	} else if paths, exs := data["paths"]; !exs {
		return nil, fmt.Errorf("swagger文件不存在paths字段，是否格式有变：%s", string(buffer))
	} else {
		return json.Marshal(paths)
	}
}

// 扫描文件夹下文件
// folderPath 指定文件夹
// suffix 指定后缀，如不需要指定为"*"
func ScanAllFilesByFolder(folderPath string, suffix string) ([]string, error) {
	if suffix != "*" {
		suffix = strings.ToLower(suffix)
	}
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}
	var retFiles []string
	for _, file := range files {
		fname := file.Name()
		if file.IsDir() {
			flAry, err := ScanAllFilesByFolder(filepath.Join(folderPath, fname), suffix)
			if err != nil {
				return nil, err
			}
			retFiles = append(retFiles, flAry...)
		} else {
			poiIdx := strings.LastIndex(fname, ".")
			if (poiIdx != -1 && strings.ToLower(fname[poiIdx + 1:]) == suffix) || suffix == "*" {
				retFiles = append(retFiles, filepath.Join(folderPath, fname))
			}
		}
	}
	return retFiles, nil
}

// 压缩文件
// files 文件数组，可以是不同dir下的文件或者文件夹
// dest 压缩文件存放地址
func Compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := compress(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func Download(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return ioutil.WriteFile(dest, bytes, os.ModePerm)
}

func CopyFolder(src string, dest string) {
	srcOriginal := src
	err := filepath.Walk(src, func(src string, file os.FileInfo, err error) error {
		if file == nil {
			return err
		}
		if !file.IsDir() {
			//fmt.Println(src)
			//fmt.Println(src_original)
			//fmt.Println(dest)

			destNew := strings.Replace(src, srcOriginal, dest, -1)
			//fmt.Println(dest_new)
			//fmt.Println("CopyFile:" + src + " to " + dest_new)
			if _, err := CopyFile(src, destNew); err != nil {
				return err
			}
		}
		//println(path)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

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

//copy file
func CopyFile(src, dst string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer srcFile.Close()
	// fmt.Println("dst:" + dst)
	separator := string(filepath.Separator)
	dstSlices := strings.Split(dst, separator)
	dstSlicesLen := len(dstSlices)
	destDir := ""
	for i := 0; i < dstSlicesLen-1; i++ {
		destDir = destDir + dstSlices[i] + separator
	}
	//dest_dir := getParentDirectory(dst)
	// fmt.Println("dest_dir:" + dest_dir)
	b, err := PathExists(destDir)
	if b == false {
		err := os.MkdirAll(destDir, os.ModePerm) //在当前目录下生成md目录
		if err != nil {
			fmt.Println(err)
		}
	}
	dstFile, err := os.Create(dst)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

// 文件插入（逐行操作，只使用文本文件）
type InsertProcFunc func(string, *bool) (string, bool, error)
func InsertTxt(fpath string, proc InsertProcFunc) error {
	// 读取import部分
	file, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	doProc := true
	code := ""
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		str := string(line)
		if scd, isBreak, err := proc(str, &doProc); err != nil {
			return err
		} else if isBreak {
			break
		} else if doProc {
			code += scd + "\n"
		} else {
			code += str + "\n"
		}
	}
	// 写入
	file.Close()
	file, err = os.OpenFile(fpath, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(code); err != nil {
		return err
	}
	return nil
}

type ReplaceProcFunc func(string) (string, error)
func ReplaceContentInFile(flPath string, rpFuns map[string]ReplaceProcFunc) error {
	file, err := os.Open(flPath)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	code := ""
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		strLin := string(line)
		matched := false
		for repWds, procFun := range rpFuns {
			if strings.Contains(strLin, repWds) {
				matched = true
				genTxt, err := procFun(strLin)
				if err != nil {
					return err
				}
				code += genTxt
			}
		}
		if !matched {
			code += strLin + "\n"
		}
	}
	file, err = os.OpenFile(flPath, os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(code); err != nil {
		return fmt.Errorf("重新写入文件失败：%v", err)
	}
	return nil
}