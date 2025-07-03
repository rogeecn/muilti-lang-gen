package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	_ "embed"
)

//go:embed index.json
var IndexJSON []byte

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export [target-dir]",
	Short: "导出内置的语言索引文件",
	Long: `导出内置的语言索引文件到指定目录。

这个命令会将编译时嵌入的 index.json 文件导出到指定目录，
可以用作自定义语言配置的起始模板。

示例:
  multilang-gen export ./langs
  multilang-gen export .`,
	Args: cobra.MaximumNArgs(1),
	RunE: runExport,
}

func init() {
	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	// 确保目标目录存在
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 导出索引文件
	indexPath := filepath.Join(targetDir, "index.json")
	if err := os.WriteFile(indexPath, IndexJSON, 0o644); err != nil {
		return fmt.Errorf("导出索引文件失败: %w", err)
	}

	fmt.Printf("已导出语言索引文件到: %s\n", indexPath)
	fmt.Println("您可以根据需要修改此文件来自定义支持的语言列表。")

	return nil
}
