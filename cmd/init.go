package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [target-dir]",
	Short: "初始化项目模板文件",
	Long: `初始化完整的项目模板文件到指定目录。

这个命令会创建一个完整的多语言项目模板，包括：
- index.tmpl: HTML 模板文件
- manifest.json: 站点基础配置
- langs/index.json: 语言索引
- langs/zh-CN.json: 中文语言包示例
- langs/en-US.json: 英文语言包示例

示例:
  multilang-gen init ./my-project
  multilang-gen init .`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	// 确保目标目录存在
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 创建 langs 子目录
	langsDir := filepath.Join(targetDir, "langs")
	if err := os.MkdirAll(langsDir, 0o755); err != nil {
		return fmt.Errorf("创建语言目录失败: %w", err)
	}

	// 导出 manifest.json
	manifest := map[string]interface{}{
		"baseURL":     "https://example.com",
		"siteName":    "My Website",
		"author":      "Website Author",
		"description": "A multilingual website",
		"version":     "1.0.0",
	}

	manifestData, _ := json.MarshalIndent(manifest, "", "    ")
	manifestPath := filepath.Join(targetDir, "manifest.json")
	if err := os.WriteFile(manifestPath, manifestData, 0o644); err != nil {
		return fmt.Errorf("导出 manifest.json 失败: %w", err)
	}

	// 导出语言索引
	languageIndex := []map[string]interface{}{
		{
			"code":        "zh",
			"name":        "中文",
			"displayName": "中文",
			"file":        "zh-CN.json",
		},
		{
			"code":        "en",
			"name":        "English",
			"displayName": "English",
			"file":        "en-US.json",
		},
	}

	indexData, _ := json.MarshalIndent(languageIndex, "", "    ")
	indexPath := filepath.Join(langsDir, "index.json")
	if err := os.WriteFile(indexPath, indexData, 0o644); err != nil {
		return fmt.Errorf("导出语言索引失败: %w", err)
	}

	// 导出中文语言包
	zhData := map[string]interface{}{
		"title":             "我的网站",
		"subtitle":          "基于模板的多语言网站",
		"description":       "这是一个多语言网站示例",
		"welcome":           "欢迎使用",
		"language_switcher": "语言切换",
		"switch_to":         "切换到",
		"site_info":         "站点信息",
		"site_name":         "站点名称",
		"version":           "版本",
		"base_url":          "基础URL",
		"author":            "作者",
		"current_language":  "当前语言",
		"language_code":     "语言代码",
		"language_name":     "语言名称",
		"footer_text":       "页脚文本",
	}

	zhJSON, _ := json.MarshalIndent(zhData, "", "    ")
	zhPath := filepath.Join(langsDir, "zh-CN.json")
	if err := os.WriteFile(zhPath, zhJSON, 0o644); err != nil {
		return fmt.Errorf("导出中文语言包失败: %w", err)
	}

	// 导出英文语言包
	enData := map[string]interface{}{
		"title":             "My Website",
		"subtitle":          "Template-based multilingual website",
		"description":       "This is a multilingual website example",
		"welcome":           "Welcome",
		"language_switcher": "Language Switcher",
		"switch_to":         "Switch to",
		"site_info":         "Site Information",
		"site_name":         "Site Name",
		"version":           "Version",
		"base_url":          "Base URL",
		"author":            "Author",
		"current_language":  "Current Language",
		"language_code":     "Language Code",
		"language_name":     "Language Name",
		"footer_text":       "Footer Text",
	}

	enJSON, _ := json.MarshalIndent(enData, "", "    ")
	enPath := filepath.Join(langsDir, "en-US.json")
	if err := os.WriteFile(enPath, enJSON, 0o644); err != nil {
		return fmt.Errorf("导出英文语言包失败: %w", err)
	}

	fmt.Printf("项目模板已初始化到: %s\n", targetDir)
	fmt.Println("包含文件:")
	fmt.Println("  - manifest.json")
	fmt.Println("  - langs/index.json")
	fmt.Println("  - langs/zh-CN.json")
	fmt.Println("  - langs/en-US.json")
	fmt.Println("\n运行 'multilang-gen gen .' 来生成多语言文件")

	return nil
}
