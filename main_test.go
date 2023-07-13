package main

import "testing"

// testing functions
func TestWriteTasksToJson(t *testing.T) {
	tasks := []Task{
		{Title: "Task 1", Done: false},
		{Title: "Task 2", Done: false},
		{Title: "Task 3", Done: false},
	}

	err := saveTasksToJson(tasks)
	if err != nil {
		t.Errorf("Error saving tasks to json: %v", err)
	}
}

func TestReadTasksFromJson(t *testing.T) {
	tasks := []Task{
		{Title: "Task 1", Done: false},
		{Title: "Task 2", Done: false},
		{Title: "Task 3", Done: false},
	}

	err := saveTasksToJson(tasks)
	if err != nil {
		t.Errorf("Error saving tasks to json: %v", err)
	}

	tasks = loadTasksFromJson()
	if len(tasks) == 0 {
		t.Errorf("Expected 0 tasks, got %v", len(tasks))
	}
}

// benchmark loading tasks from json
func BenchmarkLoadTasksFromJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		loadTasksFromJson()
	}
}

