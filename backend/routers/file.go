package routers

import (
	"backend/gorm"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func File_id_all(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string                   `json:"state"`
		Data  []map[string]interface{} `json:"data"`
	}
	data := gorm.File_id_all(jsonData)
	response := Response{
		State: "OK", // 默认状态为 "OK"
		Data:  data,
	}
	c.JSON(200, response)
}

func Picture_upload(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string `json:"state"`
	}
	gorm.Picture_upload_sql(jsonData)
	response := Response{
		State: "OK",
	}
	c.JSON(200, response)
}

func Picture_upload_query_sql(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string      `json:"state"`
		Data  models.File `json:"data"`
	}
	data := gorm.Picture_upload_query_sql(jsonData)
	response := Response{
		State: "OK",
		Data:  data,
	}
	c.JSON(200, response)
}

func Operate_file(operate string, c *gin.Context, jsonData map[string]interface{}) {

	switch operate {
	case "File_id_all":
		File_id_all(c, jsonData)
	case "Picture_upload":
		Picture_upload(c, jsonData)
	case "Picture_upload_query_sql":
		Picture_upload_query_sql(c, jsonData)
	}

}
