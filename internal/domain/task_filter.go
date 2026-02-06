package domain

type TaskFilter struct {
	Status   *TaskStatus
	Assignee *string
	Limit    int
	Offset   int
}
