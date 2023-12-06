package server

import (
	"fmt"
	"html/template"
	"net/http"
)

func homeHandleFunc(writer http.ResponseWriter, request *http.Request) {
	tem, err := template.ParseFiles(rootDir + "/template/home.html")
	if err != nil {
		fmt.Fprintln(writer, "解析模板文件失败")
		return
	}
	err = tem.Execute(writer, nil)
	if err != nil {
		fmt.Fprintln(writer, "模板执行错误")
		return
	}

}
