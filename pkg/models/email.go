package models

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

var ErrDuplicateEmail = errors.New("This email already exist")

type EmailHandler interface {
	AddEmail(email string) (string, error)
	GetALlEmails() ([]string, error)
}

type EmailJsonStorage struct {
	PathFile string
}

func (storage *EmailJsonStorage) GetALlEmails() ([]string, error) {
	jsonFile, err := os.Open(storage.PathFile)
	if err != nil {
		log.Println(err)
		os.Create(storage.PathFile)
		jsonFile, err := os.Open(storage.PathFile)
		defer jsonFile.Close()
		if err != nil {
			return nil, err
		}
		log.Println("Created new file storage for emails")
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result []string
	json.Unmarshal([]byte(byteValue), &result)
	return result, nil
}

func (storage *EmailJsonStorage) AddEmail(email string) (string, error) {
	emails, err := storage.GetALlEmails()
	if err != nil {
		return "", err
	}
	i := sort.SearchStrings(emails, email)
	if len(emails) > 0 && len(emails) > i && emails[i] == email {
		return "", ErrDuplicateEmail
	}
	emails = append(emails, "")
	if len(emails) > 1 {
		copy(emails[i+1:], emails[i:])
	}
	emails[i] = email
	content, err := json.Marshal(emails)
	err = ioutil.WriteFile(storage.PathFile, content, 0644)
	if err != nil {
		return "", err
	}
	return email, nil
}
