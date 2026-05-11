// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// )

// var queries = []string{
// 	"love", "like", "dislike", "yes", "no", "ok", "okay",
// 	"thanks", "thank you", "sorry", "please",
// 	"lol", "haha", "funny", "laugh", "lmao", "rofl", "crying laughing",
// 	"sad", "cry", "tears", "depressed", "disappointed",
// 	"heart", "kiss", "hug", "cute", "blush", "romantic",
// 	"fire", "hot", "cool", "awesome", "epic", "wow", "amazing",
// 	"surprised", "shock", "omg", "wtf", "no way",
// 	"angry", "mad", "rage", "annoyed",
// 	"thinking", "confused", "shrug", "idk",
// 	"clap", "applause", "good job", "win", "congrats",
// 	"facepalm", "mic drop", "this is fine",
// }

// const apiKey = "lYRDYLNPrX7sCQuvchqVWi7KrX5YalCj"
// const limit = 3

// type GiphyResponse struct {
// 	Data []struct {
// 		ID     string `json:"id"`
// 		Title  string `json:"title"`
// 		URL    string `json:"url"`
// 		Images struct {
// 			Original struct {
// 				URL string `json:"url"`
// 			} `json:"original"`
// 		} `json:"images"`
// 	} `json:"data"`
// }

// type Gif struct {
// 	Query string `json:"query"`
// 	ID    string `json:"id"`
// 	Title string `json:"title"`
// 	URL   string `json:"url"`
// 	Image string `json:"image"`
// }

// func fetch(query string) ([]Gif, error) {
// 	url := fmt.Sprintf(
// 		"https://api.giphy.com/v1/gifs/search?api_key=%s&q=%s&limit=%d",
// 		apiKey, query, limit,
// 	)

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	body, _ := io.ReadAll(resp.Body)

// 	var result GiphyResponse
// 	err = json.Unmarshal(body, &result)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var gifs []Gif

// 	for _, g := range result.Data {
// 		gifs = append(gifs, Gif{
// 			Query: query,
// 			ID:    g.ID,
// 			Title: g.Title,
// 			URL:   g.URL,
// 			Image: g.Images.Original.URL,
// 		})
// 	}

// 	return gifs, nil
// }

// func main() {
// 	var all []Gif

// 	for _, q := range queries {
// 		fmt.Println("Fetching:", q)

// 		gifs, err := fetch(q)
// 		if err != nil {
// 			fmt.Println("error:", err)
// 			continue
// 		}

// 		all = append(all, gifs...)
// 	}

// 	file, err := os.Create("gifs.json")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()

// 	enc := json.NewEncoder(file)
// 	enc.SetIndent("", "  ")
// 	_ = enc.Encode(all)

// 	fmt.Println("Saved to gifs.json, total:", len(all))
// }

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Gif struct {
	Query string `json:"query"`
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Image string `json:"image"`
}

func main() {
	file, err := os.ReadFile("gifs.json")
	if err != nil {
		panic(err)
	}

	var gifs []Gif
	if err := json.Unmarshal(file, &gifs); err != nil {
		panic(err)
	}

	seen := make(map[string]bool)
	unique := make([]Gif, 0, len(gifs))

	for _, g := range gifs {
		if g.ID == "" {
			continue
		}

		if seen[g.ID] {
			continue
		}

		seen[g.ID] = true
		unique = append(unique, g)
	}

	out, err := json.MarshalIndent(unique, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("gifs.unique.json", out, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Done. Unique gifs: %d\n", len(unique))
}