package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func main() {
	const numRuns = 20

	for run := 0; run < numRuns; run++ {
		var wg sync.WaitGroup
		const numMessages = 1000

		for i := 0; i < numMessages; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				err := getRequest()
				if err != nil {
					fmt.Println("-------------Error sending request:", err)
				}
			}(i)
		}
		wg.Wait()
		fmt.Printf("Run %d completed\n", run+1)
	}
	fmt.Println("All runs completed")
}

func callPostRequest(index int) {
	postData := generatePostData()

	// Convert the struct to JSON
	jsonData, err := json.Marshal(postData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Send the POST request
	err = postRequest("http://192.168.2.45:8090/producer", jsonData)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	fmt.Printf("Message %d sent\n", index)
}

func generatePostData() map[string]string {
	rand.Seed(time.Now().UnixNano())
	ordersn := fmt.Sprintf("22%02d%02dQSK8S7BX", rand.Intn(100), rand.Intn(100))
	data := map[string]string{"key": "ordersn", "value": ordersn}
	return data
}

func postRequest(url string, jsonData []byte) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}
	return nil
}

func getRequest() error {
	url := "http://192.168.38.128:3000/push_orders?key=ha&value=123"

	// Make GET request
	_, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	// defer resp.Body.Close()

	// // Read response body
	// _, err = ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println("Error reading response body:", err)
	// 	return err
	// }
	return nil
}
