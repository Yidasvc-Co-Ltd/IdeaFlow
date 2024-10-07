package routers

import (
	"backend/gorm"

	"github.com/gin-gonic/gin"
)

//创建段落document_paragraph_insert 传paragraphID (PK),知道在哪段之后插,传documentID和userID还需要知道在那个用户的那个文档里插,返回state(string默认是OK),修改数据库
func Document_paragraph_insert(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string `json:"state"`
	}

	gorm.Docuemnt_paragraph_insert_sql(jsonData)
	response := Response{
		State: "OK",
	}
	c.JSON(200, response)

}

//修改段落document_paragraph_modify 新增一行paragraphs,在documents中修改paragraphs的版本
func Document_paragraph_modify(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string `json:"state"`
	}

	gorm.Document_paragraph_modify_sql(jsonData)
	response := Response{
		State: "OK",
	}
	c.JSON(200, response)
}

//合并段落document_paragraph_merge 传paragraphID,documentID,userID 还要传number表示从这段以下要合并几段,返回state和paragraphsID 数据库新建paragraphs
func Document_paragraph_merge(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State       string `json:"state"`
		ParagraphID string `json:"paragraphID"`
	}

	gorm.Document_paragraph_merge_sql(jsonData)
	response := Response{
		State:       "OK",
		ParagraphID: "数据库中获取+1",
	}

	c.JSON(200, response)
}

//拆分段落document_paragraph_separate 传paragraphID,documentID,userID,传一个参数第一段的长度first_length,修改text  返回state和新的两个paragraphID
func Document_paragraph_separate(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State             string `json:"state"`
		FisrtParagraphID  string `json:"FirstParagraphID"`
		SecondParagraphID string `json:"secondParagraphID"`
	}

	gorm.Document_paragraph_separate_sql(jsonData)
	response := Response{
		State:             "OK",
		FisrtParagraphID:  "数据库中获取+1",
		SecondParagraphID: "数据库中获取+2",
	}

	c.JSON(200, response)
}

//查询段落 document_paragraph_query  发paragraphID,version,documentID,userID  返回state和text
func Document_paragraph_query(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string `json:"state"`
		Text  string `json:"text"`
	}

	request_data := gorm.Document_paragraph_query_sql(jsonData)
	response := Response{
		State: "OK",
		Text:  request_data["text"].(string),
	}

	c.JSON(200, response)
}

//document_paragraph_version_query paragraphID,documentID, userID 返回state和数组version,before,next
func Document_paragraph_version_query(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string                   `json:"state"`
		Data  []map[string]interface{} `json:"data"`
	}

	jsonMaps := gorm.Docuemnt_paragraph_version_query_sql(jsonData)
	response := Response{
		State: "OK",
		Data:  jsonMaps,
	}

	c.JSON(200, response)
}

//document_paragraph_version_jump 传paragraphID,version, documentID,userID 返回state
func Document_paragraph_version_jump(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string                 `json:"state"`
		Data  map[string]interface{} `json:"data"`
	}

	request_jsonData := gorm.Document_paragraph_version_jump_sql(jsonData)
	response := Response{
		State: "OK",
		Data:  request_jsonData,
	}

	c.JSON(200, response)
}

func Operate_paragraphs(operate string, c *gin.Context, jsonData map[string]interface{}) {

	switch operate {
	case "Document_paragraph_insert":
		Document_paragraph_insert(c, jsonData)
	case " Document_paragraph_modify":
		Document_paragraph_modify(c, jsonData)
	case "Document_paragraph_merge":
		Document_paragraph_merge(c, jsonData)
	case "Document_paragraph_separate":
		Document_paragraph_separate(c, jsonData)
	case " Document_paragraph_query":
		Document_paragraph_query(c, jsonData)
	case "Document_paragraph_version_query":
		Document_paragraph_version_query(c, jsonData)
	case "Document_paragraph_version_jump":
		Document_paragraph_version_jump(c, jsonData)
	}
}
