package app

import (
	db "Okedex/internal/db"

	"github.com/rivo/tview"
)

var app *tview.Application

func Run() int {
	var err error
	if err = db.Init(); err != nil {
		panic(err)
	}
	defer db.Close()

	db.AddPokemon("Pikachu")

	app = tview.NewApplication()
	list := tview.NewList().
		AddItem("Inventory", "Browse through your inventory", 'a', finder).
		AddItem("Info", "Get info about various items/pokemon", 'b', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}

	return 0
}

func start() tview.Primitive {
	return nil
}

func finder() {
	// Create the basic objects.
	categories := tview.NewList().ShowSecondaryText(false)
	categories.SetBorder(true).SetTitle("Categories")
	info := tview.NewTable()
	info.SetBorder(true).SetTitle("Info")
	articles := tview.NewList()
	articles.ShowSecondaryText(false).
		SetDoneFunc(func() {
			articles.Clear()
			info.Clear()
			app.SetFocus(categories)
		})
	articles.SetBorder(true)

	// Create the layout.
	flex := tview.NewFlex().
		AddItem(categories, 0, 1, true).
		AddItem(articles, 0, 1, false).
		AddItem(info, 0, 3, false)

	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}
