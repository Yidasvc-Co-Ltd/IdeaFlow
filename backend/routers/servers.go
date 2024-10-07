package routers

import (
	"backend/gorm"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

//servers_create 传name,userID返回state
func Servers_create(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string `json:"state"`
	}

	//修改数据库
	gorm.Servers_create_sql(jsonData)
	response := Response{
		State: "OK",
	}
	c.JSON(200, response)
}

//servers_update 传所有,返回state
func Servers_update(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string `json:"state"`
	}

	gorm.Servers_update_sql(jsonData)

	response := Response{
		State: "OK",
	}
	c.JSON(200, response)
}

//servers_query_all 传 userID 返回 name数组和state
func Servers_query_all(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string                   `json:"state"`
		Data  []map[string]interface{} `json:"data"`
	}

	jsonMaps := gorm.Servers_query_all(jsonData)
	response := Response{
		State: "OK",
		Data:  jsonMaps,
	}
	c.JSON(200, response)
}

//servers_query传name, userID返回除了name, userID之外的所有
func Servers_query(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State            string `json:"state"`
		Ssh_user         string `json:"ssh_user"`
		Ip               string `json:"ip"`
		Port             string `json:"port"`
		Auth_method      string `json:"auth_method"`
		Password         string `json:"password"`
		Key              string `json:"key"`
		JumpServerName   string `json:"jumpServerName"`
		JumpServerUserID string `json:"jumpServerUserID"`
		Login_command    string `json:"login_command"`
	}

	request_jsonData := gorm.Servers_query_sql(jsonData)
	response := Response{
		State:            "OK",
		Ssh_user:         request_jsonData["ssh_user"].(string),
		Ip:               request_jsonData["ip"].(string),
		Port:             request_jsonData["port"].(string),
		Auth_method:      request_jsonData["auth_method"].(string),
		Password:         request_jsonData["password"].(string),
		Key:              request_jsonData["key"].(string),
		JumpServerName:   request_jsonData["jumpServerName"].(string),
		JumpServerUserID: request_jsonData["jumpServerUserID"].(string),
		Login_command:    request_jsonData["login_command"].(string),
	}
	c.JSON(200, response)
}

//servers_delete传name, userID返回state
func Servers_delete(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string `json:"state"`
	}

	//修改数据库
	gorm.Servers_delete_sql(jsonData)
	response := Response{
		State: "OK",
	}
	c.JSON(200, response)
}

func Servers_run(c *gin.Context, jsonData map[string]interface{}) {
	// 使用 scp 命令上传文件至服务器
	cmdUpload := exec.Command("scp", "dynTaskServer.py", "<服务器用户名>@<服务器IP地址>:/tmp/")
	err := cmdUpload.Run()
	if err != nil {
		log.Fatal(err)
	}

	// 在服务器上执行启动命令
	cmdStart := exec.Command("ssh", "<服务器用户名>@<服务器IP地址>", "nohup python3 /tmp/dynTaskServer.py &")
	err = cmdStart.Run()
	if err != nil {
		log.Fatal(err)
	}

	// 输出启动信息
	fmt.Println("服务已成功启动！")

	// 检查是否需要修改服务端口
	if len(os.Args) > 1 && os.Args[1] == "--port" && len(os.Args) > 2 {
		port := os.Args[2]
		cmdModifyPort := exec.Command("ssh", "<服务器用户名>@<服务器IP地址>", "python3 /tmp/dynTaskServer.py --port "+port)
		err = cmdModifyPort.Run()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("已将服务端口修改为 %s\n", port)
	}
}

func Operate_servers(operate string, c *gin.Context, jsonData map[string]interface{}) {

	switch operate {
	case "Servers_create":
		Servers_create(c, jsonData)
	case "Servers_query_all":
		Servers_query_all(c, jsonData)
	case "Servers_query":
		Servers_query(c, jsonData)
	case "Servers_update":
		Servers_update(c, jsonData)
	case "Servers_delete":
		Servers_delete(c, jsonData)
	case "Servers_run":
		Servers_run(c, jsonData)
	}
}
