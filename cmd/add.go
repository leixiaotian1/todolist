package cmd

import (
	"fmt"
	"strings"
	"time"

	todolistpb "todolist/proto"
	"todolist/tasks"

	"github.com/spf13/cobra"
)

var (
	addPriority  string
	addNamespace string
	addTags      string
	addDue       string
	addRemind    string
)

var addCmd = &cobra.Command{
	Use:   "add [task description]",
	Short: "添加一项任务",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		description := args[0]

		// 解析优先级参数
		var priority todolistpb.Priority
		switch strings.ToLower(addPriority) {
		case "low":
			priority = todolistpb.Priority_LOW
		case "medium":
			priority = todolistpb.Priority_MEDIUM
		case "high":
			priority = todolistpb.Priority_HIGH
		default:
			priority = todolistpb.Priority_PRIORITY_UNSPECIFIED
		}

		// 解析标签，逗号分隔
		var tags []string
		if addTags != "" {
			tags = strings.Split(addTags, ",")
			for i, tag := range tags {
				tags[i] = strings.TrimSpace(tag)
			}
		}

		// 解析截止时间
		var dueAt *time.Time
		if addDue != "" {
			t, err := time.Parse("2006-01-02 15:04", addDue)
			if err != nil {
				fmt.Println("截止时间格式错误，请使用 'YYYY-MM-DD HH:MM'")
				return
			}
			dueAt = &t
		}

		// 解析提醒时间
		var remindAt *time.Time
		if addRemind != "" {
			t, err := time.Parse("2006-01-02 15:04", addRemind)
			if err != nil {
				fmt.Println("提醒时间格式错误，请使用 'YYYY-MM-DD HH:MM'")
				return
			}
			remindAt = &t
		}

		err := tasks.AddTask(description, priority, addNamespace, tags, dueAt, remindAt)
		if err != nil {
			fmt.Println("添加任务失败：", err)
			return
		}
		fmt.Println("任务添加成功：", description)
	},
}

func init() {

	addCmd.Flags().StringVar(&addPriority, "priority", "unspecified", "任务优先级：low, medium, high")
	addCmd.Flags().StringVar(&addNamespace, "namespace", "", "任务所属的命名空间")
	addCmd.Flags().StringVar(&addTags, "tags", "", "任务标签，多个标签用逗号分隔")
	addCmd.Flags().StringVar(&addDue, "due", "", "任务截止时间，格式: 'YYYY-MM-DD HH:MM'")
	addCmd.Flags().StringVar(&addRemind, "remind", "", "任务提醒时间，格式: 'YYYY-MM-DD HH:MM'")
}
