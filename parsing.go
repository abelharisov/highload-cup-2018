package main

import (
	"archive/zip"
	"encoding/json"
	"sort"
	"time"
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
func Parse(dataPath string, storage Storage, onlyOne bool) error {
	zipReader, err := zip.OpenReader(dataPath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	sort.Sort(byFileName(zipReader.File))
	for _, f := range zipReader.File {
		// log.Println("Start parsing file ", f.Name)
		reader, _ := f.Open()
		decoder := json.NewDecoder(reader)
		decoder.Token()
		decoder.Token()
		decoder.Token()

		accounts := make([]Account, 0, 1000)

		for decoder.More() {
			var account Account
			err := decoder.Decode(&account)
			if err == nil {
				year := time.Unix(int64(*account.Birth), 0).Year()
				account.Year = &year
				if account.Likes != nil {
					ids := make([]int, 0, len(*account.Likes))
					for _, like := range *account.Likes {
						ids = append(ids, like.Id)
					}
					account.LikeIds = &ids
				}
				accounts = append(accounts, account)
			}
		}

		reader.Close()

		// log.Println("Finish parsing file ", f.Name)
		storage.LoadAccounts(accounts)
		// log.Println("End loading parsed ", f.Name)

		if onlyOne {
			break
		}
	}

	return nil
}
