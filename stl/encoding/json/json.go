// This sample program demonstrates how to decode a JSON response
// using the json package and NewDecoder function.
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

//type (
//	// gResult maps to the result document received from the search.
//	gResult struct {
//		GsearchResultClass string `json:"GsearchResultClass"`
//		UnescapedURL       string `json:"unescapedUrl"`
//		URL                string `json:"url"`
//		VisibleURL         string `json:"visibleUrl"`
//		CacheURL           string `json:"cacheUrl"`
//		Title              string `json:"title"`
//		TitleNoFormatting  string `json:"titleNoFormatting"`
//		Content            string `json:"content"`
//	}
//
//	// gResponse contains the top level document.
//	gResponse struct {
//		ResponseData struct {
//			Results []gResult `json:"results"`
//		} `json:"responseData"`
//	}
//)
//
//func main() {
//	//uri := "http://ajax.googleapis.com/ajax/services/search/web?v=1.0&rsz=8&q=golang"
//	uri := "//https://cn.bing.com/search?q=golang&PC=U316&FORM=CHROMN"
//
//	// Issue the search against Google.
//	resp, err := http.Get(uri)
//	if err != nil {
//		log.Println("ERROR:", err)
//		return
//	}
//	defer resp.Body.Close()
//
//	// Decode the JSON response into our struct type.
//	var gr gResponse
//	err = json.NewDecoder(resp.Body).Decode(&gr)
//	if err != nil {
//		log.Println("ERROR:", err)
//		return
//	}
//
//	fmt.Println(gr)
//
//	// Marshal the struct type into a pretty print
//	// version of the JSON document.
//	pretty, err := json.MarshalIndent(gr, "", "    ")
//	if err != nil {
//		log.Println("ERROR:", err)
//		return
//	}
//
//	fmt.Println(string(pretty))
//}

// Contact represents our JSON string.
type Contact struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	Contact struct {
		Home string `json:"home"`
		Cell string `json:"cell"`
	} `json:"contact"`
}

// JSON contains a sample string to unmarshal.
var JSON = `{
	"name": "Gopher",
	"title": "programmer",
	"contact": {
		"home": "415.333.3333",
		"cell": "415.555.5555"
	}
}`

func main() {
	// Unmarshal the JSON string into our variable.
	var c Contact
	err := json.Unmarshal([]byte(JSON), &c)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	fmt.Println(c)


	// Unmarshal the JSON string into our map variable.
	var c1 map[string]interface{}
	err = json.Unmarshal([]byte(JSON), &c1)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	fmt.Println("Name:", c1["name"])
	fmt.Println("Title:", c1["title"])
	fmt.Println("Contact")
	fmt.Println("H:", c1["contact"].(map[string]interface{})["home"])
	fmt.Println("C:", c1["contact"].(map[string]interface{})["cell"])

	// Create a map of key/value pairs.
	c2 := make(map[string]interface{})
	c2["name"] = "Gopher"
	c2["title"] = "programmer"
	c2["contact"] = map[string]interface{}{
		"home": "415.333.3333",
		"cell": "415.555.5555",
	}

	// Marshal the map into a JSON string.
	data, err := json.MarshalIndent(c2, "", "    ")
	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	fmt.Println(string(data))
}
