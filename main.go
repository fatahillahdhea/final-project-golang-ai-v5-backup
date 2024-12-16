package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"

	"a21hc3NpZ25tZW50/service"

	"encoding/csv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// Initialize the services
var fileService = &service.FileService{}
var aiService = &service.AIService{Client: &http.Client{}}

// var store = sessions.NewCookieStore([]byte("my-key"))

// func getSession(r *http.Request) *sessions.Session {
// 	session, _ := store.Get(r, "chat-session")
// 	fmt.Println("SESSION", session)
// 	return session
// }

// main.go -> Router -> Controller -> Service -> Repository
// cmd/main.go -> Router -> Controller -> Service -> Repository

type Body struct {
	File     string `json:"file"`
	Question string `json:"question"`
}

/*
TODO
- User/Client -> input file and question
- BE -> buka file
	 -> read file
	 -> simpan ke struct value
	 -> kirim ke folder /Service
	 -> kirim ke repository -> thirdparty

*/

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the Hugging Face token from the environment variables
	token := os.Getenv("HUGGINGFACE_TOKEN")
	if token == "" {
		log.Fatal("HUGGINGFACE_TOKEN is not set in the .env file")
	}

	// Set up the router
	router := mux.NewRouter()

	// File upload endpoint
	router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, handler, err := r.FormFile("file")
		if err != nil {
			panic(err)
		}

		question := r.FormValue("question")

		defer file.Close()

		fmt.Println("FILE => ", file)
		fmt.Println("handler => ", handler)
		fmt.Println("question => ", question)

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			fmt.Println("Error reading records")
		}

		// Loop to iterate through
		// and print each of the string slice
		for _, eachrecord := range records {
			fmt.Println("DATA : ", eachrecord)
			fmt.Println("DATA : ", len(eachrecord))
		}

		// f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		// if err != nil {
		// 	panic(err)
		// }
		// defer f.Close()
		// _, _ = io.WriteString(w, "File "+fileName+" Uploaded successfully")
		// _, _ = io.Copy(f, file)

		w.Write([]byte("hello again"))

	}).Methods("POST")


	// Chat endpoint
	router.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			Query   string `json:"query"` // Input dari user
			Context string `json:"context"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
	
		// Panggil fungsi ChatWithAI
		answer, err := aiService.ChatWithAI(requestBody.Context, requestBody.Query, token)
		if err != nil {
			http.Error(w, fmt.Sprintf("AI chat failed: %v", err), http.StatusInternalServerError)
			return
		}
	
		// Kirim respons
		response := map[string]string{
			"status": "success",
			"answer": answer.GeneratedText,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("POST")
	

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Allow your React app's origin
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler(router)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
