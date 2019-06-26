package utils

import (
	"os"
	"fmt"
	"encoding/json"
	"archive/zip"
	"io"
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