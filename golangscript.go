package main

import (
	"bytes"
	"context"
	"encoding/json"
	"esConnection"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
)

func CreateDefaultTemplate(s string, es *elasticsearch.Client) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"index_patterns": []string{s},
		"settings": map[string]interface{}{
			"number_of_shards": 2,
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	str := buf.String()
	fmt.Println(str)
	res, err := es.Indices.PutTemplate(
		s,
		strings.NewReader(str),
	)
	fmt.Println(res, err)
}
func CreateTemplate(s string, es *elasticsearch.Client) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"index_patterns": []string{s},
		"settings": map[string]interface{}{
			"number_of_shards": 5,
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	str := buf.String()
	fmt.Println(str)
	res, err := es.Indices.PutTemplate(
		s,
		strings.NewReader(str),
	)
	fmt.Println(res, err)
}
func IndexSizeFilter(s string, t int, d int) string {

	if t >= d {
		return s
	} else {
		return ""
	}
}
func main() {
	var r []map[string]interface{}
	log.SetFlags(0)
	var es *elasticsearch.Client

	es = esConnection.MakeConnection()

	req := esapi.CatIndicesRequest{

		Bytes:  "b",
		Format: "json",
	}
	//performing request with client
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response : %s", err)
	}
	//fmt.Println(res)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("oops!you have something missing")
	}

	var insize []int
	var s []string
	var s1 []string
	json.Unmarshal([]byte(body), &r)
	for v, r := range r {
		s = append(s, r["index"].(string))
		rf := regexp.MustCompile(`-\d{4}.\d{2}.\d{2}`)
		rk := rf.ReplaceAllString(s[v], "${1}")
		s1 = append(s1, rk)
		converter := r["store.size"].(string)
		h, _ := strconv.Atoi(converter)
		insize = append(insize, h)
	}
	for i := 0; i < len(s1); i++ {
		fmt.Println(s1[i])
	}
	check := make(map[string]int)
	s2 := make([]string, 0)
	for _, val := range s1 {
		check[val] = 1
	}

	for letter, _ := range check {
		s2 = append(s2, letter)
	}
	fmt.Println("+++++++++++")
	for i := 0; i < len(s2); i++ {
		fmt.Println(s2[i])
	}
	fmt.Println("+++++++++++")

	for i := 0; i < len(s2); i++ {
		res, err := es.Indices.GetTemplate(
			es.Indices.GetTemplate.WithName(s2[i]),
			es.Indices.GetTemplate.WithFilterPath("*.version"),
		)
		if err != nil {
			log.Fatalf("Error getting response : %s", err)
		}
		fmt.Println(res.StatusCode)
		if res.StatusCode == 404 {

			CreateDefaultTemplate(s2[i], es)

		}
		if res.StatusCode == 200 {

			str := IndexSizeFilter(s2[i], insize[i], 583)
			if str == s2[i] {
				CreateTemplate(str, es)
			} else {
				str = s2[i]
				CreateDefaultTemplate(str, es)
			}
		}
	}

}
