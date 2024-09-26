package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

type ResponseHandler struct {
	handler http.Handler
	headerName string
	headerValue string
}

// middleware function

func NewResponseHandler(h http.Handler, headerName string, headerValue string) *ResponseHandler {
	return &ResponseHandler{
		handler: h, 
		headerName: headerName, 
		headerValue: headerValue,
	}
}

func (rh *ResponseHandler) ServeHTTP (w http.ResponseWriter, r *http.Request){
	w.Header().Add(rh.headerName, rh.headerValue)
	rh.handler.ServeHTTP(w, r)
}

type User struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

type Blog struct {
	BlogID     string        `json:"blog_id"`
	Title      string        `json:"title"`
	Content    string        `json:"content"`
	DayCreated time.Weekday  `json:"day_created"`
	CreatedBy  User          `json:"user"`
}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the blog API")
}

var blogs = make(map[string] *Blog)

func (u *User) Blog_detail() *Blog {
	id, err := uuid.NewRandom()

	if err != nil {
		log.Fatal(err)
	}

	blog := &Blog{
		BlogID:     id.String(),
		Title:      "My First Blog",
		Content:    "This is my first blog",
		DayCreated : time.Now().Weekday(),
		CreatedBy : *u,
	}

	blogs[blog.BlogID] = blog
	return blog
}

func Blog_encoder_json(w http.ResponseWriter, r *http.Request){
	// fmt.Println(w, "this is blog encoder")

	user := user_Account()

	blog := user.Blog_detail()

	w.Header().Set("Content-type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET")

	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(blog)

	errCheck(err)

}

func DeleteOneBlog(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "this is delete section")

	w.Header().Set("Content-type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	vars := mux.Vars(r)

	BlogID := vars["id"]

	// fmt.Fprintf(w, "Blog with ID %s is deleted", BlogID)

	if _, exists := blogs[BlogID]; exists {
		delete(blogs, BlogID)
		fmt.Fprintf(w, "Blog with ID %s is deleted", BlogID)
	} else {
		http.Error(w, "Blog not found", http.StatusNotFound)
		fmt.Fprintf(w, "Blog with ID %s is not found", BlogID)
	}

}

func user_Account() *User{
	id, err := uuid.NewRandom()

	errCheck(err)

	return &User{
		UserID: id.String(),
		UserName: "Vinayak",
	}
}

// refractor the main function to
func StartServer() {

	r := mux.NewRouter()

	r.Use(mux.CORSMethodMiddleware(r))

	CustomHandler := NewResponseHandler(r, "X-Request-Id", "12345")

	r.HandleFunc("/", greet).Methods("GET")

	r.HandleFunc("/blog", Blog_encoder_json).Methods("POST","GET")

	r.HandleFunc("/blog/{id}", DeleteOneBlog).Methods("DELETE")

	err := http.ListenAndServe(":8070", CustomHandler)

	errCheck(err)

}

var rootCmd = &cobra.Command{
	Use: "blog",
	Short: "This is a blog API",
}

var startCmd = &cobra.Command{
	Use: "start",
	Short: "Starting the blog",
	Run: func (cmd *cobra.Command, args []string) {
		fmt.Println("Starting the blog server")
		StartServer()
	},
}

func Execute() {

	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	Execute()

}


