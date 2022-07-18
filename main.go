package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ApiGatewayLambdaResponse struct {
	StatusCode int
	Body       interface{}
	Headers    map[string]string
}

func getEnvWithDefault(key string, defaultValue string) string {
	if str := os.Getenv(key); str != "" {
		return str
	}
	return defaultValue
}
func main() {
	url := getEnvWithDefault("LAMBDA_TARGET_URL", "http://localhost:9090")
	port := getEnvWithDefault("PORT", "8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// lambda rie setting
		const path = "2015-03-31/functions/function/invocations"
		var resp *http.Response
		var err error
		if r.Method == http.MethodGet {
			resp, err = http.Post(fmt.Sprintf("%s/%s", url, path), "application/json", bytes.NewBufferString("{}"))
			if err != nil {
				log.Printf("failed to get response from lambda %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Printf("failed to read request body %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Printf("%s", string(body))
			resp, err = http.Post(fmt.Sprintf("%s/%s", url, path), "application/json", bytes.NewBufferString(string(body)))
			if err != nil {
				log.Printf("failed to get response from lambda %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		defer resp.Body.Close()
		byteArray, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("failed to read response from lambda %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s := string(byteArray)
		m := &ApiGatewayLambdaResponse{}
		json.Unmarshal([]byte(s), &m)
		for key, value := range m.Headers {
			w.Header().Set(key, value)
		}
		w.WriteHeader(m.StatusCode)
		obj, err := json.Marshal(m.Body)
		if err != nil {
			log.Printf("failed to read response from lambda %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println(string(obj))
		fmt.Fprintln(w, string(obj))
	})
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Listen and Serve failse: %v", err)
	}
}
