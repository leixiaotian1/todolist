package cmd

import (
	"os"
	"strconv"

	todolistpb "todolist/proto"
	"todolist/tasks"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [n]",
	Short: "列出任务。如果提供数字 n，则只显示最新 n 条任务",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskList, err := tasks.LoadTasks()
		if err != nil {
			cmd.Println("读取任务失败：", err)
			return
		}
		if len(taskList.Tasks) == 0 {
			cmd.Println("当前没有任务。")
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Description", "Priority", "Namespace", "Status", "CreatedAt", "DueAt", "RemindAt", "Tags"})
		table.SetBorder(true)
		table.SetCenterSeparator("|")

		var tasksToDisplay []*todolistpb.Task
		if len(args) == 1 {
			count, err := strconv.Atoi(args[0])
			if err != nil || count < 1 {
				cmd.Println("参数错误：请提供一个正整数")
				return
			}
			if count > len(taskList.Tasks) {
				count = len(taskList.Tasks)
			}
			// 取最新 count 条任务（任务按添加顺序存储，最新任务在列表末尾）
			tasksToDisplay = taskList.Tasks[len(taskList.Tasks)-count:]
		} else {
			tasksToDisplay = taskList.Tasks
		}

		for _, t := range tasksToDisplay {
			createdAt := ""
			if t.CreatedAt != nil {
				createdAt = t.CreatedAt.AsTime().Format("2006-01-02 15:04:05")
			}
			dueAt := ""
			if t.DueAt != nil {
				dueAt = t.DueAt.AsTime().Format("2006-01-02 15:04:05")
			}
			remindAt := ""
			if t.RemindAt != nil {
				remindAt = t.RemindAt.AsTime().Format("2006-01-02 15:04:05")
			}
			tags := ""
			if len(t.Tags) > 0 {
				for i, tag := range t.Tags {
					if i > 0 {
						tags += ", "
					}
					tags += tag
				}
			}
			table.Append([]string{
				strconv.FormatInt(t.Id, 10),
				t.Description,
				t.Priority.String(),
				t.Namespace,
				t.Status.String(),
				createdAt,
				dueAt,
				remindAt,
				tags,
			})
		}
		table.Render()
	},
}
