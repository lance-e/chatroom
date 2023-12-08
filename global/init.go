package global

import (
	"os"
	"path/filepath"
	"sync"
)

func init() {
	Init()
}

var once = sync.Once{}

func Init() {
	once.Do(func() {
		inferRootDir()
		initConfig()
	})
}

// rootDir 根目录路径
var RootDir string

// inferRootDir 推断出项目根目录
func inferRootDir() {
	dir, err := os.Getwd() //返回当前工作目录
	if err != nil {
		panic(err)
	}
	RootDir = inferWhereHaveTemplateDir(dir)

}

// inferWhereHaveTemplateDir 确保根目录下存在 template 目录,递归调用
func inferWhereHaveTemplateDir(dir string) string {
	if exist(dir) {
		return dir
	}
	return inferWhereHaveTemplateDir(filepath.Dir(dir)) //传入上一级目录进行判断是否是根目录
}

// exist 表示文件是否存在
func exist(filename string) bool {
	_, err := os.Stat(filename + "/template")
	return err == nil || os.IsExist(err)
}
