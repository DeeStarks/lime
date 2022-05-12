package main

import (
	"github.com/DeeStarks/lime/configs"
	"github.com/DeeStarks/lime/internal/editor"
	app "github.com/DeeStarks/lime/init"
)

func main() {

	version := app.NewVersion(configs.V1)
	application := app.NewApplication(version)
	screen := application.Screen
	screen.ShowBox() //
	editor := editor.NewEditor(screen)
	editor.Launch()
}