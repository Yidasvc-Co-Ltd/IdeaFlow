package routers

import (
	"backend/gorm"

	"github.com/gin-gonic/gin"
)

//dynFigures_create传documentID,userID,name 返回state和dynFigureID
func DynFigures_create(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State       string `json:"state"`
		DynFigureID int    `json:"dynFigureID"`
	}

	//修改数据库
	dynFigureID := gorm.DynFigures_create_sql(jsonData)

	response := Response{
		State:       "OK",
		DynFigureID: dynFigureID,
	}
	c.JSON(200, response)
}

//dynFigures_query_all传documentID,userID,返回dynFigueID,name构成的数组
func DynFigures_query_all(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string                   `json:"state"`
		Data  []map[string]interface{} `json:"data"`
	}

	//修改数据库
	jsonMap := gorm.DynFigures_query_all_sql(jsonData)
	response := Response{
		State: "OK",
		Data:  jsonMap,
	}
	c.JSON(200, response)
}

//	dynFigures_query传dynFigueID,documentID,userID返回name, currentTag, codeGenTask,codeFig,tagQueue
func DynFigures_query(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State       string `json:"state"`
		Name        string `"json:"name"`
		CurrentTag  string `json:"currentTag"`
		CodeGenTask string `json:"codeGenTask"`
		CodeFig     string `json:"codeFig"`
		TagQueue    string `json:"tagQueue"`
	}

	//修改数据库
	querydata := gorm.DynFigures_query_sql(jsonData)
	response := Response{
		State:       "OK",
		Name:        querydata["Name"].(string),
		CurrentTag:  querydata["currentTag"].(string),
		CodeGenTask: querydata["codeGenTask"].(string),
		CodeFig:     querydata["codeFig"].(string),
		TagQueue:    querydata["tagQueue"].(string),
	}
	c.JSON(200, response)
}

//dynFigures_update传dynFigueID,documentID,userID,name ,currentTag, codeGenTask ,codeFig,tagQueue 返回state
func DynFigures_update(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string `json:"state"`
	}

	//修改数据库
	gorm.DynFigures_update_sql(jsonData)
	response := Response{
		State: "OK",
	}
	c.JSON(200, response)
}

//dynFigures_delete传3.1,3.2,3.3返回state
func DynFigures_delete(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string `json:"state"`
	}

	//修改数据库
	gorm.DynFigures_delete_sql(jsonData)

	response := Response{
		State: "OK",
	}
	c.JSON(200, response)
}

func Operate_dynFigures(operate string, c *gin.Context, jsonData map[string]interface{}) {

	switch operate {
	case "DynFigures_create":
		DynFigures_create(c, jsonData)
	case "DynFigures_query_all":
		DynFigures_query_all(c, jsonData)
	case "DynFigures_query":
		DynFigures_query(c, jsonData)
	case "DynFigures_update":
		DynFigures_update(c, jsonData)
	case "DynFigures_delete":
		DynFigures_delete(c, jsonData)
	}
}
