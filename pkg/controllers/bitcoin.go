package controllers

import (
	"bitcoin-service/pkg/config"
	"bitcoin-service/pkg/models"
	"bitcoin-service/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
)

func Subscribe(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var storage models.EmailHandler = &models.EmailJsonStorage{PathFile: "storage/emails.json"}
	addr, err := mail.ParseAddress(r.Form.Get("email"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Incorrect email"))
		return
	}
	res, err := storage.AddEmail(addr.Address)
	if errors.Is(err, models.ErrDuplicateEmail) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	} else {
		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(res)
		w.Write(res)
	}
}

func GetRate(w http.ResponseWriter, r *http.Request) {
	var converter utils.BitcoinReader = &utils.BitcoinConverterCoingate{Domain: config.BitcoinCoingateDomain}
	rate, err := converter.ExchangeRate("UAH")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprint(rate)))

}

func SendEmails(w http.ResponseWriter, r *http.Request) {
	var converter utils.BitcoinReader = &utils.BitcoinConverterCoingate{Domain: config.BitcoinCoingateDomain}
	rate, err := converter.ExchangeRate("UAH")
	var notifier utils.EmailNotifier = &utils.EmailBTCtoUAHNotifier{
		Host:     config.EmailHost,
		Port:     config.EmailPort,
		From:     config.EmailName,
		Password: config.EmailPass,
		Rate:     rate,
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	var storage models.EmailHandler = &models.EmailJsonStorage{PathFile: config.EmailsStoragePath}
	emails, err := storage.GetALlEmails()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	go notifier.SendEmails(emails)
	w.WriteHeader(http.StatusOK)
}
