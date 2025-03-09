// cmd/interactive.go
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	todolistpb "todolist/proto"
	"todolist/tasks"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "进入交互式任务管理界面",
	Run: func(cmd *cobra.Command, args []string) {
		model, err := tea.NewProgram(initialModel()).Run()
		if err != nil {
			fmt.Println("启动交互模式失败:", err)

			os.Exit(1)
		}
		model.Init()
	},
}

// 自定义列表项类型
type taskItem struct {
	*todolistpb.Task
}

func (i taskItem) Title() string {
	if i.Task == nil {
		return "无效任务"
	}
	return i.Task.Description
}
func (i taskItem) Description() string {
	if i.Task == nil {
		return ""
	}

	// 安全处理时间字段
	formatTime := func(t *timestamppb.Timestamp) string {
		if t == nil {
			return "未设置"
		}
		return t.AsTime().Format("2006-01-02 15:04")
	}

	// 构造状态描述
	status := "待办"
	if i.Status == todolistpb.Status_COMPLETED {
		status = "已完成"
	}

	return fmt.Sprintf(
		"ID: %d | 状态: %s | 优先级: %s | 截止: %s",
		i.Id,
		status,
		strings.ToLower(i.Priority.String()),
		formatTime(i.DueAt),
	)
}
func (i taskItem) FilterValue() string { return i.Description() }

// 主模型
type model struct {
	list  list.Model
	tasks []*todolistpb.Task
}

func initialModel() model {
	// 带错误处理的加载
	taskList, err := tasks.LoadTasks()
	if err != nil {
		log.Printf("任务加载失败: %v", err)
		taskList = &todolistpb.TaskList{} // 返回空列表
	}

	// 转换任务为列表项（带空列表保护）
	items := make([]list.Item, 0, len(taskList.Tasks))
	for _, t := range taskList.Tasks {
		if t != nil { // 过滤空指针
			items = append(items, taskItem{t})
		}
	}

	// 自定义列表样式
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("62")).
		Foreground(lipgloss.Color("212"))

	// 初始化列表
	l := list.New(items, delegate, 80, 20)
	l.Title = "待办事项列表"
	l.Styles.Title = lipgloss.NewStyle().
		MarginLeft(2).
		Foreground(lipgloss.Color("42")).
		Bold(true)
	l.SetStatusBarItemName("任务", "个任务")

	return model{
		list:  l,
		tasks: taskList.Tasks,
	}
}
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			// 切换任务状态
			selectedItem := m.list.SelectedItem().(taskItem)
			for i, t := range m.tasks {
				if t.Id == selectedItem.Id {
					if m.tasks[i].Status == todolistpb.Status_PENDING {
						m.tasks[i].Status = todolistpb.Status_COMPLETED
					} else {
						m.tasks[i].Status = todolistpb.Status_PENDING
					}
					// 保存修改
					if err := tasks.SaveTasks(&todolistpb.TaskList{Tasks: m.tasks}); err != nil {
						fmt.Println("保存失败:", err)
					}
					// 更新列表显示
					items := make([]list.Item, len(m.tasks))
					for j, task := range m.tasks {
						items[j] = taskItem{task}
					}
					m.list.SetItems(items)
				}
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.list.View() + "\n" +
		"  ↑/↓: 导航 | 回车: 切换状态 | q: 退出\n"
}
