package cmd

import (
	"fmt"

	"todolist/tasks"

	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "清空所有任务",
	Run: func(cmd *cobra.Command, args []string) {
		err := tasks.ClearTasks()
		if err != nil {
			fmt.Println("清空任务失败：", err)
			return
		}
		fmt.Println("所有任务已清空")
	},
}
