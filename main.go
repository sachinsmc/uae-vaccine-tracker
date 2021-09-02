package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Data []struct {
	PeopleFullyVaccinatedPerHundred float64 `json:"people_fully_vaccinated_per_hundred,omitempty"`
}

type Vaccinations []struct {
	Country string `json:"country"`
	Data    Data   `json:"data"`
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func tweet(tweet string) {
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)
	client.Statuses.Update(tweet, nil)
}

func main() {

	countries := []string{"United States", "India", "United Kingdom", "United Arab Emirates", "Philippines", "Israel", "Australia"}
	resp, err := http.Get("https://raw.githubusercontent.com/owid/covid-19-data/master/public/data/vaccinations/vaccinations.json")
	if err != nil {
		log.Fatal(err)
	}
	var vaccinations Vaccinations
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal([]byte(body), &vaccinations)
	for _, vaccination := range vaccinations {
		if contains(countries, vaccination.Country) {
			lastUpdated := len(vaccination.Data) - 1
			var fullyVaccinated float64
			peopleFullyVaccinated := vaccination.Data[lastUpdated].PeopleFullyVaccinatedPerHundred
			if peopleFullyVaccinated > 0 {
				fullyVaccinated = peopleFullyVaccinated * 20 / 100
			} else {
				peopleFullyVaccinated = getLatestData(vaccination.Data)
				fullyVaccinated = peopleFullyVaccinated * 20 / 100
			}
			tweetStr := vaccination.Country + " \n " + strings.Repeat("▓", int(math.Round(fullyVaccinated))) + strings.Repeat("░", int(math.Round(20-fullyVaccinated))) + " " + fmt.Sprintf("%.2f", peopleFullyVaccinated) + "%"
			fmt.Println(tweetStr)
			tweet(tweetStr)
		}
	}
}

func getLatestData(d Data) float64 {
	var peopleFullyVaccinated float64
	for i := len(d) - 1; i > 0; i-- {
		peopleFullyVaccinated = d[i].PeopleFullyVaccinatedPerHundred
		if peopleFullyVaccinated != 0 {
			fmt.Println(peopleFullyVaccinated)
			break
		}
	}
	return peopleFullyVaccinated
}
