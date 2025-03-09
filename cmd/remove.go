package cmd

import (
	"fmt"
	"strconv"

	"todolist/tasks"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [task index]",
	Short: "删除指定的任务",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("任务索引必须为数字")
			return
		}
		err = tasks.RemoveTask(index)
		if err != nil {
			fmt.Println("删除任务失败：", err)
			return
		}
		fmt.Println("任务删除成功")
	},
}
