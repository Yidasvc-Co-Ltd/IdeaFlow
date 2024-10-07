package gorm

import (
	"backend/models"
	"encoding/json"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect_To_Database_Fourth() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("/root/server.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func DynTasks_create_sql(jsonData map[string]interface{}) {
	db := Connect_To_Database_Fourth()
	// 迁移模型
	db.AutoMigrate(&models.DynTasks{})
	dyntasks := models.DynTasks{
		Name:        jsonData["name"].(string),
		DynFigureID: jsonData["dynFigureID"].(int),
		DocumentID:  jsonData["documentID"].(int),
		UserID:      jsonData["userID"].(string),
		CodeShell:   "",
		TimeStart:   "",
		TimeEnd:     "",
		Tag:         "",
	}
	// 在数据库中创建记录
	result := db.Model(&models.DynTasks{}).Create(&dyntasks)
	if result.Error != nil {
		panic(result.Error)
	}
}

func DynTasks_delete_sql(jsonData map[string]interface{}) {
	db := Connect_To_Database_Fourth()
	// 迁移模型
	db.AutoMigrate(&models.DynTasks{})
	deleteCondition := models.DynTasks{
		Name:        jsonData["name"].(string),
		DynFigureID: jsonData["dynFigureID"].(int),
		DocumentID:  jsonData["documentID"].(int),
		UserID:      jsonData["userID"].(string),
		CodeShell:   "",
		TimeStart:   "",
		TimeEnd:     "",
		Tag:         "",
	}
	// 在数据库中创建记录
	result := db.Where(&deleteCondition).Delete(&models.DynTasks{})
	if result.Error != nil {
		panic(result.Error)
	}
}

func DynTasks_update_sql(jsonData map[string]interface{}) {
	db := Connect_To_Database_Fourth()

	// 迁移模型
	db.AutoMigrate(&models.DynTasks{})

	type UpdateCondition struct {
		Name        string `"json:"name"`
		DynFigureID int    `json:"dynFigureID"`
		DocumentID  int    `json"documentID"`
		UserID      string `json"userID"`
	}
	updateCondition := UpdateCondition{
		DocumentID:  jsonData["documentID"].(int),
		UserID:      jsonData["userID"].(string),
		DynFigureID: jsonData["dynFigureID"].(int),
		Name:        jsonData["name"].(string),
	}

	// 更新记录
	result_codeShell := db.Model(&models.DynTasks{}).Where(&updateCondition).Update("code_shell", jsonData["codeShell"].(string))
	result_tag := db.Model(&models.DynTasks{}).Where(&updateCondition).Update("tag", jsonData["tag"].(string))
	result_name := db.Model(&models.DynTasks{}).Where(&updateCondition).Update("name", jsonData["new_name"].(string))

	if result_name.Error != nil {
		panic(result_name.Error)
	}
	if result_codeShell.Error != nil {
		panic(result_codeShell.Error)
	}
	if result_tag.Error != nil {
		panic(result_tag.Error)
	}
}

func DynTasks_query_all_sql(jsonData map[string]interface{}) []map[string]interface{} {
	db := Connect_To_Database_Fourth()
	var dynTasks []models.DynTasks
	db.Where("document_id = ? AND user_id = ? AND dyn_figure_id= ? ", jsonData["documentID"].(int), jsonData["userID"].(string), jsonData["dynFigureID"]).Find(&dynTasks)
	var jsonArray []string
	for _, doc := range dynTasks {
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
			"name":      jsonMap["Name"],
			"codeShell": jsonMap["code_shell"],
			"timeStart": jsonMap["time_start"],
			"timeEnd":   jsonMap["time_end"],
			"tag":       jsonMap["tag"],
		}
		jsonMaps = append(jsonMaps, jsonFilter)
	}
	return jsonMaps
}

func DynTasks_query_name_sql(jsonData map[string]interface{}) map[string]interface{} {
	db := Connect_To_Database_Fourth()
	// 迁移模型
	db.AutoMigrate(&models.DynTasks{})

	type QueryCondition struct {
		DynFigureID int    `json:"dynFigureID"`
		DocumentID  int    `json:"documentID"`
		UserID      string `json:"userID"`
		Name        string `json:"name"`
	}
	queryCondition := QueryCondition{
		DynFigureID: jsonData["dynFigureID"].(int),
		DocumentID:  jsonData["documentID"].(int),
		UserID:      jsonData["userID"].(string),
		Name:        jsonData["name"].(string),
	}

	// 查询记录
	var result models.DynTasks
	if err := db.Where(&queryCondition).First(&result).Error; err != nil {
		panic(err)
	}

	// 将结果转换为 JSON 格式
	request, err := json.Marshal(result)
	if err != nil {
		// 处理 JSON 转换错误
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
