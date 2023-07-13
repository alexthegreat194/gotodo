package main

import (
	"fmt"
	"io"
	"os"

	"github.com/eiannone/keyboard"
	"github.com/gosuri/uilive"
)

type Task struct {
	title string
	done bool
}

func exitProgram() {
	err := keyboard.Close()
	if err != nil {
		panic(err)
	}
	os.Exit(1)
}

func menu(items []Task) int {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	current := 0
	for {

		// draw menu
		for i, v := range items {
			done := " "
			arrow := " "
	
			if current == i {
				arrow = ">"
			}
			if v.done == true {
				done = "X"
			}
	
			fmt.Printf("%v [%v] %v\n", arrow, done, v.title)
		}

		_, key, err := keyboard.GetKey()
		if err != nil {
			return -1
		}

		if key == keyboard.KeyEsc { return -1 }

		if key == keyboard.KeyArrowUp { current-- }
		if key == keyboard.KeyArrowDown { current++ }
		if key == keyboard.KeyEnter {
			status := items[current].done
			items[current].done = !status
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
		writer.Stop()
		for range writerList {
			writer.Stop()
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
			endWriters()
			return -1
		}

		if key == keyboard.KeyEsc { 
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

func main() {
	buffer := []Task{ 
		{title: "one", done: false},
		{title: "two", done: false},
		{title: "three", done: false},
	}
	// menu(buffer)
	fmt.Println(menuSelect([]string{"one", "two", "three"}))
	fmt.Println(buffer)
}