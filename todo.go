package godo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
)

type item struct {
	Task      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (i *item) UpdateStatus(status string) {
	i.Status = strings.ToUpper(status)
	i.UpdatedAt = time.Now()
}

func (i item) Is(status string) bool {
	return strings.ToUpper(status) == strings.ToUpper(i.Status)
}

const DateTimeFormat = "January _2 2006, 15:04"

func (i item) ToRow(index int) []*simpletable.Cell {
	task := blue(i.Task)
	status := blue(i.Status)
	if i.Is("Done") {
		task = green(fmt.Sprintf("\u2705 %s", i.Task))
		status = green(strings.ToUpper(i.Status))
	}
	updatedAt := ""
	if !i.UpdatedAt.IsZero() {
		updatedAt = i.UpdatedAt.Format(DateTimeFormat)
	}
	return []*simpletable.Cell{
		{Text: fmt.Sprintf("%d", index)},
		{Text: task},
		{Text: status},
		{Text: i.CreatedAt.Format(DateTimeFormat)},
		{Text: updatedAt},
	}
}

type Todos []item

const todoFileName = ".godo.json"

func getAbsolutePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, todoFileName), nil
}

func (t *Todos) Add(task string) {
	*t = append(*t, item{
		Task:      task,
		Status:    "TODO",
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{},
	})
}

func (t *Todos) UpdateStatus(taskId int, status string) error {
	if taskId <= 0 || taskId > len(*t) {
		return errors.New("invalid index")
	}

	(*t)[taskId-1].UpdateStatus(status)
	return nil
}

func (t *Todos) Delete(taskId int) error {
	if taskId <= 0 || taskId > len(*t) {
		return errors.New("invalid index")
	}

	*t = append((*t)[:(taskId-1)], (*t)[taskId:]...)
	return nil
}

func Load() (*Todos, error) {
	var todos Todos

	fileName, err := getAbsolutePath()
	if err != nil {
		return &todos, err
	}

	file, err := os.ReadFile(fileName)
	switch {
	case errors.Is(err, os.ErrNotExist), len(file) == 0:
		return &todos, nil
	case err != nil:
		return &todos, err
	}

	err = json.Unmarshal(file, &todos)
	if err != nil {
		return &todos, err
	}

	return &todos, nil
}

func (t *Todos) Store() error {
	fileName, err := getAbsolutePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, data, 0644)
}

func (t *Todos) header() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Todo"},
			{Align: simpletable.AlignCenter, Text: "Status"},
			{Align: simpletable.AlignCenter, Text: "Created"},
			{Align: simpletable.AlignCenter, Text: "Updated"},
		},
	}
}

func (t *Todos) Print() {
	var rows [][]*simpletable.Cell
	for idx, item := range *t {
		idx++
		rows = append(rows, item.ToRow(idx))
	}

	table := simpletable.Table{
		Header: t.header(),
		Body:   &simpletable.Body{Cells: rows},
		Footer: &simpletable.Footer{
			Cells: []*simpletable.Cell{
				{
					Align: simpletable.AlignCenter,
					Span:  5,
					Text:  red(fmt.Sprintf("You have %d pending todos.", t.CountPending())),
				},
			},
		},
	}

	table.SetStyle(simpletable.StyleUnicode)
	table.Println()
}

func (t *Todos) PrintTodo(id int) {
	var rows [][]*simpletable.Cell
	for idx, item := range *t {
		idx++
		if idx == id {
			rows = append(rows, item.ToRow(idx))
		}
	}

	table := simpletable.Table{
		Header: t.header(),
		Body:   &simpletable.Body{Cells: rows},
	}

	table.SetStyle(simpletable.StyleUnicode)
	table.Println()
}

func (t *Todos) CountPending() int {
	count := 0
	for _, item := range *t {
		if !item.Is("Done") {
			count++
		}
	}
	return count
}
