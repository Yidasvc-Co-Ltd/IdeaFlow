package routers

import (
	"backend/gorm"

	"github.com/gin-gonic/gin"
)

//paragraph_templates_create 传 name,userID返回state
func Paragraph_templates_create(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string `json:"state"`
	}

	gorm.Paragraph_templates_create_sql(jsonData)
	response := Response{
		State: "OK",
	}
	c.JSON(200, response)
}

//paragraph_templates_update传所有,返回state
func Paragraph_templates_update(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string `json:"state"`
	}

	gorm.Paragraph_templates_update_sql(jsonData)
	response := Response{
		State: "OK",
	}
	c.JSON(200, response)
}

// paragraph_templates_query_all 传userID返回name组成的数组
func Paragraph_templates_query_all(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string                   `json:"state"`
		Data  []map[string]interface{} `json:"data"`
	}

	jsonMaps := gorm.Paragraph_templates_query_all_sql(jsonData)
	response := Response{
		State: "OK",
		Data:  jsonMaps,
	}
	c.JSON(200, response)
}

// paragraph_templates_query传name,userID返回 tag,text
func Paragraph_templates_query(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string `json:"state"`
		Tag   string `json:"tag"`
		Text  string `json:"text"`
	}

	request_jsonData := gorm.Paragraph_templates_query_sql(jsonData)
	response := Response{
		State: "OK",
		Tag:   request_jsonData["tag"].(string),
		Text:  request_jsonData["Text"].(string),
	}
	c.JSON(200, response)
}

// paragraph_templates_delete传name,userID返回state
func Paragraph_templates_delete(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string `json:"state"`
	}

	gorm.Paragraph_templates_delete_sql(jsonData)
	response := Response{
		State: "OK",
	}
	c.JSON(200, response)
}

func Operate_paragraph_templates(operate string, c *gin.Context, jsonData map[string]interface{}) {

	switch operate {
	case "Paragraph_templates_create":
		Paragraph_templates_create(c, jsonData)
	case "Paragraph_templates_update":
		Paragraph_templates_update(c, jsonData)
	case "Paragraph_templates_query_all":
		Paragraph_templates_query_all(c, jsonData)
	case "Paragraph_templates_query":
		Paragraph_templates_query(c, jsonData)
	case "Paragraph_templates_delete":
		Paragraph_templates_delete(c, jsonData)
	}
}
