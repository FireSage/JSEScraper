package main

import (
	"encoding/json"
	"log"
	"os"
)

func write_json_file(data string, file_name string) {

	file, err := os.OpenFile(file_name+".json", os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		log.Print(err)
		return
		// log.Fatal(err)
	}
	if _, err := file.Write([]byte(data)); err != nil {
		file.Close()
		log.Print(err)
		return
		// log.Fatal(err)
	}

	if err := file.Close(); err != nil {
		log.Print()
		return
		// log.Fatal()
	}

}
func read_stocks_from_json_file(file_name string, stocklist *[]Stock) {
	data, file_err := os.ReadFile(file_name)
	if file_err != nil {
		log.Print(file_err)
		return
	}
	json_err := json.Unmarshal(data, &stocklist)
	if json_err != nil {
		log.Print(json_err)
	}
}
