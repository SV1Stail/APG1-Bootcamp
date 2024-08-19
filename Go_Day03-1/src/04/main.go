package main

import (
	"04/server"
	"net/http"
)

func main() {

	rootMux := http.NewServeMux()
	getTokenMux := http.NewServeMux()
	rootMux.Handle("/api/", http.StripPrefix("/api", getTokenMux))

	getTokenMux.Handle("/recommend", server.Middleware(http.HandlerFunc(server.Recommend)))
	getTokenMux.HandleFunc("/get_token", server.Token)
	http.ListenAndServe(":8888", rootMux)

}

// curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIiLCJleHAiOjE3MjM4MTk2NTV9.t_ACQtJl5ZMP4499cbBn_cWn1rRyUBOVvUgp1GZdjsY" -X GET http://127.0.0.1:8888/api/recommend?lon=55.674&lat=37.666
