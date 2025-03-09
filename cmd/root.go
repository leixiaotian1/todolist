package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "todolist",
	Short: "Todolist 是一个简单的待办事项管理工具",
	Long:  `使用 Golang 和 Cobra 框架实现的命令行工具，用于管理带有提醒、优先级、命名空间等扩展信息的待办任务。`,
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// 添加子命令
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(clearCmd)

	rootCmd.AddCommand(interactiveCmd)
}
