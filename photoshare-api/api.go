package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Tag int

const (
	Selfie Tag = iota
	Animal
	Scenary
	Vacation
)

func (t Tag) String() string {
	return [...]string{"Selfie", "Animal", "Scenary", "Vacation"}[t]
}

type Image struct {
	Path	string `json:"path"`
	Tags	[]string
}

//Fields with 'omitempty' are optional fields.
type Response struct {
	Message	string	`json:"message"`
	Images *[]Image `json:"images,omitempty"`
	Image *Image `json:"image,omitempty"`
}


func getImage() {
	
}

func createImage(res http.ResponseWriter) {
	//Generated a new uuid for the image
	uuid := uuid.New().String()

	//
	resp, _ := json.Marshal(Response{Message: fmt.Sprintf("Image %s added", uuid), Image: &Image{Path: "test", Tags: []string{Selfie.String(), Vacation.String()}}})
	res.Write(resp)
}

func deleteImage() {
	
}

func send404() {

}

func imageHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		getImage()
	case http.MethodPost:
		createImage(res)
	case http.MethodDelete:
		deleteImage()
	default:
		send404()
	}
}

func startAPI() {
	http.HandleFunc("/image", imageHandler);
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	startAPI()
}