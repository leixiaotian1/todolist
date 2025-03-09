package tasks

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	todolistpb "todolist/proto"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var tasksFile = filepath.Join(os.Getenv("HOME"), ".todolist.pb")

// LoadTasks 从文件中加载任务列表，如果文件不存在则返回一个空的 TaskList
func LoadTasks() (*todolistpb.TaskList, error) {
	if _, err := os.Stat(tasksFile); os.IsNotExist(err) {
		return &todolistpb.TaskList{}, nil
	}
	data, err := os.ReadFile(tasksFile)
	if err != nil {
		return nil, err
	}
	var taskList todolistpb.TaskList
	if err := proto.Unmarshal(data, &taskList); err != nil {
		return nil, err
	}
	return &taskList, nil
}

// SaveTasks 将任务列表保存到文件
func SaveTasks(taskList *todolistpb.TaskList) error {
	data, err := proto.Marshal(taskList)
	if err != nil {
		return err
	}
	return os.WriteFile(tasksFile, data, 0644)
}

// generateNewID 计算新任务的 id（取现有任务 id 最大值+1）
func generateNewID(taskList *todolistpb.TaskList) int64 {
	var maxID int64 = 0
	for _, t := range taskList.Tasks {
		if t.Id > maxID {
			maxID = t.Id
		}
	}
	return maxID + 1
}

// AddTask 添加一个新任务，dueAt 和 remindAt 为可选时间参数（传入 nil 则为空）
func AddTask(description string, priority todolistpb.Priority, namespace string, tags []string, dueAt, remindAt *time.Time) error {
	taskList, err := LoadTasks()
	if err != nil {
		return err
	}
	now := time.Now()
	newTask := &todolistpb.Task{
		Id:          generateNewID(taskList),
		Description: description,
		CreatedAt:   timestamppb.New(now),
		UpdatedAt:   timestamppb.New(now),
		Priority:    priority,
		Namespace:   namespace,
		Status:      todolistpb.Status_PENDING,
		Tags:        tags,
	}
	if dueAt != nil {
		newTask.DueAt = timestamppb.New(*dueAt)
	}
	if remindAt != nil {
		newTask.RemindAt = timestamppb.New(*remindAt)
	}
	taskList.Tasks = append(taskList.Tasks, newTask)
	return SaveTasks(taskList)
}

// RemoveTask 根据传入的任务序号（从 1 开始）删除任务
func RemoveTask(index int) error {
	taskList, err := LoadTasks()
	if err != nil {
		return err
	}
	if index < 1 || index > len(taskList.Tasks) {
		return errors.New("invalid task index")
	}
	taskList.Tasks = append(taskList.Tasks[:index-1], taskList.Tasks[index:]...)
	return SaveTasks(taskList)
}

// ClearTasks 清空所有任务
func ClearTasks() error {
	taskList := &todolistpb.TaskList{}
	return SaveTasks(taskList)
}
