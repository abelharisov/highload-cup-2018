package main

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

func EnrichAccount(account *Account, now int) {
	if account.Birth != nil {
		birthYear := time.Unix(int64(*account.Birth), 0).Year()
		account.BirthYear = &birthYear
	}

	if account.Joined != nil {
		joinedYear := time.Unix(int64(*account.Joined), 0).Year()
		account.JoinedYear = &joinedYear
	}

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
	}

	account.StatusId = ParseStatus(account.Status)
}

// Parse files in dataPath and put to Storage
func Parse(dataPath string, optionsPath string, storage Storage, onlyOne bool) error {
	optionsFile, err := os.Open(optionsPath)
	if err != nil {
		return err
	}
	defer optionsFile.Close()

	optionsReader := bufio.NewReader(optionsFile)
	timeNow, err := optionsReader.ReadString('\n')
	if err != nil {
		return err
	}
	timeNow = timeNow[0 : len(timeNow)-1]
	now, err := strconv.Atoi(timeNow)
	if err != nil {
		return err
	}
	storage.SetNow(now)

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

		accounts := make([]Account, 0, 1000)

		for decoder.More() {
			var account Account
			if err := decoder.Decode(&account); err == nil {
				EnrichAccount(&account, now)
				accounts = append(accounts, account)
			}
		}

		reader.Close()

		storage.LoadAccounts(accounts)

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
	case 'в':
		return StatusWtf
	default:
		return 0
	}
}

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
