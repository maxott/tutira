package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/maxott/welisp/core"
	relalg "github.com/maxott/welisp/relalg"
	repl "github.com/openengineer/go-repl"
)

type MyHandler struct {
	r *repl.Repl
}

var helpMessage = `help              display this message
add <int> <int>   add two numbers
quit              quit this program`

func main() {

	var prog string
	var query string

	flag.StringVar(&prog, "e", "", "program to execute immediately")
	flag.StringVar(&query, "q", "", "run query")

	flag.Parse()

	if query != "" {
		if res, err := relalg.EvalQuery(query, ""); err == nil {
			fmt.Println(res.Print(false))
		} else {
			log.Fatal(err)
		}
		return
	}

	if prog != "" {
		if res, err := core.Eval(prog); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(res.Print(false))
		}
	} else {
		fmt.Println("Welcome, type \"help\" for more info")
		h := &MyHandler{}
		h.r = repl.NewRepl(h)

		// start the terminal loop
		if err := h.r.Loop(); err != nil {
			log.Fatal(err)
		}
	}

	// list := core.parse(prog, t)
	// scope := NewScope()
	// return list.Eval(scope)

}

func (h *MyHandler) Prompt() string {
	return "> "
}

func (h *MyHandler) Tab(buffer string) string {
	return "" // do nothing
}

func (h *MyHandler) Eval(line string) string {
	fields := strings.Fields(line)

	if len(fields) == 0 {
		return ""
	} else {
		cmd := fields[0]
		//args :=fields[1:]

		switch cmd {
		case "help":
			return helpMessage
		// case "add":
		// 	if len(args) != 2 {
		// 		return "\"add\" expects 2 args"
		// 	} else {
		// 		return add(args[0], args[1])
		// 	}
		case "quit":
			h.r.Quit()
			return ""
		default:
			return fmt.Sprintf("unrecognized command \"%s\"", cmd)
		}
	}
}
