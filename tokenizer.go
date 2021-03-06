package tokenizer

import (
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/registry"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	Name = "sego"
)

var SegoTokenizerHost = "localhost:3000"

type SegoTokenizer struct{}

func NewSegoTokenizer() (*SegoTokenizer, error) {
	return &SegoTokenizer{}, nil
}

func (this *SegoTokenizer) Tokenize(b []byte) (stream analysis.TokenStream) {
	stream = make(analysis.TokenStream, 0)
	var form = url.Values{}
	form.Add("data", string(b))
	var url = fmt.Sprintf("http://%s/api/tokenizer/", SegoTokenizerHost)

	var req, _ = http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("http.DefaultClient.Do() failed (%s)\n", err)
		return
	}
	defer rsp.Body.Close()
	if int(rsp.StatusCode/100) != 2 {
		log.Printf("tokenizer failed\n")
		return
	}

	decoder := json.NewDecoder(rsp.Body)
	if err = decoder.Decode(&stream); err != nil {
		log.Printf("json.NewDecoder().Decode() failed (%s)", err)
		return
	}

	return
}

func SegoTokenizerConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.Tokenizer, error) {
	return NewSegoTokenizer()
}

func init() {
	registry.RegisterTokenizer(Name, SegoTokenizerConstructor)
}
