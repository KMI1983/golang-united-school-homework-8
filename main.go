package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Arguments map[string]string

type Item struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Perform(args Arguments, writer io.Writer) error {

	if args["operation"] == "" {
		return fmt.Errorf("-operation flag has to be specified")
	}
	if args["fileName"] == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}

	operation := args["operation"]
	switch operation {
	case "add":
		return Add(args, writer)
	case "remove":
		return Remove(args, writer)
	case "findById":
		return FindById(args, writer)
	case "list":
		return List(args, writer)
	default:
		return fmt.Errorf("Operation %v not allowed!", operation)
	}
}

func Add(args Arguments, writer io.Writer) error {

	if args["item"] == "" {
		return fmt.Errorf("-item flag has to be specified")
	}
	item := Item{}
	json.Unmarshal([]byte(args["item"]), &item)
	fileName := args["fileName"]

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer file.Close()
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	var data []Item
	json.Unmarshal(content, &data)
	if len(data) > 0 {
		for _, it := range data {
			fmt.Println(it)
			if it.Id == item.Id {

				fmt.Fprintf(writer, "Item with id %v already exists", item.Id)
				return nil
			}
		}
	}

	data = append(data, item)
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}
	file.Write(res)
	return err
}

func Remove(args Arguments, writer io.Writer) error {

	if args["id"] == "" {
		return fmt.Errorf("-id flag has to be specified")
	}
	id, _ := args["id"]
	fileName := args["fileName"]
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer file.Close()
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	var data []Item
	json.Unmarshal(content, &data)
	found := false
	var data2 []Item
	if len(data) > 0 {
		for _, it := range data {
			if it.Id != id {
				data2 = append(data2, it)
			} else {
				found = true
			}
		}
	}
	if !found {
		fmt.Fprintf(writer, "Item with id %v not found", id)
		return nil
	}

	res, err := json.Marshal(data2)
	if err != nil {
		return err
	}
	file.Truncate(0)
	file.Seek(0, 0)
	file.Write(res)
	return err
}

func FindById(args Arguments, writer io.Writer) error {

	if args["id"] == "" {
		return fmt.Errorf("-id flag has to be specified")
	}
	id, _ := args["id"]
	fileName := args["fileName"]
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer file.Close()
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	var data []Item
	json.Unmarshal(content, &data)
	found := Item{}
	if len(data) > 0 {
		for _, it := range data {
			if it.Id == id {
				found = it
			}
		}
	}
	if found.Id != "" {
		res, _ := json.Marshal(found)
		writer.Write([]byte(res))
	} else {
		writer.Write([]byte(""))
	}
	return nil
}

func List(args Arguments, writer io.Writer) error {

	fileName := args["fileName"]
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer file.Close()
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	var data []Item
	json.Unmarshal(content, &data)
	if len(data) > 0 {
		res, _ := json.Marshal(data)
		writer.Write([]byte(res))
	} else {
		writer.Write([]byte(""))
	}
	return nil
}

func parseArgs() Arguments {

	operation := flag.String("operation", "", "an operation")
	item := flag.String("item", "", "an item")
	fileName := flag.String("fileName", "", "a fileName")
	id := flag.String("id", "", "an id")
	flag.Parse()

	r := strings.NewReplacer(
		"«", "\"",
		"»", "\"",
	)

	s := r.Replace(*item)
	res := Item{}
	json.Unmarshal([]byte(s), &res)

	args := Arguments{

		"id":        *id,
		"operation": strings.Trim(r.Replace(*operation), "\""),
		"item":      r.Replace(*item),
		"fileName":  strings.Trim(r.Replace(*fileName), "\""),
	}
	return args

}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
