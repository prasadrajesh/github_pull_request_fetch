package main

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	//"log"
	"net/http"
	"net/smtp"

	log "github.com/sirupsen/logrus"
	gojsonq "github.com/thedevsaddam/gojsonq/v2"
)

func main() {
	temp := pull_open("open")
	os.WriteFile("temp.json", temp, 0777)

	jq := gojsonq.New().File("./temp.json")

	date, _ := jq.PluckR("created_at")
	title, _ := jq.PluckR("title")

	title_slice, _ := title.StringSlice()
	date_slice, _ := date.StringSlice()

	data_map := make(map[string]string)
	for i := 0; i < len(date_slice); i += 2 {
		data_map[date_slice[i]] = title_slice[i+1]
	}

	final_msg := ""
	for key, value := range data_map {

		date := key[:strings.IndexByte(key, 'T')]
		converted_time, err := time.Parse("2006-01-02", date)
		if err != nil {
			log.Fatal(err)
		}
		now := time.Now()
		difference := now.Sub(converted_time)
		if difference.Hours() < 193 {
			final_msg = final_msg + key + " " + value + ", "
		}
	}

	sendMail(final_msg)
}

func sendMail(body string) {
	from := "prasadrajesh_cc@live.com"
	password := "<paste Google password or app password here>"

	toEmailAddress := "<paste the email address you want to send to>"
	to := []string{toEmailAddress}

	host := "smtp-mail.outlook.com"
	port := "587"
	address := host + ":" + port

	subject := "Subject: Last 7 days PR\n"
	message := []byte(subject + body)

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		panic(err)
	}
}

func pull_open(state string) []byte {

	make_url := "https://api.github.com/repos/docker/docker.github.io/pulls?state=" + state + "&sort=created&direction=dec&page=1"
	data := SendGet(make_url)
	return data
}

func DoRequest(request *http.Request) []byte {

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("ERROR - Got error in response")
	}

	return body
}

func SendGet(endpoint string) []byte {

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Println("\t\t\t ERROR in SendGet")
	}
	return DoRequest(req)
}
