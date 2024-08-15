package main

import (
	"03/server"
	"fmt"
	"net/http"
)

func main() {
	// client, _ := elasticsearch.NewDefaultClient()

	// store := db.NewElasticsearchStore(client, "places")
	// places, err := store.GetNearesRestaurant(55.674, 37.666, 3)
	// if err != nil {
	// 	fmt.Printf("%v", err)
	// 	os.Exit(1)
	// }

	// for _, place := range places {
	// 	fmt.Println(place)
	// }

	http.HandleFunc("/", server.Handler)
	fmt.Println("Server is listening on port 8888...")
	http.ListenAndServe(":8888", nil)
}
