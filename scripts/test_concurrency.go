package scripts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func UpdatePost(postId int, p UpdatePostPayload, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("http://localhost:8080/v1/posts/%d", postId)

	b, _ := json.Marshal(p)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(b))

	if err != nil {
		fmt.Println("Error creating request: ", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending request: ", err)
		return
	}

	defer resp.Body.Close()

	fmt.Println("Response Status: ", resp.Status)
}

func main() {
	var wg sync.WaitGroup

	postId := 13

	wg.Add(2)
	content := "updated content"
	title := "updated title"

	go UpdatePost(postId, UpdatePostPayload{Title: &title, Content: &content}, &wg)
	go UpdatePost(postId, UpdatePostPayload{Title: &title, Content: &content}, &wg)

	wg.Wait()
}
