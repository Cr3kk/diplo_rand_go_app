package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Countries struct {
	PrevCountries map[string]string `json:"prevCountries"`
}

func logic(prevCountries map[string]string) map[string]string {
	names := []string{"Ben", "Jan-Jan", "Lock", "Koen", "Niels", "Casper", "Wouter"}
	if len(names) != len(prevCountries) {
		fmt.Println("Error: Names and countries count mismatch")
		return prevCountries
	}

	rand.Seed(time.Now().UnixNano())
	for {
		rand.Shuffle(len(names), func(i, j int) {
			names[i], names[j] = names[j], names[i]
		})

		newCountries := make(map[string]string)
		i := 0
		for _, country := range prevCountries {
			newCountries[names[i]] = country
			i++
		}

		isDifferent := false
		for key, value := range newCountries {
			if prevCountries[key] != value {
				isDifferent = true
				break
			}
		}
		if isDifferent {
			return newCountries
		}
	}
}

func main() {
	file, err := os.Open("setup.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var data Countries
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	newCountries := logic(data.PrevCountries)

	myApp := app.New()
	myWindow := myApp.NewWindow("Randomized Countries")

	var countryOptions []string
	for _, country := range data.PrevCountries {
		countryOptions = append(countryOptions, country)
	}

	prevCountriesSelector := container.NewVBox()
	for name, prevCountry := range data.PrevCountries {
		selector := widget.NewSelect(countryOptions, func(selectedCountry string) {
			newCountries[name] = selectedCountry
		})
		selector.SetSelected(prevCountry)
		prevCountriesSelector.Add(container.NewHBox(
			widget.NewLabel(fmt.Sprintf("%s:", name)),
			selector,
		))
	}

	newCountriesText := "New Allocations:\n"
	for name, country := range newCountries {
		newCountriesText += fmt.Sprintf("%s: %s\n", name, country)
	}

	newAllocationsLabel := widget.NewLabel(newCountriesText)

	saveButton := widget.NewButton("Save", func() {
		file, err := os.Create("setup.json")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(Countries{PrevCountries: newCountries}); err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}
		fmt.Println("New countries saved to setup.json")
	})

	runLogicButton := widget.NewButton("Run Logic", func() {
		newCountries = logic(data.PrevCountries)

		newCountriesText := ""
		for name, country := range newCountries {
			newCountriesText += fmt.Sprintf("%s: %s\n", name, country)
		}

		newAllocationsLabel.SetText(newCountriesText)
	})

	leftContainer := container.NewVBox(
		runLogicButton,
		newAllocationsLabel,
	)

	saveButtonContainer := container.NewHBox(
		layout.NewSpacer(),
		saveButton,
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
	)
	saveButtonContainer.Objects[1].Resize(fyne.NewSize(0, 100))

	prevCountriesSelector.Add(saveButtonContainer)

	content := container.NewHBox(
		container.NewVBox(
			widget.NewLabel("Select previous countries for each person:"),
			prevCountriesSelector,
		),
		layout.NewSpacer(),
		leftContainer,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(600, 400))
	myWindow.ShowAndRun()
}
