package server

import (
	"chatroom/global"
	"chatroom/logic"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func homeHandleFunc(writer http.ResponseWriter, request *http.Request) {
	tem, err := template.ParseFiles(global.RootDir + "/template/home.html")
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

func UserListHandleFunc(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	writer.WriteHeader(http.StatusOK)

	userList := logic.BroadCaster.UserList()
	users, err := json.Marshal(userList)
	if err != nil {
		fmt.Fprint(writer, `[]`)
	} else {
		fmt.Fprint(writer, string(users))
	}
}
