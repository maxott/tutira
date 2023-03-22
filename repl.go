package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/maxott/tutiro/core"
	relalg "github.com/maxott/tutiro/relalg"
	repl "github.com/openengineer/go-repl"
)

var helpMessage = `
help              display this message
quit              quit this program
`

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
		repl.DEBUG = "/tmp/repl.log"
		fmt.Println("Welcome, type \"?help\" for more info")
		r := NewRepl()
		if err := r.Loop(); err != nil {
			log.Fatal(err)
		}
	}
}

type Handler struct {
	stream     *core.ByteStream
	reader     core.ReaderI
	machine    core.MachineI
	level      int
	openString bool
	r          *repl.Repl
}

func NewRepl() *repl.Repl {
	stream := core.NewByteStream([]byte{})
	reader := core.Reader(stream)
	machine := core.NewMachine(reader)

	h := &Handler{
		stream:  stream,
		reader:  reader,
		machine: machine,
	}
	h.r = repl.NewRepl(h)
	return h.r
}
func (h *Handler) Prompt() string {
	if h.level > 0 {
		return fmt.Sprintf("%d> ", h.level)
	} else if h.openString {
		return "=> "
	}
	return " > "
}

func (h *Handler) Tab(buffer string) string {
	return "" // do nothing
}

func (h *Handler) Eval(line string) string {
	if len(line) == 0 {
		return ""
	}
	if line[0] == '?' {
		return h.commands(line[1:])
	}

	c := context.Background()
	h.stream.AppendData([]byte(line)).Mark()
	_, err := h.reader.Next(c, true)
	h.stream.Reset(c)
	h.openString = false
	h.level = 0
	if err == nil {
		//fmt.Printf("-- '%s' \n", h.stream.ToString())
		if term, err2 := h.machine.Eval(c); err2 == nil {
			return term.Print(false)
		} else {
			return fmt.Sprintf("Error: %v", err2)
		}
	} else {
		// not finished form
		h.stream.AppendData([]byte("\n"))
		//fmt.Printf("err-- '%v' \n", err)
		if _, ok := err.(*core.OpenStringReaderError); ok {
			h.openString = true
		} else if uerr, ok := err.(*core.UnbalancedReaderError); ok {
			h.level = uerr.Level
		}
		return ""
	}

}

func (h *Handler) commands(line string) string {
	//fmt.Printf("cmd--- '%s'\n", line)
	fields := strings.Fields(line)

	if len(fields) == 0 {
		return ""
	} else {
		cmd := fields[0]
		//args :=fields[1:]

		switch cmd {
		case "help":
			return helpMessage
		case "quit":
			h.r.Quit()
			return ""
		default:
			return fmt.Sprintf("unrecognized command \"%s\"", cmd)
		}
	}
}
