package main

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
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
func Parse(dataPath string, optionsPath string, storage Storage, onlyOne bool) error {
	optionsFile, err := os.Open(optionsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer optionsFile.Close()

	optionsReader := bufio.NewReader(optionsFile)
	timeNow, err := optionsReader.ReadString('\n')
	if err != nil {
		return err
	}
	timeNow = timeNow[0:len(timeNow) - 1]
	now, err := strconv.Atoi(timeNow)
	if err != nil {
		return err
	}

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
				birthYear := time.Unix(int64(*account.Birth), 0).Year()
				account.BirthYear = &birthYear
				joinedYear := time.Unix(int64(*account.Joined), 0).Year()
				account.JoinedYear = &joinedYear
				if account.Likes != nil {
					ids := make([]int, 0, len(*account.Likes))
					for _, like := range *account.Likes {
						ids = append(ids, like.Id)
					}
					account.LikeIds = &ids
				}
				if account.Premium == nil {
					account.PremiumStatus = PremiumNull
				} else {
					if account.Premium.Start < int64(now) && int64(now) < account.Premium.Finish {
						account.PremiumStatus = PremiumActive
					} else {
						account.PremiumStatus = PremiumNotActive
					}
				}
				if account.Phone != nil {
					r := regexp.MustCompile(`\((?P<Code>\d.+)\)`)
					matches := r.FindStringSubmatch(*account.Phone)
					codeStr := matches[1]
					if code, err := strconv.Atoi(codeStr); err == nil {
						account.PhoneCode = code
					}
					// log.Println(*account.Phone, codeStr, account.PhoneCode)
				}
				account.StatusId = ParseStatus(account.Status)
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

func ParseStatus(status string) int {
	char := []rune(status)[0]
	switch char {
	case 'с':
		return StatusFree
	case 'з':
		return StatusBusy
	default:
		return StatusWtf
	}
}
