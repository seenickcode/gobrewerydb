package brewerydb

import (
	"fmt"
	"os"
	"testing"
)

func TestSearchBreweryDB(t *testing.T) {

	apiKey := os.Getenv("BREWERYDB_KEY")
	if len(apiKey) == 0 {
		panic("BREWERYDB_KEY API key not set, go and register to get one at brewerydb.com")
	}
	query := "guinness"
	page := 0
	collected := 0
	client := NewClient(apiKey)
	client.VerboseMode = true

	for {
		response := client.SearchBeers(query, page)
		if len(response.Beers) == 0 {
			break
		}

		for _, beer := range response.Beers {
			fmt.Printf("got beer '%v'\n", beer.Name)
			collected++
		}
		page++
	}

	if collected == 0 {
		t.Errorf("no beers returned for query '%v' when there should have been\n", query)
	}
}
