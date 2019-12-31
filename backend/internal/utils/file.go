package utils

import (
	"archive/zip"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
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

type StorageConfig struct {
	Url       string
	Bucket    string
	AccessKey string
	SecretKey string
}

type ProgressRecord struct {
	Progresses []storage.BlkputRet `json:"progresses"`
}

func Upload(absFile string, sc StorageConfig) (string, error) {
	putPolicy := storage.PutPolicy{
		Scope: sc.Bucket,
	}
	mac := qbox.NewMac(sc.AccessKey, sc.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuadong
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false

	finfo, err := os.Stat(absFile)
	if err != nil {
		return "", err
	}

	fsize := finfo.Size()
	lstIdx := strings.LastIndex(absFile, string(filepath.Separator))
	dname := ""
	fname := absFile
	if lstIdx != -1 {
		dname = absFile[:lstIdx]
		fname = absFile[lstIdx + 1:]
	}
	flmd := finfo.ModTime().UnixNano()
	recordKey := Md5Hex(fmt.Sprintf("%s:%s:%s:%d", sc.Bucket, fname, absFile, flmd)) + ".progress"
	recordPath := filepath.Join(dname, recordKey)
	if recordPath[0] != filepath.Separator {
		recordPath = string(filepath.Separator) + recordPath
	}

	pgsRcd := ProgressRecord{}

	if rcdFile, err := os.Open(recordPath); err != nil {

	} else if pgsByte, err := ioutil.ReadAll(rcdFile); err != nil {
		return "", err
	} else if err := json.Unmarshal(pgsByte, &pgsRcd); err != nil {
		return "", err
	} else {
		for _, item := range pgsRcd.Progresses {
			if storage.IsContextExpired(item) {
				fmt.Println(item.ExpiredAt)
				pgsRcd.Progresses = make([]storage.BlkputRet, storage.BlockCount(fsize))
				break
			}
		}
		rcdFile.Close()
	}

	if len(pgsRcd.Progresses) == 0 {
		pgsRcd.Progresses = make([]storage.BlkputRet, storage.BlockCount(fsize))
	}

	resumeUploader := storage.NewResumeUploader(&cfg)
	ret := storage.PutRet{}
	pgsLock := sync.RWMutex{}
	putExtra := storage.RputExtra{
		Progresses: pgsRcd.Progresses,
		Notify: func(blkIdx int, blkSize int, ret *storage.BlkputRet) {
			pgsLock.Lock()
			pgsLock.Unlock()

			pgsRcd.Progresses[blkIdx] = *ret
			pgsBytes, _ := json.Marshal(pgsRcd)
			fmt.Println("write progress file", blkIdx, recordPath)
			if err := ioutil.WriteFile(recordPath, pgsBytes, 0644); err != nil {
				panic(err)
			}
		},
	}
	if err := resumeUploader.PutFile(context.Background(), &ret, upToken, fname, absFile, &putExtra); err != nil {
		return "", err
	}
	if err := os.Remove(recordPath); err != nil {
		return "", err
	}
	url := sc.Url + "/" + ret.Key
	fmt.Printf("%s上传成功，哈希：%s，通过%s可下载\n", ret.Key, ret.Hash, url)
	return url, nil
}

func Download(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	lstIdx := strings.LastIndex(dest, string(filepath.Separator))
	if lstIdx != -1 {
		dirPath := dest[:lstIdx]
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
				return err
			}
		}
	}
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
// @Param{1}: 读取到的一行文本
// @Param{2}: 中断处理标识，在回调中设为false会导致之后的文本均不加处理
// @Return{1}: 处理之后的行文本
// @Return{2}: 是否直接停止处理
// @Return{3}: 处理过程发生的错误
type InsertProcFunc func(string, *bool) (string, bool, error)
// @Param{1}: 文件路径
// @Param{2}: 处理回调
// @Param{3}: 默认是否开启中断处理
func InsertTxt(fpath string, proc InsertProcFunc, first bool) error {
	// 读取import部分
	file, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	doProc := first
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
			if len(scd) != 0 {
				code += scd + "\n"
			}
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

// 删除tag包围的文字块（连同tag一起删除）
// @Param{1}：文件位置
// @Param{2}：tag字段，会在其后追加（_BEG）和（_END）形成开始块和结束块
// @Param{3}：是否删除文字块，如果选择不删除，该函数会删除tag字段，但保留块内文字
func DelByTagInFile(flPath string, tag string, delete bool) error {
	file, err := os.Open(flPath)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	code := ""
	beg := fmt.Sprintf("[%s_BEG]", tag)
	end := fmt.Sprintf("[%s_END]", tag)
	skip := false
	for {
		byteLine, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		strLine := string(byteLine)
		if strings.Index(strLine, end) != -1 {
			skip = false
			continue
		}
		if skip {
			continue
		}
		if strings.Index(strLine, beg) != -1 {
			skip = delete
			continue
		}
		code += strLine + "\n"
	}
	if err := os.Truncate(flPath, 0); err != nil {
		return err
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