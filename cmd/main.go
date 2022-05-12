package main

import (
	"github.com/DeeStarks/lime/configs"
	app "github.com/DeeStarks/lime/init"
	"github.com/DeeStarks/lime/internal/editor"
)

func main() {

	version := app.NewVersion(configs.V1)
	application := app.NewApplication(version)
	application.Screen.ShowBox() //
	editor := editor.NewEditor(application.Screen)
	editor.Launch()
}
