package app

import (
	db "Okedex/internal/db"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Pokemon struct {
	Name      string                         `json:"name"`
	Height    int                            `json:"height"`
	Weight    int                            `json:"weight"`
	Xp        int                            `json:"base_experience"`
	Abilities []map[string]map[string]string `json:"abilities"`
	HeldItems []map[string]map[string]string `json:"held_items"`
}

type Item struct {
	Cost          int                 `json:"cost"`
	FlingPower    int                 `json:"fling_power"`
	FlingEffect   map[string]string   `json:"fling_effect"`
	Attributes    []map[string]string `json:"attributes"`
	EffectEntries []map[string]string `json:"effect_entries"`
}

var app *tview.Application

func Run() int {
	var err error
	if err = db.Init(); err != nil {
		log.Fatal("Failed to init db")
	}
	defer db.Close()

	// Testing purposes
	/*
		db.AddItem("antidote")
		db.AddItem("apple")
		db.AddPokemon("Pikachu")
		db.AddPokemon("Pichu")
		db.AddPokemon("Bulbasaur")
	*/

	app = tview.NewApplication()
	startMenu()

	return 0
}

func startMenu() {
	var menu *tview.List
	menu = tview.NewList().
		AddItem("Inventory", "Browse through your inventory", 'a', func() { inventorySection(menu) }).
		AddItem("Add", "Add to your inventory", 'b', func() { addSection(menu) }).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	if err := app.SetRoot(menu, true).SetFocus(menu).EnableMouse(true).Run(); err != nil {
		log.Fatalf("Failed to run menu: %v", err)
	}
}

func addSection(menu *tview.List) {
	var flex *tview.Flex
	name := ""
	t := ""

	flex_log := tview.NewFlex().SetDirection(tview.FlexRow)
	form := tview.NewForm().
		AddDropDown("Type", []string{"Pokemon", "Item"}, 0, func(option string, optionIndex int) { t = option }).
		AddInputField("Name", "", 20, nil, func(text string) { name = text }).
		AddButton("Add", func() { addToInventory(name, t, flex_log) })
	form.SetBorder(true).SetTitle("Add to inventory").SetTitleAlign(tview.AlignCenter)
	form.SetCancelFunc(func() {
		if err := app.SetRoot(menu, true).EnableMouse(true).SetFocus(menu).Run(); err != nil {
			log.Fatalf("Failed to run menu: %v", err)
		}
	})

	flex = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(form, 0, 9, true)
	flex.AddItem(flex_log, 5, 1, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		log.Fatalf("Failed to run form: %v", err)
	}
}

func addToInventory(name string, t string, flex *tview.Flex) {
	resp, err := http.Get("https://pokeapi.co/api/v2/" + t + "/" + name)
	if err != nil {
		log.Fatalf("Error creating http request: %v", err)
	}
	defer resp.Body.Close()

	flex.Clear()
	if resp.StatusCode != http.StatusOK {
		flex.AddItem(tview.NewTextView().SetText(t+" '"+name+"' not found. Can't add to inventory.").SetTextColor(tcell.ColorRed), 1, 1, false)
		return
	}

	flex.AddItem(tview.NewTextView().SetText(t+" '"+name+"' found. Succesfully added to inventory").SetTextColor(tcell.ColorGreen), 1, 1, false)

	switch t {
	case "Pokemon":
		db.AddPokemon(name)
	case "Item":
		db.AddItem(name)
	}
}

func inventorySection(menu *tview.List) {
	// Create the basic objects.
	categories := tview.NewList().ShowSecondaryText(false)
	categories.SetBorder(true).SetTitle("Categories")
	categories.SetDoneFunc(func() {
		if err := app.SetRoot(menu, true).EnableMouse(true).SetFocus(menu).Run(); err != nil {
			log.Fatalf("Failed to run menu: %v", err)
		}
	})

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

	// Add items into the objects
	categories.AddItem("Pokemon", "", 0, func() { onCategorySelect(articles, info, pokemonCategory, db.GetOwnedPokemon()) })
	categories.AddItem("Items", "", 0, func() { onCategorySelect(articles, info, ItemCategory, db.GetOwnedItems()) })

	// Create the layout.
	flex := tview.NewFlex().
		AddItem(categories, 0, 1, true).
		AddItem(articles, 0, 1, false).
		AddItem(info, 0, 3, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).SetFocus(flex).Run(); err != nil {
		log.Fatalf("Failed to run flex: %v", err)
	}
}

func onCategorySelect(articles *tview.List, info *tview.Table,
	categoryFunc func(articles *tview.List, info *tview.Table, articleList *map[string]int), articleList *map[string]int) {

	articles.Clear()
	info.Clear()

	info.SetCell(0, 0, &tview.TableCell{Text: "Waiting for API"})
	go categoryFunc(articles, info, articleList)
}

func pokemonCategory(articles *tview.List, info *tview.Table, articleList *map[string]int) {
	articles.SetChangedFunc(func(i int, articleName string, t string, s rune) {
		var pokemon Pokemon

		resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + strings.Split(articleName, "x ")[1])
		if err != nil {
			log.Fatalf("Error creating http request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response: %v", err)
		}

		json.Unmarshal(body, &pokemon)

		info.Clear()
		info.SetCell(0, 0, &tview.TableCell{Text: "Height: ", Align: tview.AlignLeft})
		info.SetCell(0, 1, &tview.TableCell{Text: strconv.Itoa(pokemon.Height), Align: tview.AlignLeft})
		info.SetCell(1, 0, &tview.TableCell{Text: "Weight: ", Align: tview.AlignLeft})
		info.SetCell(1, 1, &tview.TableCell{Text: strconv.Itoa(pokemon.Weight), Align: tview.AlignLeft})
		info.SetCell(2, 0, &tview.TableCell{Text: "XP: ", Align: tview.AlignLeft})
		info.SetCell(2, 1, &tview.TableCell{Text: strconv.Itoa(pokemon.Xp), Align: tview.AlignLeft, Color: tcell.ColorYellow})
		info.SetCell(3, 0, &tview.TableCell{Text: "Abilities: ", Align: tview.AlignLeft})

		var end int
		for i, ability := range pokemon.Abilities {
			info.SetCell(4+i, 0, &tview.TableCell{Text: "(" + strconv.Itoa(i) + ") ", Align: tview.AlignRight})
			info.SetCell(4+i, 1, &tview.TableCell{
				Text:  ability["ability"]["name"],
				Align: tview.AlignLeft,
				Color: tcell.ColorGreen,
			})
			end = 4 + i
		}
		info.SetCell(end+1, 0, &tview.TableCell{Text: "Held items: ", Align: tview.AlignLeft})
		for i, item := range pokemon.HeldItems {
			info.SetCell(end+i+2, 0, &tview.TableCell{Text: "(" + strconv.Itoa(i) + ") ", Align: tview.AlignRight})
			info.SetCell(end+i+2, 1, &tview.TableCell{
				Text:  item["item"]["name"],
				Align: tview.AlignLeft,
				Color: tcell.ColorBlue,
			})
		}
	})

	log.Print(articleList)
	ownedPokemon := articleList
	for name, count := range *ownedPokemon {
		articles.AddItem(strconv.Itoa(count)+"x "+name, "", 0, func() {})
	}

	app.SetFocus(articles)
	articles.SetCurrentItem(0)
	app.Draw()
}

func ItemCategory(articles *tview.List, info *tview.Table, articleList *map[string]int) {
	articles.SetChangedFunc(func(i int, articleName string, t string, s rune) {
		var item Item

		resp, err := http.Get("https://pokeapi.co/api/v2/item/" + strings.Split(articleName, "x ")[1])
		if err != nil {
			log.Fatalf("Error creating http request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response: %v", err)
		}

		info.Clear()
		json.Unmarshal(body, &item)

		info.SetCell(0, 0, &tview.TableCell{Text: "Cost: ", Align: tview.AlignLeft})
		info.SetCell(0, 1, &tview.TableCell{Text: strconv.Itoa(item.Cost), Align: tview.AlignLeft, Color: tcell.ColorYellow})
		info.SetCell(1, 0, &tview.TableCell{Text: "Fling power: ", Align: tview.AlignLeft})
		info.SetCell(1, 1, &tview.TableCell{Text: strconv.Itoa(item.FlingPower), Align: tview.AlignLeft, Color: tcell.ColorYellow})
		info.SetCell(2, 0, &tview.TableCell{Text: "Fling effect: ", Align: tview.AlignLeft})
		info.SetCell(2, 1, &tview.TableCell{Text: item.FlingEffect["name"], Align: tview.AlignLeft, Color: tcell.ColorYellow})
		info.SetCell(3, 0, &tview.TableCell{Text: "Attributes: ", Align: tview.AlignLeft})

		for i, attribute := range item.Attributes {
			info.SetCell(4+i, 0, &tview.TableCell{Text: "(" + strconv.Itoa(i) + ") ", Align: tview.AlignRight})
			info.SetCell(4+i, 1, &tview.TableCell{
				Text:  attribute["name"],
				Align: tview.AlignLeft,
				Color: tcell.ColorGreen,
			})
		}
	})

	ownedItems := articleList
	for name, count := range *ownedItems {
		articles.AddItem(strconv.Itoa(count)+"x "+name, "", 0, func() {})
	}

	app.SetFocus(articles)
	articles.SetCurrentItem(0)
	app.Draw()
}
