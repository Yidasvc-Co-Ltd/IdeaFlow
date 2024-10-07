package app

import (
	"backend/routers"

	"github.com/gin-gonic/gin"
)

func Receive_data(c *gin.Context) map[string]interface{} {
	var jsonData map[string]interface{}
	// 使用 ShouldBindJSON 方法将 JSON 数据绑定到 map 中
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	return jsonData
}
func App_go(c *gin.Context) {
	jsonData := Receive_data(c)
	operate_type := jsonData["operate_type"].(string)
	operate := jsonData["operate"].(string)
	switch operate_type {
	case "Operate_documents":
		routers.Operate_documents(operate, c, jsonData)
	case "Operate_paragraphs":
		routers.Operate_paragraphs(operate, c, jsonData)
	case "Operate_dynFigures":
		routers.Operate_dynFigures(operate, c, jsonData)
	case "Operate_dynTasks":
		routers.Operate_dynTasks(operate, c, jsonData)
	case "Operate_servers":
		routers.Operate_servers(operate, c, jsonData)
	case "Operate_paragraph_templates":
		routers.Operate_paragraph_templates(operate, c, jsonData)
	case "Operate_pdf_create":
		routers.Operate_pdf_create(operate, c, jsonData)
	case "Operate_file":
		routers.Operate_file(operate, c, jsonData)
	}
}
