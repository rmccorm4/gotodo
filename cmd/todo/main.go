package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/rmccorm4/gotodo/todo"
)

const dbPath = "db.pb"

func add(text string) error {
	task := &todo.Task{
		Text: text,
		Done: false,
	}
	bs, err := proto.Marshal(task)
	if err != nil {
		return fmt.Errorf("Couldn't encode task: %v", err)
	}

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("Couldn't open file %s: %v", dbPath, err)
	}

	// Write task length to database to know how many bytes to read in list()
	if err := gob.NewEncoder(f).Encode(int64(len(bs))); err != nil {
		return fmt.Errorf("Couldn't encode message: %v", err)
	}
	_, err = f.Write(bs)
	if err != nil {
		return fmt.Errorf("Couldn't write task to file %s: %v", dbPath, err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("Couldn't close file %s: %v", dbPath, err)
	}

	return nil
}

func list() error {
	bs, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return fmt.Errorf("Couldn't read file %s: %v", dbPath, err)
	}

	for {
		if len(bs) == 0 {
			return nil
		} else if len(bs) < 4 {
			return fmt.Errorf("Bytes missing length header, only %d bytes remaining in byte slice.", len(bs))
		}
		// Get first 4 bytes containing the length of our message
		var length int64
		if err := gob.NewDecoder(bytes.NewReader(bs[:4])).Decode(&length); err != nil {
			return fmt.Errorf("Couldn't decode message length: %v", err)
		}
		// Remove length bytes from the beginning so we can read current task
		bs = bs[4:]
		var task todo.Task
		err = proto.Unmarshal(bs[:length], &task)
		if err == io.EOF {
			fmt.Println("Reached EOF.")
			return nil
		} else if err != nil {
			return fmt.Errorf("Couldn't read task: %v", err)
		}
		// Remove message bytes from the beginning so we can read next length/task
		bs = bs[length:]

		if task.Done {
			fmt.Printf("✔️")
		} else {
			fmt.Printf("❌")
		}
		fmt.Printf(" %s\n", task.Text)
	}
	return nil
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		subcommands := []string{"add", "list"}
		fmt.Println("Missing subcommand:", subcommands)
		os.Exit(1)
	}

	var err error
	switch cmd := flag.Arg(0); cmd {
	case "add":
		err = add(strings.Join(flag.Args()[1:], " "))
	case "list":
		err = list()
	default:
		err = fmt.Errorf("Unknown subcommand: %s", cmd)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
