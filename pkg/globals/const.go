package globals

const (
	// task priority
	LOW    = "LOW"
	MEDIUM = "MEDIUM"
	HIGH   = "HIGH"

	//task status
	PENDING   = "PENDING"
	COMPLETED = "COMPLETED"
)

var TaskPriority = map[string]int{
	LOW:    1,
	MEDIUM: 2,
	HIGH:   3,
}

var TaskStatus = map[string]int{
	PENDING:   1,
	COMPLETED: 2,
}

var TaskPriorityReverse = map[int]string{
	1: LOW,
	2: MEDIUM,
	3: HIGH,
}

var TaskStatusReverse = map[int]string{
	1: PENDING,
	2: COMPLETED,
}
