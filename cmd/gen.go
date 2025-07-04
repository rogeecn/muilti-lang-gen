package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	outputPattern string
	langCodes     []string
)

// Language 表示语言配置
type Language struct {
	Code        string `json:"code"`        // 语言代码，如 "zh", "en"
	Name        string `json:"name"`        // 语言名称，如 "中文", "English"
	DisplayName string `json:"displayName"` // 显示名称，用于链接文本
	File        string `json:"file"`        // 对应的语言文件名
	URL         string `json:"-"`           // 生成的文件URL，运行时添加
	Current     bool   `json:"-"`           // 是否为当前语言，运行时添加
}

// Manifest 表示站点基础配置
type Manifest struct {
	BaseURL     string `json:"baseURL"`
	SiteName    string `json:"siteName"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen [directory]",
	Short: "根据指定目录生成多语言文件",
	Long: `根据指定目录生成多语言文件。目录结构应包含：
- index.tmpl: 模板文件
- langs/index.json: 语言索引文件
- langs/*.json: 语言数据文件
- manifest.json: 站点基础配置（可选）
- outputs/: 输出目录

示例:
  multilang-gen gen .
  multilang-gen gen ./project --output "{lang}.html"
  multilang-gen gen . --lang zh
  multilang-gen gen . --lang zh,en
  multilang-gen gen . --lang zh --lang en --output "page-{lang}.html"`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGen,
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.Flags().StringVarP(&outputPattern, "output", "o", "{lang}.html", "输出文件名模式，{lang} 为语言替代符")
	genCmd.Flags().
		StringSliceVarP(&langCodes, "lang", "l", []string{}, "只生成指定语言代码的文件，支持多个语言（如: zh,en 或 --lang zh --lang en）")
}

func runGen(cmd *cobra.Command, args []string) error {
	// 确定项目目录
	projectDir := "."
	if len(args) > 0 {
		projectDir = args[0]
	}

	// 检查项目目录是否存在
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		return fmt.Errorf("项目目录不存在: %s", projectDir)
	}

	// 固定的文件路径
	templatePath := filepath.Join(projectDir, "index.tmpl")
	langDir := filepath.Join(projectDir, "langs")
	outputDir := filepath.Join(projectDir, "outputs")
	manifestPath := filepath.Join(projectDir, "manifest.json")

	// 检查必需文件是否存在
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return fmt.Errorf("模板文件不存在: %s", templatePath)
	}

	if _, err := os.Stat(filepath.Join(langDir, "index.json")); os.IsNotExist(err) {
		return fmt.Errorf("语言索引文件不存在: %s", filepath.Join(langDir, "index.json"))
	}

	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 1. 读取 manifest.json（可选）
	manifest, err := loadManifest(manifestPath)
	if err != nil {
		fmt.Printf("警告: 无法读取 manifest.json，将使用默认值: %v\n", err)
		manifest = &Manifest{
			BaseURL:     "",
			SiteName:    "Website",
			Author:      "",
			Description: "",
			Version:     "1.0.0",
		}
	}

	// 2. 读取语言索引文件
	languages, err := loadLanguageIndex(langDir)
	if err != nil {
		return fmt.Errorf("读取语言索引失败: %w", err)
	}

	if len(languages) == 0 {
		return fmt.Errorf("在索引文件中未找到任何语言配置")
	}

	// 如果指定了语言代码，过滤语言列表
	if len(langCodes) > 0 {
		var filteredLanguages []Language
		langCodeSet := make(map[string]bool)

		// 创建语言代码集合便于查找
		for _, code := range langCodes {
			langCodeSet[code] = true
		}

		// 过滤出指定的语言
		for _, lang := range languages {
			if langCodeSet[lang.Code] {
				filteredLanguages = append(filteredLanguages, lang)
				delete(langCodeSet, lang.Code) // 标记已找到
			}
		}

		// 检查是否有未找到的语言代码
		if len(langCodeSet) > 0 {
			var notFound []string
			for code := range langCodeSet {
				notFound = append(notFound, code)
			}
			return fmt.Errorf("未找到以下语言代码的配置: %v", notFound)
		}

		if len(filteredLanguages) == 0 {
			return fmt.Errorf("指定的语言代码都未找到对应的配置")
		}

		languages = filteredLanguages
		fmt.Printf("只生成指定语言 (%d 种): ", len(languages))
		for i, lang := range languages {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%s(%s)", lang.DisplayName, lang.Code)
		}
		fmt.Println()
	} else {
		fmt.Printf("找到 %d 种语言: ", len(languages))
		for i, lang := range languages {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%s(%s)", lang.DisplayName, lang.Code)
		}
		fmt.Println()
	}

	// 3. 解析模板文件
	tmpl, err := parseTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("解析模板文件失败: %w", err)
	}

	// 4. 生成语言链接（使用过滤后的语言列表）
	langLinks := generateLanguageLinksFromIndex(languages)

	// 5. 为每种语言生成文件
	for _, lang := range languages {
		if err := generateLanguageFileFromIndex(tmpl, lang, langLinks, langDir, outputDir, manifest); err != nil {
			return fmt.Errorf("生成语言文件 %s 失败: %w", lang.Code, err)
		}

		outputFile := strings.ReplaceAll(outputPattern, "{lang}", lang.Code)
		fmt.Printf("生成文件: %s (%s)\n", filepath.Join(outputDir, outputFile), lang.DisplayName)
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
	langLinks map[string]Language,
	langDir string,
	outputDir string,
	manifest *Manifest,
) error {
	// 读取当前语言的数据文件
	langData, err := loadLanguageDataFromFile(langDir, currentLang.File)
	if err != nil {
		return fmt.Errorf("加载语言数据失败: %w", err)
	}

	// 将语言数据合并到当前语言结构中
	currentLang.URL = strings.ReplaceAll(outputPattern, "{lang}", currentLang.Code)
	currentLang.Current = true

	// 生成语言链接列表（使用传入的语言链接）
	var allLangLinks []Language
	for _, lang := range langLinks {
		// 为每个语言添加输出文件路径和当前状态
		lang.URL = strings.ReplaceAll(outputPattern, "{lang}", lang.Code)
		lang.Current = (lang.Code == currentLang.Code)
		allLangLinks = append(allLangLinks, lang)
	}

	// 序列化 I18N 数据为 JSON
	i18nJson, err := json.Marshal(langData)
	if err != nil {
		return fmt.Errorf("序列化语言数据失败: %w", err)
	}

	// 准备模板数据
	templateData := struct {
		Lang      Language
		LangLinks []Language
		I18N      map[string]interface{}
		I18NJson  string
		Base      *Manifest
	}{
		Lang:      currentLang,
		LangLinks: allLangLinks,
		I18N:      langData,
		I18NJson:  string(i18nJson),
		Base:      manifest,
	}

	// 生成输出文件名
	outputFile := strings.ReplaceAll(outputPattern, "{lang}", currentLang.Code)
	outputPath := filepath.Join(outputDir, outputFile)

	// 创建输出文件
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer outFile.Close()

	// 执行模板渲染
	if err := tmpl.Execute(outFile, templateData); err != nil {
		return fmt.Errorf("模板渲染失败: %w", err)
	}

	return nil
}

// loadManifest 加载站点配置文件
func loadManifest(manifestPath string) (*Manifest, error) {
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("读取 manifest.json 失败: %w", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(content, &manifest); err != nil {
		return nil, fmt.Errorf("解析 manifest.json 失败: %w", err)
	}

	return &manifest, nil
}

// loadLanguageIndex 加载语言索引文件
func loadLanguageIndex(langDir string) ([]Language, error) {
	indexPath := filepath.Join(langDir, "index.json")

	content, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("读取语言索引文件失败: %w", err)
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

// loadLanguageDataFromFile 从指定文件加载语言数据
func loadLanguageDataFromFile(langDir, filename string) (map[string]interface{}, error) {
	filePath := filepath.Join(langDir, filename)

	// 获取文件扩展名
	ext := filepath.Ext(filename)

	return parseLanguageFile(filePath, ext)
}
