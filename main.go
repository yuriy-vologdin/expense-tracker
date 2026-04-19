package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const filename = "expenses.json"

type expense struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Amount      int       `json:"amount"`
	Description string    `json:"description"`
}

type storage struct {
	Nextid   int       `json:"nextid"`
	Expenses []expense `json:"expenses"`
}

func save(s storage) error {
	jsonData, err := json.MarshalIndent(s, "", "  ")
	if err != nil {

		return fmt.Errorf("Marshaling error: %v", err)
	}
	os.WriteFile(filename, jsonData, 0644)

	return nil
}

func load() (storage, error) {
	_, err := os.Stat(filename)
	if err != nil {

		return storage{Nextid: 1}, nil
	}
	data, err := os.ReadFile(filename)
	if err != nil {

		return storage{Nextid: 1}, fmt.Errorf("File reading error: %v", err)
	}
	var s storage
	err = json.Unmarshal(data, &s)
	if err != nil {

		return storage{Nextid: 1}, fmt.Errorf("File unmarshal error: %v", err)
	}
	return s, nil
}

func main() {
	data, err := load()
	if err != nil {
		fmt.Println(err)
	}
	if len(os.Args) > 1 && os.Args[1] == "list" {
		for _, exp := range data.Expenses {
			fmt.Printf("%-4d %-20s %-5d %s\n", exp.ID,
				exp.Date.Format("02/01/2006 15:04"), exp.Amount, exp.Description)
		}
	}

}
