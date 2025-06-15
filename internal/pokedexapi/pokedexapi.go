package pokedexapi

import(
    "fmt"
	"net/http"
	"io"
	"encoding/json"
)

type LocationsResponse struct {
	Count	int			`json:"count"`
 	Next	string		`json:"next"`
	Prev	string		`json:"previous"`
	Results	[]Location	`json:"results"`
}

type Location struct {
	Name	string	`json:"name"`
	Url		string	`json:"url"`
}

func QueryPokedexApi(url string) (*LocationsResponse, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return nil, err
	}

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	var locations LocationsResponse
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal failed: %v", err)
	}

	return &locations, nil
}
