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

func (s storage) save() error {
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

func (s *storage) list() {
	for _, exp := range s.Expenses {
		fmt.Printf("%-4d %-12s %-5d %s\n", exp.ID,
			exp.Date.Format("02-01-2006"), exp.Amount, exp.Description)
	}
}

func (s *storage) add() {
	var exp expense
	if len(os.Args) != 4 {
		fmt.Printf("Error: invalid number of arguments, use add \"description\" amount\n")
	} else {
		amount, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Printf("Error: incorrect amount format, amount must be integer\n")
		}

		description := os.Args[2]
		exp.ID = s.Nextid
		exp.Date = time.Now()
		exp.Amount = amount
		exp.Description = description

		s.Expenses = append(s.Expenses, exp)
		fmt.Printf("Expense added successfully (ID: %d)\n", s.Nextid)
		s.Nextid++
		s.save()
	}
}

func (s *storage) summary() {
	var sum int
	if len(os.Args) == 2 {
		for _, exp := range s.Expenses {
			sum += exp.Amount
		}
		fmt.Printf("Total expenses: %d\n", sum)

	} else if len(os.Args) == 3 {
		var month string
		var monthNumber int
		months := []string{"January", "February", "March", "April",
			"May", "June", "July", "August", "September",
			"October", "November", "December"}
		monthIndx, err := strconv.Atoi(os.Args[2])

		if err == nil && monthIndx >= 1 && monthIndx <= 12 {
			month = months[monthIndx-1]
			monthNumber = monthIndx
		} else {
			for i, monthName := range months {
				if monthName == os.Args[2] {
					month = monthName
					monthNumber = i + 1
					break
				}
			}
		}

		if month != "" {
			for _, exp := range s.Expenses {
				if int(exp.Date.Month()) == monthNumber {
					sum += exp.Amount
				}
			}
			fmt.Printf("Total expenses for %s: %d\n", month, sum)
		} else {
			fmt.Println(`Error: invalid month name or number, use "summary [number of month]" or "summary [name of month]"`)
		}
	}
}

func (s *storage) delete() {
	newExpenses := make([]expense, 0)
	if len(os.Args) != 3 {
		fmt.Printf("Error: invalid number of arguments, use delete id\n")
	} else {
		id, err := strconv.Atoi(os.Args[2])

		if err != nil {
			fmt.Printf("Error: incorrect ID format, ID must be integer\n")
		}
		oldLen := len(s.Expenses)
		for _, exp := range s.Expenses {
			if exp.ID != id {
				newExpenses = append(newExpenses, exp)
			}
		}
		s.Expenses = newExpenses
		if oldLen == len(s.Expenses) {
			fmt.Printf("No expense with ID %d\n", id)
		} else {
			s.save()
			fmt.Printf("Expense deleted successfully\n")
		}
	}
}

func main() {
	data, err := load()
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "list":
			data.list()
		case "add":
			data.add()
		case "summary":
			data.summary()
		case "delete":
			data.delete()
		}
	}
}
