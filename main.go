package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

type config struct {
	google struct {
		apiKey string `yaml:"api_key"`
	} `yaml:"google"`
}

func (cfg *config) home(w http.ResponseWriter, r *http.Request) {
	res, err := translateText("pl", "Hello from the other side")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, res)
}

func translateText(targetLanguage, text string) (string, error) {
	ctx := context.Background()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", fmt.Errorf("language.Parse: %v", err)
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", fmt.Errorf("Translate: %v", err)
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("Translate returned empty response to text: %s", text)
	}
	return resp[0].Text, nil
}

func parseConfig() (*config, error) {
	f, err := os.Open("./config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	conf := config{}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func main() {
	config, err := parseConfig()
	mux := http.NewServeMux()
	mux.HandleFunc("/", config.home)
	err = http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
