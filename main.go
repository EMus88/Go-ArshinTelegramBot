package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func main() {
	botToken := "2122958020:AAF9C9hC5nmCrfuqF6GWMgVD4Pq-6uAPa3s"
	botAPI := "https://api.telegram.org/bot"
	botUrl := botAPI + botToken
	arshinAPIUrl := "https://fgis.gost.ru/fundmetrology/eapi/vri?year=2021&rows=100&search=*"

	offset := 0
	for {

		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("Smth went wrong: ", err.Error())
		}
		for _, update := range updates {

			arshinSearchInputText := url.QueryEscape(update.Message.Text)
			searchResults, err := getInfoFromArshin(arshinAPIUrl, arshinSearchInputText)
			if err != nil {
				log.Println("Arshin error: ", err.Error())
			}
			buffer := bytes.Buffer{}
			for _, searchResult := range searchResults {
				buffer.WriteString("Организация поверитель " + searchResult.Organization)
				buffer.WriteString("Тип СИ " + searchResult.TypeOfDevice)
				buffer.WriteString("Заводсткой номер " + searchResult.DeviceNumber)
				buffer.WriteString("Действительно до " + searchResult.ValidDate)
			}

			err = respond(botUrl, buffer, update)
			if err != nil {
				log.Println("Smth went wrong: ", err.Error())
			}
			offset = update.UpdateId + 1
		}
		fmt.Println(updates)
	}

}

func getUpdates(botUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

func getInfoFromArshin(arshinApiUrl string, searchText string) ([]Device, error) {
	resp, err := http.Get(arshinApiUrl + searchText + "*&rows=100")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponseDevices
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	fmt.Println(restResponse.SearchResult.Devices)
	return restResponse.SearchResult.Devices, nil
}

func respond(botUrl string, bufDevices bytes.Buffer, update Update) error {
	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	botMessage.Text = bufDevices.String()
	fmt.Println("---------------------")
	fmt.Println(botMessage.Text)
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	fmt.Println("отправлено")
	return nil
}
