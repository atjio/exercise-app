package handlers

import (
	"log"
	"net/http"
	"time"

	"SimpleApp/global"
)

func RegisterService(url string) error {
	log.Println(global.GetHealthcheckUrl())
	for {
		resp, err := http.Get(global.GetHealthcheckUrl())
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(1 * time.Second)
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Client-Port", global.PORT)

	_, err = client.Do(req)
	if err != nil {
		log.Fatal("Failed to register service: ", err)
	}

	return err

}
