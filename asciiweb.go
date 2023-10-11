package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Art struct {
	Output string
	Color  string
	Back   string
}

func main() {
	styles := http.FileServer(http.Dir("./template/css"))
	http.Handle("/css/", http.StripPrefix("/css/", styles))

	// Register HTTP endpoints and corresponding handler functions
	http.HandleFunc("/style.css", stylee)
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ascii-art", asciipage)
	http.HandleFunc("/error.css", csserror) // Register the "/error.css" endpoint here

	fmt.Println("Listening on port :8800...")
	// Start the server on port 8800
	err := http.ListenAndServe(":8800", nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}

	log.Fatal(http.ListenAndServe(":8800", nil))
}

func stylee(w http.ResponseWriter, r *http.Request) {
	// Serve the "style.css" file
	http.ServeFile(w, r, "template/css/style.css")
}
func csserror(w http.ResponseWriter, r *http.Request) {
	// Serve the "error.css" file
	http.ServeFile(w, r, "template/css/error.css")
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	// Check if the URL path is not "/"
	

	if r.URL.Path != "/" {
		
		w.WriteHeader(http.StatusNotFound)
		// Serve the "error.html" file
		http.ServeFile(w, r, "template/error.html")
		return
	}
	 

	if r.Method != http.MethodGet {
		// http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.ServeFile(w, r, "template/error1.html")
		return
	}

	// Parse and execute the "webhtml.html" template
	tmplt, err := template.ParseFiles("webhtml.html")
	if err != nil {
		// http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "template/error2.html")
		return
	}
	tmplt.Execute(w, nil)
}

func asciipage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.ServeFile(w, r, "template/error.html")
		return
	} else {

	
	// Your existing code for handling the ASCII art conversion
	// Perform the ASCII art conversion using handleASCII function
	theTEXT := handleASCII(w, r)
	picker := r.FormValue("colorPicker")
	BackGroundColor := r.FormValue("background")

	
	// Create an instance of Art struct to hold the ASCII art
	art := Art{Output: theTEXT,
		Color: picker,
		Back:  BackGroundColor}

	// Parse and execute the "webhtml.html" template with the art data
	tmplt, err := template.ParseFiles("webhtml.html")
	if err != nil {
		// http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		w.WriteHeader(505)
		http.ServeFile(w, r, "template/error2.html")
		return
	}
	tmplt.Execute(w, art)
}

}
