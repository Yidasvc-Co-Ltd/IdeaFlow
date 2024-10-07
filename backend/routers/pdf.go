package routers

import (
	"backend/gorm"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func compileLatex(latexSource string, outputDirectory string, pdfName string) {
	// Write LaTeX source code to a temporary .tex file
	tempTexFile := "temp.tex"
	err := ioutil.WriteFile(tempTexFile, []byte(latexSource), 0644)
	if err != nil {
		fmt.Printf("Failed to write LaTeX source to temporary file: %v\n", err)
		return
	}
	defer os.Remove(tempTexFile) // Clean up: remove temporary .tex file

	// Execute pdflatex to compile the LaTeX document
	cmd := exec.Command("pdflatex", tempTexFile)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to compile LaTeX document: %v\n", err)
		return
	}

	// Get the name of the generated PDF file
	pdfFile := filepath.Join(".", tempTexFile[:len(tempTexFile)-len(filepath.Ext(tempTexFile))]+".pdf")

	// Rename the PDF file according to specified rules
	newPDFName := pdfName + ".pdf" // Modify as needed
	err = os.Rename(pdfFile, filepath.Join(outputDirectory, newPDFName))
	if err != nil {
		fmt.Printf("Failed to rename PDF file: %v\n", err)
		return
	}

	fmt.Printf("PDF file has been successfully compiled and renamed to: %s\n", newPDFName)
}

func pdf_create(c *gin.Context, jsonData map[string]interface{}) {
	type Response struct {
		State string `json:"state"`
	}
	// LaTeX 代码
	latexSourceCode := gorm.Document_paragraph_query_all_text_sql(jsonData)

	// 文件放置位置
	outputDirectory := "/root/data/download/file/"

	//将latex代码翻译成pdf文件
	compileLatex(latexSourceCode, outputDirectory, jsonData["pdfName"].(string))

	response := Response{
		State: "OK",
	}
	c.JSON(200, response)
}

func Operate_pdf_create(operate string, c *gin.Context, jsonData map[string]interface{}) {
	switch operate {
	case "pdf_create":
		pdf_create(c, jsonData)
	}
}
