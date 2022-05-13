package main

import (
	"fmt"
	"github.com/DeeStarks/lime/configs"
	app "github.com/DeeStarks/lime/init"
	"github.com/DeeStarks/lime/internal/editor"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("No filename specified")
		return
	}
	filepath := args[1]
	// Make sure a filenmae is specified
	var file *os.File
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			// File does not exist
			// Create a new file
			file, err = os.Create(filepath)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()
		} else {
			// Something else went wrong
			fmt.Println(err)
			return
		}
	} else {
		// File exists
		// Open it
		file, err = os.Open(filepath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
	}

	v := app.NewVersion(configs.V1)
	a := app.NewApplication(v)
	a.Screen.ShowBox() //
	e := editor.NewEditor(a.Screen)
	e.Launch(file)
}
