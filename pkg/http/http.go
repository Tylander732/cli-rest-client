package http

import (
	"io"
	"log"
	"net/http"
)

func GetRequest(requestBody string, url string) string {
    client := &http.Client{}

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Println("Error creating request:", err)
        return ""
    }


    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error sending request:", err)
        return ""
    }

    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error reading response:", err)
        return ""
    }

    log.Println("Response from ", url, ":", string(body))

    return string(body)
}
