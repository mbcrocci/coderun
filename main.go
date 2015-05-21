package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/codegangsta/cli"
)

type CodeFile struct {
	filename  string
	extension string
	name      string
	ftype     string
}

func newCodeFile(filename string) CodeFile {
	var f CodeFile

	f.filename = filename
	f.extension = filepath.Ext(filename)
	f.name = filename[0 : len(filename)-len(f.extension)]

	if f.extension != "" {
		f.ftype = f.extension[1:]
	}
}

func (f *CodeFile) show() {
	fmt.Printf("[Filename]: %s\n[Name]: %s\n[Extension]: %s\n[Type]: %s\n",
		f.filename, f.name, f.extension, f.ftype)
}

func (f *CodeFile) Compile() error {
	// (TODO): Support different filetype
	// (TODO): Check operating system

	cmd := exec.Command("gcc", "-o", f.name, f.filename)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (f *CodeFile) Run() error {
	var cmd *exec.Cmd

	switch f.ftype {
	case "":
		cmd = exec.Command("./", f.name)
	case "c":
		f.Compile()
		cmd = exec.Command("./" + f.name)
	case "go":
		cmd = exec.Command("go", "run", f.filename)
	case "py":
		// (TODO): different python versions
		cmd = exec.Command("python", f.filename)
	case "lua":
		cmd = exec.Command("lua", f.filename)
	case "ruby":
		cmd = exec.Command("ruby", f.filename)

	default:
		return errors.New("Not a valid filetype!\n[Supported] bin, C",
			+"Go, Python, Lua, Ruby")
	}
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func compile(c *cli.Context) {
	file := newCodeFile(c.Args().First())

	if err := file.Compile(); err != nil {
		fmt.Errorf("%v\n", err)
	}
}

func run(c *cli.Context) {
	file := newCodeFile(c.Args().First())

	if err := file.Run(); err != nil {
		fmt.Errorf("%v\n", err)
	}
}

func main() {
	app := cli.NewApp()
	app.Usage = "Compile and/or run your code files"
	app.Commands = []cli.Command{
		{
			Name:      "compile",
			ShortName: "c",
			Usage:     "compiles a file",
			Action, compile,
		},
		{
			Name:   "run",
			Usage:  "Run a file/script, it compiles it if needed",
			Action: run,
		},
	}
	app.RunAndExitOnError()
}
