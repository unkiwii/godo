package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	todo "github.com/ritikdhasmana/godo"
)

const usage = `
Usage of godo:
  -h, --help          Prints help information 
  -a, --add           Add a new todo (Type the task as flag arguement or in new line)
  -l, --list          List all todos
  -f, --finish        Mark todo as finished whose id matches with the id passed
  -i, --id            List the todo whose id matches with the id passed
  -ss, --set-status   Set a custom status for todo with matching id (Pass the new status as flag arguement after passing id). Example: 'godo --set-status 1 apple'
  -d, --delete        Deletes the todo with the passed id
`

func main() {
	add := flag.Bool("add", false, "Add a new todo (Type the task as flag arguement or in new line)")
	flag.BoolVar(add, "a", false, "Add a new todo (Type the task as flag arguement or in new line)")

	finish := flag.Int("finish", 0, "Mark todo as finished whose id matches with the id passed")
	flag.IntVar(finish, "f", 0, "Mark todo as finished whose id matches with the id passed")

	id := flag.Int("id", 0, "List the todo whose id matches with the id passed")
	flag.IntVar(id, "i", 0, "List the todo whose id matches with the id passed")

	setStatus := flag.Int(
		"set-status",
		0,
		"Set a custom status for todo with matching id (Pass the new status as flag arguement after passing id). Example: 'godo --set-status 1 apple'",
	)
	flag.IntVar(
		setStatus,
		"ss",
		0,
		"Set a custom status for todo with matching id (Pass the new status as flag arguement after passing id). Example: 'godo --set-status 1 apple'",
	)

	delete := flag.Int("delete", 0, "Deletes the todo with the passed id")
	flag.IntVar(delete, "d", 0, "Deletes the todo with the passed id")

	list := flag.Bool("list", false, "List all todos")
	flag.BoolVar(list, "l", false, "List all todos")

	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	todos := must2(todo.Load)

	switch {
	case *add:
		task := must2(getInput(os.Stdin, flag.Args()...))

		todos.Add(task)
		must(todos.Store)
		fmt.Println("Added!")

	case *setStatus > 0:
		status := must2(getStatus(flag.Args()...))
		must(func() error { return todos.UpdateStatus(*setStatus, status) })
		must(todos.Store)
		fmt.Println("Updated!")

	case *finish > 0:
		must(func() error { return todos.UpdateStatus(*finish, "Done") })
		must(todos.Store)
		fmt.Println("Updated!")

	case *id > 0:
		todos.PrintTodo(*id)

	case *list:
		todos.Print()

	case *delete > 0:
		must(func() error { return todos.Delete(*delete) })
		must(todos.Store)
		fmt.Println("Deleted!")

	default:
		fmt.Fprintln(os.Stdout, "Invalid command! Type `godo --help` to see all available commands.")
		os.Exit(0)
	}
}

func getInput(r io.Reader, args ...string) func() (string, error) {
	return func() (string, error) {
		if len(args) > 0 {
			return strings.Join(args, " "), nil
		}

		scanner := bufio.NewScanner(r)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			return "", err
		}

		text := scanner.Text()
		if len(text) == 0 {
			return "", errors.New("empty todo is not allowed")
		}

		return text, nil
	}
}

func getStatus(args ...string) func() (string, error) {
	return func() (string, error) {
		if len(args) > 0 {
			return strings.Join(args, " "), nil
		}

		return "", errors.New("empty status is not allowed, pass status along with task id")
	}
}

func must(f func() error) {
	err := f()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func must2[T any](f func() (T, error)) T {
	v, err := f()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	return v
}
