package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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

func add(s *storage, amount int, description string) {
	var exp expense
	exp.ID = s.Nextid
	exp.Date = time.Now()
	exp.Amount = amount
	exp.Description = description
	s.Expenses = append(s.Expenses, exp)
	fmt.Printf("Expense added successfully (ID: %d)\n", s.Nextid)
	s.Nextid++
}

func delete(s *storage, id int) {
	newExpences := make([]expense, 0)
	for _, exp := range s.Expenses {
		if exp.ID != id {
			newExpences = append(newExpences, exp)
		}
	}
	s.Expenses = newExpences
}

func main() {
	data, err := load()
	if err != nil {
		fmt.Println(err)
	}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "list":
			for _, exp := range data.Expenses {
				fmt.Printf("%-4d %-12s %-5d %s\n", exp.ID,
					exp.Date.Format("02-01-2006"), exp.Amount, exp.Description)
			}
		case "add":
			if len(os.Args) != 4 {
				fmt.Printf("Error: invalid number of arguments, use add \"description\" amount\n")
			} else {
				amount, err := strconv.Atoi(os.Args[3])
				if err != nil {
					fmt.Printf("Error: incorrect amount format, amount must be integer\n")
				}
				description := os.Args[2]
				add(&data, amount, description)
				save(data)
			}
		case "summary":
			var sum int
			for _, exp := range data.Expenses {
				sum += exp.Amount
			}
			fmt.Printf("Total expences: %d\n", sum)
		case "delete":
			if len(os.Args) != 3 {
				fmt.Printf("Error: invalid number of arguments, use delete id\n")
			} else {
				id, err := strconv.Atoi(os.Args[2])
				if err != nil {
					fmt.Printf("Error: incorrect ID format, ID must be integer\n")
				}
				oldLen := len(data.Expenses)
				delete(&data, id)
				save(data)
				if oldLen == len(data.Expenses) {
					fmt.Printf("No expense with ID %d\n", id)
				} else {
					fmt.Printf("Expense deleted successfully!\n")
				}
			}
		}
	}

}
