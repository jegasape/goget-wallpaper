package main

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type ImageManager struct {
	mu        sync.Mutex
	hashes    map[string]bool
	counter   int
	repeated  int
	limit     int
	maxUnique int
}

func NewImageManager(repeatLimit int, maxUnique int) *ImageManager {
	return &ImageManager{
		hashes:    make(map[string]bool),
		limit:     repeatLimit,
		maxUnique: maxUnique,
	}
}

func (im *ImageManager) SaveIfUnique(data []byte) (bool, string) {
	hash := fmt.Sprintf("%x", md5.Sum(data))

	im.mu.Lock()
	defer im.mu.Unlock()

	if im.hashes[hash] {
		im.repeated++
		return false, ""
	}

	im.counter++
	im.hashes[hash] = true
	im.repeated = 0
	fileName := fmt.Sprintf("image_%d.jpg", im.counter)
	err := os.WriteFile("./images/"+fileName, data, 0644)
	if err != nil {
		log.Printf("Error saving image: %v", err)
	}
	return true, fileName
}

func (im *ImageManager) ShouldStop() bool {
	im.mu.Lock()
	defer im.mu.Unlock()
	return im.repeated >= im.limit || (im.maxUnique > 0 && im.counter >= im.maxUnique)
}

func downloadWorker(ctx context.Context, id int, baseUrl string, im *ImageManager,
	sem chan struct{}, wg *sync.WaitGroup) {

	defer wg.Done()

	sem <- struct{}{}
	defer func() { <-sem }()

	url := baseUrl + strings.Repeat("?", id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Worker %d: Error creating request: %v", id, err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Worker %d: Error downloading image: %v", id, err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("[#%d] Invalid HTTP Status %s", id, res.Status)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("[#%d] Can´t read body: %v", id, err)
		return
	}

	if ok, filename := im.SaveIfUnique(body); ok {
		log.Printf("[#%d] New image saved it! %s", id, filename)
	} else {
		log.Printf("[#%d] Repeated image", id)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to get .env")
	}

	encodeURL := os.Getenv("URL")
	if encodeURL == "" {
		log.Fatal("URL is not define")
	}

	urlBytes, err := base64.StdEncoding.DecodeString(encodeURL)
	if err != nil {
		log.Fatalf("Error decoding URL %v", err)
	}

	baseUrl := string(urlBytes)

	const (
		maxConcurrentDownloads = 30
		repeatLimit            = 100
		maxUniqueImages        = 0
	)

	im := NewImageManager(repeatLimit, maxUniqueImages)
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	var wg sync.WaitGroup

	sem := make(chan struct{}, maxConcurrentDownloads)

	for i := 0; ; i++ {
		if im.ShouldStop() {
			break
		}

		wg.Add(1)
		go downloadWorker(ctx, i, baseUrl, im, sem, &wg)

		time.Sleep(50 * time.Millisecond)
	}

	wg.Wait()
	log.Printf("Total images downloaded: %d", im.counter)
}
