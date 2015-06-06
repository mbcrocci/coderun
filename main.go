package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type CodeFile struct {
	filename  string
	extension string
	name      string
	ftype     string
}

// newCodefile returns a new CodeFile,
// removes extension from filename,
// removes name form filename,
// and sets type from the extension
func newCodeFile(filename string) CodeFile {
	var f CodeFile

	f.filename = filename
	f.extension = filepath.Ext(filename)
	f.name = filename[0 : len(filename)-len(f.extension)]

	if f.extension != "" {
		f.ftype = f.extension[1:]
	}
	return f
}

// Show printf the CodeFile fields for debug purposes
func (f *CodeFile) show() {
	fmt.Printf("[Filename]: %s\n[Name]: %s\n[Extension]: %s\n[Type]: %s\n",
		f.filename, f.name, f.extension, f.ftype)
}

// Compile
// Supports C and Go files
func (f *CodeFile) Compile() error {
	var cmd *exec.Cmd

	switch f.ftype {
	case "c":
		cmd = exec.Command("gcc", "-o", f.name, f.filename)

	case "go":
		cmd = exec.Command("go", "build", f.filename)

	default:
		return errors.New("Not a supported filetype!")
	}

	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// Run the file/script
// Supports C, Go, Python, Lua, Ruby
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
		return errors.New("Not a valid filetype!\n[Supported] bin, C" +
			"Go, Python, Lua, Ruby")
	}
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func helpMessage() {
	fmt.Println("coderun filename [option]")
	fmt.Println("Options (optional):\n\tcompile\n\trun")
	fmt.Println("It defaults to run")
}

func main() {
	file := newCodeFile(os.Args[1])

	if len(os.Args) > 2 {
		switch os.Args[2] {
		case "compile", "c", "-c", "--compile":
			if err := file.Compile(); err != nil {
				fmt.Errorf("%v\n", err)
			}
		case "run", "r", "-r", "--run", "":
			if err := file.Run(); err != nil {
				fmt.Errorf("%v\n", err)
			}

		default:
			helpMessage()
			os.Exit(1)
		}
	} else {
		if err := file.Run(); err != nil {
			fmt.Errorf("%v\n", err)
		}
	}
}
