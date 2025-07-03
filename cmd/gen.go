package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var outputPattern string

// Language 表示语言配置
type Language struct {
	Code        string `json:"code"`        // 语言代码，如 "zh", "en"
	Name        string `json:"name"`        // 语言名称，如 "中文", "English"
	DisplayName string `json:"displayName"` // 显示名称，用于链接文本
	File        string `json:"file"`        // 对应的语言文件名
}

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen [template] [language-dir]",
	Short: "根据指定模板生成多语言文件",
	Long: `根据指定模板和语言文件目录生成多语言文件。

示例:
  multilang-gen gen template.html ./langs
  multilang-gen gen template.html ./langs --output "{lang}.html"`,
	Args: cobra.ExactArgs(2),
	RunE: runGen,
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.Flags().StringVarP(&outputPattern, "output", "o", "{lang}.html", "输出文件名模式，{lang} 为语言替代符")
}

func runGen(cmd *cobra.Command, args []string) error {
	templatePath := args[0]
	langDir := args[1]

	// 1. 读取语言索引文件
	languages, err := loadLanguageIndex(langDir)
	if err != nil {
		return fmt.Errorf("读取语言索引失败: %w", err)
	}

	if len(languages) == 0 {
		return fmt.Errorf("在索引文件中未找到任何语言配置")
	}

	fmt.Printf("找到 %d 种语言: ", len(languages))
	for i, lang := range languages {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%s(%s)", lang.DisplayName, lang.Code)
	}
	fmt.Println()

	// 2. 解析模板文件
	tmpl, err := parseTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("解析模板文件失败: %w", err)
	}

	// 3. 生成语言链接
	langLinks := generateLanguageLinksFromIndex(languages)

	// 4. 为每种语言生成文件
	for _, lang := range languages {
		if err := generateLanguageFileFromIndex(tmpl, lang, langLinks, langDir); err != nil {
			return fmt.Errorf("生成语言文件 %s 失败: %w", lang.Code, err)
		}

		outputFile := strings.ReplaceAll(outputPattern, "{lang}", lang.Code)
		fmt.Printf("生成文件: %s (%s)\n", outputFile, lang.DisplayName)
	}

	fmt.Println("多语言文件生成完成!")
	return nil
}

// parseTemplate 使用 template/html 解析模板文件
func parseTemplate(templatePath string) (*template.Template, error) {
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("读取模板文件失败: %w", err)
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(templateContent))
	if err != nil {
		return nil, fmt.Errorf("解析模板失败: %w", err)
	}

	return tmpl, nil
}

// generateLanguageLinksFromIndex 从索引生成语言链接
func generateLanguageLinksFromIndex(languages []Language) map[string]Language {
	links := make(map[string]Language)

	for _, lang := range languages {
		links[lang.Code] = lang
	}

	return links
}

// generateLanguageFileFromIndex 为指定语言生成文件（基于索引）
func generateLanguageFileFromIndex(
	tmpl *template.Template,
	currentLang Language,
	allLangs map[string]Language,
	langDir string,
) error {
	// 读取当前语言的数据文件
	langData, err := loadLanguageDataFromFile(langDir, currentLang.File)
	if err != nil {
		return fmt.Errorf("加载语言数据失败: %w", err)
	}

	// 生成其他语言链接的HTML
	var langLinksHTML strings.Builder
	for _, lang := range allLangs {
		if lang.Code != currentLang.Code {
			outputFile := strings.ReplaceAll(outputPattern, "{lang}", lang.Code)
			langLinksHTML.WriteString(fmt.Sprintf(`<a href="%s">%s</a> `, outputFile, lang.DisplayName))
		}
	}

	// 准备模板数据
	templateData := struct {
		Language  string
		LangCode  string
		LangName  string
		Data      map[string]interface{}
		LangLinks string
	}{
		Language:  currentLang.Code,
		LangCode:  currentLang.Code,
		LangName:  currentLang.DisplayName,
		Data:      langData,
		LangLinks: strings.TrimSpace(langLinksHTML.String()),
	}

	// 生成输出文件名
	outputFile := strings.ReplaceAll(outputPattern, "{lang}", currentLang.Code)

	// 创建输出文件
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer outFile.Close()

	// 执行模板渲染
	if err := tmpl.Execute(outFile, templateData); err != nil {
		return fmt.Errorf("模板渲染失败: %w", err)
	}

	// 处理 {__LANG_LINKS__} 替换
	return replaceLanguageLinks(outputFile, langLinksHTML.String())
}

// loadLanguageData 加载语言数据文件（仅支持JSON格式）
func loadLanguageData(langDir, language string) (map[string]interface{}, error) {
	// 只支持 JSON 格式
	filePath := filepath.Join(langDir, language+".json")
	if _, err := os.Stat(filePath); err == nil {
		// 文件存在，解析JSON文件
		return parseLanguageFile(filePath, ".json")
	}

	// 如果没有找到数据文件，返回错误
	return nil, fmt.Errorf("未找到语言文件 %s.json", language)
}

// loadLanguageIndex 加载语言索引文件
func loadLanguageIndex(langDir string) ([]Language, error) {
	indexPath := filepath.Join(langDir, "index.json")

	var content []byte
	var err error

	// 首先尝试读取外部文件
	if _, statErr := os.Stat(indexPath); statErr == nil {
		content, err = os.ReadFile(indexPath)
		if err != nil {
			return nil, fmt.Errorf("读取外部索引文件失败: %w", err)
		}
	} else {
		return nil, fmt.Errorf("外部索引文件不存在，请先创建 index.json")
	}

	var languages []Language
	if err := json.Unmarshal(content, &languages); err != nil {
		return nil, fmt.Errorf("解析索引文件失败: %w", err)
	}

	return languages, nil
}

// parseLanguageFile 解析语言文件（仅支持JSON格式）
func parseLanguageFile(filePath, ext string) (map[string]interface{}, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})

	switch ext {
	case ".json":
		// JSON 格式解析
		if err := json.Unmarshal(content, &result); err != nil {
			return nil, fmt.Errorf("解析 JSON 文件失败: %w", err)
		}
	default:
		return nil, fmt.Errorf("不支持的文件格式 %s，仅支持 .json 格式", ext)
	}

	return result, nil
}

// replaceLanguageLinks 替换文件中的 {__LANG_LINKS__} 占位符
func replaceLanguageLinks(filename, langLinks string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// 替换 {__LANG_LINKS__} 占位符
	updatedContent := strings.ReplaceAll(string(content), "{__LANG_LINKS__}", langLinks)

	return os.WriteFile(filename, []byte(updatedContent), 0o644)
}

// loadLanguageDataFromFile 从指定文件加载语言数据
func loadLanguageDataFromFile(langDir, filename string) (map[string]interface{}, error) {
	filePath := filepath.Join(langDir, filename)

	// 获取文件扩展名
	ext := filepath.Ext(filename)

	return parseLanguageFile(filePath, ext)
}
