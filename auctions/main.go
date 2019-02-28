package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"
)

// BidData : bid-data struct
type BidData struct {
	ID          string `json:"id"`
	PlacementID string `json:"placementID"`
	BidPrice    int    `json:"bidPrice"`
	Currency    string `json:"currency"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/call-bidders", CallBidders)

	log.Println("Starting Server")
	log.Fatal(http.ListenAndServe(":8000", router))
}

// CallBidders : this is the handler for calling bidders.
// makes multipe calls to the bidders at the same time.
func CallBidders(w http.ResponseWriter, r *http.Request) {
	AllBids := []BidData{}
	ch := make(chan BidData)
	// url := os.Getenv("BIDDER_URL")
	for i := 1; i <= 10; i++ {
		go MakeRequest("http://127.0.0.1:5000/make-bid?placement-id=123", ch)
	}

	for i := 1; i <= 10; i++ {
		AllBids = append(AllBids, <-ch)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SortBids(AllBids))
}

// MakeRequest : function to actually call bidding service.
func MakeRequest(url string, ch chan BidData) {
	var Bid BidData
	timeOver := time.Duration(200 * time.Millisecond)
	client := http.Client{
		Timeout: timeOver,
	}
	resp, err := client.Get(url)
	if err, ok := err.(net.Error); ok && err.Timeout() {
		ch <- Bid
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode == 204 || resp.StatusCode == 400 {
		ch <- Bid
		return
	}
	if err := json.NewDecoder(resp.Body).Decode(&Bid); err != nil {
		ch <- Bid
		return
	}
	ch <- Bid
}

// SortBids : this sorts the bids and returns highest bid
func SortBids(allBids []BidData) []BidData {
	sort.Slice(allBids[:], func(i, j int) bool {
		return allBids[i].BidPrice > allBids[j].BidPrice
	})
	return allBids
}
