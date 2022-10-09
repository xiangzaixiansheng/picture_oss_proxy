package util

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

//https://blog.csdn.net/qq_44675969/article/details/123929148?spm=1001.2101.3001.6650.2&utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7ECTRLIST%7Edefault-2-123929148-blog-119918880.pc_relevant_default&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7ECTRLIST%7Edefault-2-123929148-blog-119918880.pc_relevant_default&utm_relevant_index=5
//相对路径和绝对路径操作相关

// 最终方案-全兼容
func GetCurrentAbPath() string {
	dir := GetCurrentAbPathByExecutable()
	if strings.Contains(dir, GetTmpDir()) {
		return GetCurrentAbPathByCaller()
	}
	return dir
}

// 获取系统临时目录，兼容go run
func GetTmpDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
	}
	res, _ := filepath.EvalSymlinks(dir)
	return res
}

// 获取当前执行文件绝对路径
func GetCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func GetCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

// FileExist 判断文件是否存在
func FileExist(path string) bool {
	fi, err := os.Lstat(path)
	if err == nil {
		return !fi.IsDir()
	}
	return !os.IsNotExist(err)
}

//删除文件
func DeLFile(filePath string) error {
	return os.RemoveAll(filePath)
}
