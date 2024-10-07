package gorm

import (
	"backend/models"
	"encoding/json"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect_To_Database_Third() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("/root/server.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

//根据同一个documentID实现dyFigures的自增
func DynFigures_self_increasing(db *gorm.DB, documentID int, userID string) int {
	db.AutoMigrate(&models.DynFigures{})
	// 执行查询，按照创建dyFigures降序排序，获取最新一条符合条件的记录的文档ID
	var dynfigure models.DynFigures
	if err := db.Where("document_id = ? AND user_id = ?", documentID, userID).Order("dyn_figure_id DESC").First(&dynfigure).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			//首次创建
			return 1
		} else {
			panic(err)
		}
	}
	return dynfigure.DynFigureID + 1
}

func DynFigures_create_sql(jsonData map[string]interface{}) int {
	db := Connect_To_Database_Third()
	// 迁移模型
	db.AutoMigrate(&models.DynFigures{})
	dynfigure := DynFigures_self_increasing(db, int(jsonData["documentID"].(float64)), jsonData["userID"].(string))
	dynFigures := models.DynFigures{
		DynFigureID: dynfigure,
		DocumentID:  int(jsonData["documentID"].(float64)),
		UserID:      jsonData["userID"].(string),
		Name:        jsonData["name"].(string),
		CurrentTag:  "nil",
		CodeGenTask: "请输入python代码",
		CodeFig:     "请输入python代码",
		TagQueue:    "111",
	}
	// 在数据库中创建记录
	result := db.Model(&models.DynFigures{}).Create(&dynFigures)
	if result.Error != nil {
		panic(result.Error)
	}
	return dynfigure
}

func DynFigures_query_all_sql(jsonData map[string]interface{}) []map[string]interface{} {
	db := Connect_To_Database_Third()
	var dynFigures []models.DynFigures
	documentID := int(jsonData["documentID"].(float64))
	db.Where("document_id = ? AND user_id = ?", documentID, jsonData["userID"].(string)).Find(&dynFigures)
	var jsonArray []string
	for _, doc := range dynFigures {
		jsonDateConversion, err := json.Marshal(doc)
		if err != nil {
			panic(err)
		}
		jsonArray = append(jsonArray, string(jsonDateConversion))
	}
	// 创建一个空的 []map[string]interface{} 类型的切片，用于存储转换后的数据
	var jsonMaps []map[string]interface{}
	// 遍历每个 JSON 字符串
	for _, jsonString := range jsonArray {
		// 创建一个 map 用于存储解析后的 JSON 数据
		var jsonMap map[string]interface{}

		// 解析 JSON 字符串到 map[string]interface{} 类型
		errs := json.Unmarshal([]byte(jsonString), &jsonMap)
		if errs != nil {
			panic(errs)
		}
		// 将解析后的 map 添加到 jsonMaps 切片中
		jsonFilter := map[string]interface{}{
			"dynFigureID": int(jsonMap["dynFigureID"].(float64)),
			"name":        jsonMap["Name"],
			"codeFig":     jsonMap["codeFig"].(string),
			"codeGenTask": jsonMap["codeGenTask"].(string),
			"currentTag":  jsonMap["currentTag"].(string),
			"tagQueue":    jsonMap["tagQueue"].(string),
		}
		jsonMaps = append(jsonMaps, jsonFilter)
	}
	return jsonMaps
}

func DynFigures_query_sql(jsonData map[string]interface{}) map[string]interface{} {
	db := Connect_To_Database_Third()
	// 迁移模型
	db.AutoMigrate(&models.DynFigures{})

	type QueryCondition struct {
		DynFigureID int    `json:"dynFigureID"`
		DocumentID  int    `json:"documentID"`
		UserID      string `json:"userID"`
	}
	queryCondition := QueryCondition{
		DynFigureID: int(jsonData["dynFigureID"].(float64)),
		DocumentID:  int(jsonData["documentID"].(float64)),
		UserID:      jsonData["userID"].(string),
	}

	// 查询记录
	var result models.DynFigures
	if err := db.Where(&queryCondition).First(&result).Error; err != nil {
		panic(err)
	}

	// 将结果转换为 JSON 格式
	request, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	// 将 []byte 转换为 map[string]interface{}
	var request_jsonData map[string]interface{}
	err = json.Unmarshal(request, &request_jsonData)
	if err != nil {
		panic(err)
	}

	return request_jsonData
}

func DynFigures_update_sql(jsonData map[string]interface{}) {
	db := Connect_To_Database_Third()

	// 迁移模型
	db.AutoMigrate(&models.DynFigures{})

	type UpdateCondition struct {
		DocumentID  int    `json:"documentID"`
		UserID      string `json:"userID"`
		DynFigureID int    `json:"dynFigureID"`
	}
	updateCondition := UpdateCondition{
		DocumentID:  int(jsonData["documentID"].(float64)),
		UserID:      jsonData["userID"].(string),
		DynFigureID: jsonData["dynFigureID"].(int),
	}

	// 更新记录
	result_name := db.Model(&models.DynFigures{}).Where(&updateCondition).Update("name", jsonData["name"].(string))
	result_currentTag := db.Model(&models.DynFigures{}).Where(&updateCondition).Update("current_tag", jsonData["currentTag"].(string))
	result_codeGenTask := db.Model(&models.DynFigures{}).Where(&updateCondition).Update("code_gen_task", jsonData["codeGenTask"].(string))
	result_codeFig := db.Model(&models.DynFigures{}).Where(&updateCondition).Update("code_fig", jsonData["codeFig"].(string))
	result_tagQueue := db.Model(&models.DynFigures{}).Where(&updateCondition).Update("tag_queue", jsonData["tagQueue"].(string))

	if result_name.Error != nil {
		panic(result_name.Error)
	}
	if result_currentTag.Error != nil {
		panic(result_currentTag.Error)
	}
	if result_codeGenTask.Error != nil {
		panic(result_codeGenTask.Error)
	}
	if result_codeFig.Error != nil {
		panic(result_codeFig.Error)
	}
	if result_tagQueue.Error != nil {
		panic(result_tagQueue.Error)
	}
}

func DynFigures_delete_sql(jsonData map[string]interface{}) {
	db := Connect_To_Database_Third()
	// 迁移模型
	db.AutoMigrate(&models.DynFigures{})

	deleteCondition := models.DynFigures{DynFigureID: int(jsonData["dynFigureID"].(float64)), DocumentID: int(jsonData["documentID"].(float64)), UserID: jsonData["userID"].(string)}

	// 删除记录
	result := db.Where(&deleteCondition).Delete(&models.DynFigures{})
	if result.Error != nil {
		panic(result.Error)
	}
}
