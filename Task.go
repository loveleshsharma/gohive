package GoHive

//represents an executable task
type Task struct {
	executable func()
}

func (t *Task) setTask(function func()) {
	t.executable = function
}

func (t *Task) getTask() (func()) {
	return t.executable
}

func NewTask(fun func()) Task {
	return Task{
		executable: fun,
	}
}