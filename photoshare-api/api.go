package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type Tag int

var AllImagesByUuid map[string]Image
var AllImagesByTag map[string]map[string][]Image

type Image struct {
	Path string   `json:"path"`
	Tags []string `json:"tags"`
}

//Fields with 'omitempty' are optional fields.
type Response struct {
	Message string   `json:"message"`
	Images  *[]Image `json:"images,omitempty"`
	Image   *Image   `json:"image,omitempty"`
}

func getImage() {

}

func createImage(res http.ResponseWriter, req *http.Request) {
	//Generated a new uuid for the image
	uuid := uuid.New().String()
	//Makes the max size of form 32MB
	req.ParseMultipartForm(32 << 20)

	//Upload the Image
	file, _, err := req.FormFile("image")
	if err != nil {
		send500(res, "Error reading file. Please try again in 1 minute.")
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		send500(res, "Error reading file. Please try again in 1 minute.")
	}
	path := fmt.Sprintf("images/%s.png", uuid)
	err = ioutil.WriteFile(path, fileBytes, 0477)
	if err != nil {
		send500(res, "Error uploading file. Please try again in 1 minute.")
	}
	tags := req.Form["tags"]
	image := Image{Path: path, Tags: tags}

	//Add image to maps
	AllImagesByUuid[uuid] = image
	for _, tag := range tags {
		if AllImagesByTag[tag] == nil {
			AllImagesByTag[tag] = make(map[string][]Image)
			AllImagesByTag[tag][uuid] = make([]Image, 0)
		}
		AllImagesByTag[tag][uuid] = append(AllImagesByTag[tag][uuid], image)
	}

	//Respond to user
	resp, _ := json.Marshal(Response{Message: fmt.Sprintf("Image %s added", uuid), Image: &image})
	res.Write(resp)
}

func deleteImage(res http.ResponseWriter, req *http.Request) {
	//Delete image from maps
	req.ParseForm()
	uuid, ok := req.URL.Query()["uuid"]
	if !ok || len(uuid[0]) < 1 {
		send400(res, "Url Param 'uuid' is missing")
		return
	}
	delete(AllImagesByUuid, uuid[0])

	for tag := range AllImagesByTag {
		delete(AllImagesByTag[tag], uuid[0])
	}

	//Delete image from system
	err := os.Remove(fmt.Sprintf("images/%s.png", uuid[0]))
	if err != nil {
		log.Println(fmt.Sprintf("ERROR: Failed to delete image %s", uuid[0]))
	}

	//Respond to user
	resp, _ := json.Marshal(Response{Message: fmt.Sprintf("Image %s deleted", uuid[0])})
	res.Write(resp)
}

func send500(res http.ResponseWriter, msg string) {
	res.WriteHeader(http.StatusInternalServerError)
	res.Write([]byte(msg))
}

func send400(res http.ResponseWriter, msg string) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte(msg))
}

func send404(res http.ResponseWriter, msg string) {
	res.WriteHeader(http.StatusNotFound)
	res.Write([]byte(msg))
}

func imageHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		getImage()
	case http.MethodPost:
		createImage(res, req)
	case http.MethodDelete:
		deleteImage(res, req)
	default:
		send404(res, "Page not Found")
	}
}

func startAPI() {
	http.HandleFunc("/image", imageHandler)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	AllImagesByTag = make(map[string]map[string][]Image)
	AllImagesByUuid = make(map[string]Image)
	startAPI()
}
