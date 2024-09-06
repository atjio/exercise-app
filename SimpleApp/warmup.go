package main

import (
	"net/http"
	"io"
	"log"
	"time"
)

func registerService(url string) error {
	log.Println("http://localhost" + port + "/healthcheck")
	for {
        resp, err := http.Get("http://localhost" + port + "/healthcheck")
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
	req.Header.Set("X-Client-Port", port)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to register service: ", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Println("TEST")
	log.Println(string(body))


	return err
}