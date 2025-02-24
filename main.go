package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var images []string
var validExtensions = getValidExtensions()

func main() {
	dirPath := os.Getenv("IMAGE_DIR")
	if dirPath == "" {
		dirPath = "/images"
	}
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Image directory: %s\n", dirPath)
	for _, file := range files {
		if fileIsValid(path.Join(dirPath, file.Name())) {
			images = appendIfNotExists(images, path.Join(dirPath, file.Name()))
		} else {
			log.Printf("File %s unreadable or not an image, ignoring.\n", file.Name())
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				file := filenameFromEvent(event)
				if event.Has(fsnotify.Create) {
					// happens on rename and new file add
					if fileIsValid(file) {
						images = appendIfNotExists(images, file)
					}
				} else if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
					// happens on rm or old file rename
					images = remove(images, file)
				} else if event.Has(fsnotify.Chmod) {
					// happens on chmod and rm if file descriptors are open
					// file is deleted or unreadable, remove
					if !fileIsValid(file) {
						images = remove(images, file)
					} else {
						// file was unreadable but now is, add
						images = appendIfNotExists(images, file)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", serveImage)
	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveImage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, images[rand.Intn(len(images))])
}

// not fast at all but whatever
func remove(s []string, item string) []string {
	for i, fn := range s {
		if fn == item {
			log.Printf("REMOVE %s\n", item)
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func appendIfNotExists(slice []string, element string) []string {
	for _, ele := range slice {
		if ele == element {
			return slice
		}
	}
	log.Printf("ADD %s\n", element)
	return append(slice, element)
}

func filenameFromEvent(event fsnotify.Event) string {
	return strings.Split(event.String(), "\"")[1]
}

func fileIsValid(file string) bool {
	return slices.Contains(validExtensions, filepath.Ext(file)) && isReadable(file)
}

func getValidExtensions() []string {
	env := os.Getenv("ALLOWED_EXTENSIONS")
	if env == "" {
		return []string{".png", ".jpg", ".jpeg", ".webp"}
	} else {
		return strings.Split(env, ",")
	}
}

// maybe expensive but works
// tried file info mode perms, syscall.Access already, didnt work
func isReadable(file string) bool {
	_, err := os.Open(file)
	return err == nil
}
