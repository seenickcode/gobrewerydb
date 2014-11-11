package brewerydb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const defaultBaseUrl string = "http://api.brewerydb.com/v2"

type breweryDBClient struct {
	apiKey      string
	baseUrl     string
	VerboseMode bool
}

type SearchResponse struct {
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
}

type BeerSearchResponse struct {
	SearchResponse
	Beers []Beer `json:"Data"`
}

type Beer struct {
	Name           string
	ABV            string
	IBU            string
	Style          Style
	Available      Available
	Breweries      []Brewery
	SocialAccounts []SocialAccount
}

type Style struct {
	Name string
}

type Available struct {
	Name string
}

type Brewery struct {
	Name      string
	Website   string
	Locations []Location
}

type Location struct {
	Locality  string
	Region    string
	IsPrimary string
}

type SocialAccount struct {
	Link string
}

func NewClient(apiKey string) (c *breweryDBClient) {
	c = new(breweryDBClient)
	c.apiKey = apiKey
	c.baseUrl = defaultBaseUrl
	c.VerboseMode = false
	return c
}

func (c *breweryDBClient) SearchBeers(q string, pg int) (resp BeerSearchResponse) {

	// set up query string then url
	v := url.Values{}
	v.Set("type", "beer")
	v.Add("withBreweries", "Y")      // add premium features even
	v.Add("withSocialAccounts", "Y") // if user isn't signed up for them
	v.Add("q", q)
	// brewerydb's page param isn't zero based unfortunately
	// (but if you pass in 0 just in case it will still assume it's '1')
	v.Add("p", strconv.Itoa(pg+1))
	v.Add("key", c.apiKey)
	url := c.baseUrl + "/search?" + v.Encode()

	// perform request and convert response to an object
	data, err := c.performGetRequest(url)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	// deserialize to objects
	err = json.Unmarshal(data, &resp)
	if err != nil {
		fmt.Printf("json err: %v\n", err)
		return
	}

	// report our search results
	c.log("fetched pg %d (%d results spanning %d pages)\n", resp.CurrentPage, resp.TotalResults, resp.NumberOfPages)

	return
}

func (c *breweryDBClient) performGetRequest(url string) (buf []byte, err error) {
	c.log("fetching: %v\n", url)
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("http err: %v\n", err)
		return
	}
	if res == nil {
		err = fmt.Errorf("err, response is nil")
		return
	}
	buf, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("ioutil err: %v\n", err)
	}
	return
}

func (c *breweryDBClient) log(format string, a ...interface{}) {
	if c.VerboseMode {
		fmt.Printf(format, a...)
	}
}
