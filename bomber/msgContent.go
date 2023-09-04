package bomber

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Article struct {
	Source      Source `json:"source"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	URLToImage  string `json:"urlToImage"`
	PublishedAt string `json:"publishedAt"`
	Content     string `json:"content"`
}

type APIResponse struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

func GetMsgContent(apiKey string) (APIResponse, error) {
	curDate := time.Now()
	prevMonth := curDate.AddDate(0, -1, 1)
	prevMonthDate := prevMonth.Format("2006-01-02")
	var apiUrl string
	if len(apiKey) == 0 {
		apiUrl = fmt.Sprintf("https://newsapi.org/v2/everything?q=football&from=%s&sortBy=publishedAt&apiKey=300e5445bd884438bafd4685f9a536e5", prevMonthDate)
	} else {
		apiUrl = fmt.Sprintf("https://newsapi.org/v2/everything?q=football&from=%s&sortBy=publishedAt&apiKey=%s", prevMonthDate, apiKey)
	}
	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return APIResponse{}, err
	}
	bodyByte, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return APIResponse{}, err
	}
	var body APIResponse
	err = json.Unmarshal(bodyByte, &body)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return APIResponse{}, err
	}
	return body, nil
}
