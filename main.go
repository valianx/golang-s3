package main
import (
	"fmt"
	"log"
	"net/http"
)
func handler(w http.ResponseWriter, r *http.Request) {
	maxSize := int64(1024000) // allow only 1MB of file size
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Image too large. Max Size: %v", maxSize)
		return
	}
	file, fileHeader, err := r.FormFile("profile_picture")
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Could not get uploaded file")
		return
	}
	defer file.Close()
	// create an AWS session which can be
	// reused if we're uploading many files
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials(
			"secret-id", // id
			"secret-key",   // secret
			""),  // token can be left blank for now
	})
	if err != nil {
		fmt.Fprintf(w, "Could not upload file")
	}
	fileName, err := UploadFileToS3(s, file, fileHeader)
	if err != nil {
		fmt.Fprintf(w, "Could not upload file")
	}
	fmt.Fprintf(w, "Image uploaded successfully: %v", fileName)
}
func main() {
	http.HandleFunc("/", handler)
	log.Println("Upload server started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
