package domain

type TaskFilter struct {
	Status   *TaskStatus
	Assignee *string
	Limit    int
	Offset   int
}

func IsValidStatus(s TaskStatus) bool {
	switch s {
	case StatusTodo,
		StatusInProgress,
		StatusDone:
		return true
	default:
		return false
	}
}
