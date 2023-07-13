package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/eiannone/keyboard"
	"github.com/gosuri/uilive"
)

type Task struct {
	Title string "json:\"title\""
	Done bool	"json:\"done\""
}

func exitProgram() {
	err := keyboard.Close()
	if err != nil {
		panic(err)
	}
	os.Exit(1)
}

func menuCheckbox(items []Task) int {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	writer := uilive.New()
	writer.Start()
	writerList := []io.Writer{}
	for range items {
		buffer := writer.Newline()
		writerList = append(writerList, buffer)
	}

	endWriters := func() {
		writer.Stop()
		for range writerList {
			writer.Stop()
		}
	}

	current := 0
	for {

		// draw menu
		for i, v := range items {
			done := " "
			arrow := " "
	
			if current == i {
				arrow = ">"
			}
			if v.Done == true {
				done = "X"
			}
	
			fmt.Fprintf(writerList[i], "%v [%v] %v\n", arrow, done, v.Title)
		}

		_, key, err := keyboard.GetKey()
		if err != nil {
			endWriters()
			panic(err)
		}

		if key == keyboard.KeyEsc {
			return -1 
		}

		if key == keyboard.KeyArrowUp { current-- }
		if key == keyboard.KeyArrowDown { current++ }
		if key == keyboard.KeyEnter {
			status := items[current].Done
			items[current].Done = !status
		}

		if current > len(items) - 1 { current = len(items) - 1 }
		if current < 0 { current = 0 }
	}
}

// menu for selecting options
func menuSelect(options []string) int {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	writer := uilive.New()
	writer.Start()
	writerList := []io.Writer{}
	for range options {
		buffer := writer.Newline()
		writerList = append(writerList, buffer)
	}

	endWriters := func() {
		writer.Flush()
		for _, v := range writerList {
			fmt.Fprintf(v, "\n")
		}
	}

	current := 0
	for {

		// draw menu
		for i, v := range options {
			arrow := " "
	
			if current == i {
				arrow = ">"
			}
			
			fmt.Fprintf(writerList[i], "%v %v\n", arrow, v)
		}

		_, key, err := keyboard.GetKey()
		if err != nil {
			writer.Flush()
			endWriters()
			return -1
		}

		if key == keyboard.KeyEsc { 
			writer.Flush()
			endWriters()
			return -1 
		}

		if key == keyboard.KeyArrowUp { current-- }
		if key == keyboard.KeyArrowDown { current++ }
		if key == keyboard.KeyEnter {
			return current
		}

		if current > len(options) - 1 { current = len(options) - 1 }
		if current < 0 { current = 0 }
	}
}

func menuAddTask() string {
	fmt.Printf("Enter new task: ")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	text = text[:len(text)-1]
	return text
}

func menuDisplay(items []Task) {
	for _, v := range items {
		done := " "
		if v.Done == true {
			done = "X"
		}
		fmt.Printf("[%v] %v\n", done, v.Title)
	}
}

func saveTasksToJson(tasks []Task) {
	file, err := os.Create("tasks.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// serialize to json
	json.NewEncoder(file).Encode(tasks)
}

func loadTasksFromJson() []Task {
	file, err := os.Open("tasks.json")
	if err != nil {
		return []Task{}
	}
	defer file.Close()

	var tasks []Task
	// deserialize from json
	err = json.NewDecoder(file).Decode(&tasks)
	if err != nil {
		panic(err)
	}

	return tasks
}

func main() {
	buffer := loadTasksFromJson()

	done := false
	for !done {
		os.Stdout.WriteString("\x1b[3;J\x1b[H\x1b[2J")
		menuDisplay(buffer)
		option := menuSelect([]string{"Edit Tasks", "Add Task", "Remove Task", "Exit"})
		os.Stdout.WriteString("\x1b[3;J\x1b[H\x1b[2J")
		switch option {
			case 0:
				// edit tasks
				menuCheckbox(buffer)
			case 1:
				// add task
				task := menuAddTask()
				buffer = append(buffer, Task{Title: task, Done: false})
			case 2:
				// remove task
				strings := []string{}
				for _, v := range buffer {
					strings = append(strings, v.Title)
				}
				strings = append(strings, "Cancel")
				index := menuSelect(strings)
				if index != len(strings) - 1 {
					buffer = append(buffer[:index], buffer[index+1:]...)
				}

			case 3:
				// exit
				done = true
			case -1:
				// exit
				done = true
		}
	}

	fmt.Println("Tasks saved to tasks.json...")
	saveTasksToJson(buffer)
	return
}