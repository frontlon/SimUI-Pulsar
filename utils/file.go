package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// 复制文件
func FileCopy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	dir := GetFilePath(dst)
	if !DirExists(dir) {
		CreateDir(dir)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func FolderCopy(from, to string) error {

	if from == "" || to == "" {
		return nil
	}

	var err error

	f, err := os.Stat(from)
	if err != nil {
		return err
	}

	fn := func(fromFile string) error {
		//复制文件的路径
		rel, err := filepath.Rel(from, fromFile)
		if err != nil {
			return err
		}
		toFile := filepath.Join(to, rel)

		//创建复制文件目录
		if err = os.MkdirAll(filepath.Dir(toFile), 0777); err != nil {
			return err
		}

		//读取源文件
		file, err := os.Open(fromFile)
		if err != nil {
			return err
		}

		defer file.Close()
		bufReader := bufio.NewReader(file)
		// 创建复制文件用于保存
		out, err := os.Create(toFile)
		if err != nil {
			return err
		}

		defer out.Close()
		// 然后将文件流和文件流对接起来
		_, err = io.Copy(out, bufReader)
		return err
	}

	//转绝对路径
	pwd, _ := os.Getwd()
	if !filepath.IsAbs(from) {
		from = filepath.Join(pwd, from)
	}
	if !filepath.IsAbs(to) {
		to = filepath.Join(pwd, to)
	}

	//复制
	if f.IsDir() {
		return filepath.WalkDir(from, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				return fn(path)
			} else {
				if err = os.MkdirAll(path, 0777); err != nil {
					return err
				}
			}
			return err
		})
	} else {
		return fn(from)
	}
}

// 移动文件
func FileMove(oldFile string, newFile string) error {
	if FileExists(oldFile) {
		err := os.Rename(oldFile, newFile)
		return err
	}
	return nil
}

// 移动文件夹
func FolderMove(oldFolder string, newFolder string) error {
	if DirExists(oldFolder) {
		err := os.Rename(oldFolder, newFolder)
		return err
	}
	return nil
}

/*
重命名文件
*/
func FileRename(oldpath string, filename string) (string, error) {
	if !FileExists(oldpath) {
		return "", nil
	}
	newpath := ReplaceFileNameByPath(oldpath, filename)
	return newpath, os.Rename(oldpath, newpath)
}

/*
重命名文件夹
*/
func FolderRename(oldpath string, filename string) error {
	if !DirExists(oldpath) {
		return nil
	}
	newpath := filepath.Dir(oldpath) + "/" + filename
	return os.Rename(oldpath, newpath)
}

// 删除文件
func FileDelete(src string) error {
	if FileExists(src) {
		err := os.Remove(src)
		return err
	}
	return nil
}

// 删除目录
func DeleteDir(src string) error {
	if DirExists(src) {
		err := os.RemoveAll(src)
		return err
	}
	return nil
}

/**
 * 返回文件大小 + 单位
 */
func GetFileSizeString(size int64) string {

	//小于2k不记录
	if size < 2048 {
		return ""
	}

	if size < 1024 {
		return fmt.Sprintf("%.2fB", float64(size)/float64(1))
	} else if size < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(size)/float64(1024))
	} else if size < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(size)/float64(1024*1024))
	} else if size < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(size)/float64(1024*1024*1024))
	} else {
		return fmt.Sprintf("%.2fEB", float64(size)/float64(1024*1024*1024*1024*1024))
	}
}

/**
 * 创建多层文件夹
 **/
func CreateDir(p string) error {

	if DirExists(p) {
		return nil
	}

	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		return err
	}
	if err := os.Chmod(p, os.ModePerm); err != nil {
		return err
	}
	return nil
}

/**
 * 创建一个文件
 **/
func CreateFile(p string, content string) error {
	if p == "" {
		return nil
	}
	f, err := os.Create(p)
	defer f.Close()

	if content != "" {
		if _, err = f.Write([]byte(content)); err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}
	return nil
}

/**
 * 写文件（覆盖写）
 **/
func OverlayWriteFile(p, content string) error {
	if p == "" {
		return nil
	}
	if err := os.WriteFile(p, []byte(content), 0664); err != nil {
		return err
	}
	return nil
}

/*
检查并读取文件内容
check:读取前是否检查文件存在
*/
func ReadFile(p string, check bool) (string, error) {
	if p == "" {
		return "", nil
	}

	if check && !FileExists(p) {
		return "", nil
	}
	bytes, err := os.ReadFile(p)
	if err != nil {
		return "", nil
	}

	return string(bytes), nil
}

// 读取文件的更新时间
func GetFileUpdateDate(p string) time.Time {
	f, err := os.Open(p)
	defer f.Close()
	if err != nil {
		return time.Time{}
	}
	stat, _ := f.Stat()
	return stat.ModTime()
}
