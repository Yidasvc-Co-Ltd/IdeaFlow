package routers

import (
	"backend/gorm"

	"github.com/gin-gonic/gin"
)

//dynTasks_create传name,dynFigueID,documentID,userID 返回state
func DynTasks_create(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string `json:"state"`
	}

	gorm.DynTasks_create_sql(jsonData)
	response := Response{
		State: "OK",
	}
	c.JSON(200, response)
}

//dynTasks_delete 传 name,dynFigueID,documentID,userID 返回state
func DynTasks_delete(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string `json:"state"`
	}

	gorm.DynTasks_delete_sql(jsonData)
	response := Response{
		State: "OK",
	}
	c.JSON(200, response)
}

//3dynTasks_update传name,dynFigueID,documentID,userID,codeShell, tag 返回state
func DynTasks_update(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string `json:"state"`
	}

	gorm.DynTasks_update_sql(jsonData)
	response := Response{
		State: "OK",
	}
	c.JSON(200, response)
}

//dynTasks_query_all 传dynFigueID,documentID,userID 返回state 和name, codeShell, timeStart,timeEnd,tag组成的数组
func DynTasks_query_all(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string                   `json:"state"`
		Data  []map[string]interface{} `json:"data"`
	}
	//修改数据库
	jsonMaps := gorm.DynTasks_query_all_sql(jsonData)
	response := Response{
		State: "OK",
		Data:  jsonMaps,
	}
	c.JSON(200, response)
}

//dynTasks_query_name传4.1,4.2,4.3,4.4 返回state和count
func DynTasks_query_name(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string                 `json:"state"`
		Count int                    `json:"count"`
		Data  map[string]interface{} `json:"data"`
	}

	request_jsonData := gorm.DynTasks_query_name_sql(jsonData)
	response := Response{
		State: "OK",
		Count: 666,
		Data:  request_jsonData,
	}
	c.JSON(200, response)
}

func Operate_dynTasks(operate string, c *gin.Context, jsonData map[string]interface{}) {

	switch operate {
	case "DynTasks_create":
		DynTasks_create(c, jsonData)
	case "DynTasks_query_all":
		DynTasks_query_all(c, jsonData)
	case "DynTasks_query_name":
		DynTasks_query_name(c, jsonData)
	case "DynTasks_update":
		DynTasks_update(c, jsonData)
	case "DynTasks_delete":
		DynTasks_delete(c, jsonData)
	}
}
