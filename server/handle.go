package server

import (
	"chatroom/logic"
	"net/http"
	"os"
	"path/filepath"
)

// rootDir 根目录路径
var rootDir string

func RegisterHandle() {
	//负责推断出项目根目录
	inferRootDir()
	//接收广播消息
	go logic.broadCast.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", websocketHandleFunc)
}

// inferRootDir 推断出项目根目录
func inferRootDir() {
	dir, err := os.Getwd() //返回当前工作目录
	if err != nil {
		panic(err)
	}
	rootDir = inferWhereHaveTemplateDir(dir)

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
