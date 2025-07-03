/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "multilang-gen",
	Short: "多语言文件生成器",
	Long: `多语言文件生成器 - 根据指定模板和语言文件生成对应的多语言页面。

支持功能:
- 初始化项目模板文件
- 根据模板文件生成多语言页面
- 自定义输出文件名格式
- 自动生成语言间链接
- 支持多种语言数据格式

示例:
  multilang-gen init ./my-project
  multilang-gen gen template.html ./langs
  multilang-gen gen template.html ./langs --output "{lang}.html"`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.muilti-lang-gen.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
