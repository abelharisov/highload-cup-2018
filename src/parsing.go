package main

import (
	"archive/zip"
	"encoding/json"
	// "fmt"
	"sort"
)

type byFileName []*zip.File

func (files byFileName) Len() int {
	return len(files)
}

func (files byFileName) Swap(i, j int) {
	files[i], files[j] = files[j], files[i]
}

func (files byFileName) Less(i, j int) bool {
	return files[i].Name < files[j].Name
}

// Parse files in dataPath and put to Storage
func Parse(dataPath string, storage Storage) error {
	zipReader, err := zip.OpenReader(dataPath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	sort.Sort(byFileName(zipReader.File))
	for _, f := range zipReader.File {
		reader, _ := f.Open()
		decoder := json.NewDecoder(reader)
		decoder.Token()
		decoder.Token()
		decoder.Token()

		var account Account
		accounts := make([]Account, 0, 1000)

		for decoder.More() {
			err := decoder.Decode(&account)
			if err == nil {
				accounts = append(accounts, account)
			}
		}

		reader.Close()

		storage.LoadAccounts(accounts)
	}

	return nil
}
