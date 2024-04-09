package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type urlShortener struct {
	urls map[string]string
}

func (u *urlShortener) HandleShorten(w http.ResponseWriter, r *http.Request) {
	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortKey := generateShortKey()
	u.urls[shortKey] = originalURL

	shortenedURL := fmt.Sprintf("http://localhost:8080/%s", shortKey)

	w.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(`
		<h2>URL Shortener</h2>
		<p>Original URL: %s</p>
		<p>Shortened URL: <a href="%s">%s</a></p>
		<form method="post" action="/shorten">
			<input type="text" name="url" placeholder="Enter a URL">
			<input type="submit" value="Shorten">
		</form>
	`, originalURL, shortenedURL, shortenedURL)
	fmt.Fprint(w, responseHTML)
}

func (u *urlShortener) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := r.PathValue("url")

	originalURL, found := u.urls[shortKey]
	if !found {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>URL Shortener</title>
		</head>
		<body>
			<h2>URL Shortener</h2>
			<form method="post" action="/shorten">
				<input type="url" name="url" placeholder="Enter a URL" required>
				<input type="submit" value="Shorten">
			</form>
		</body>
		</html>
	`)
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 3

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[r.Intn(len(charset))]
	}
	return string(shortKey)
}

func main() {
	shortener := &urlShortener{
		urls: make(map[string]string),
	}

	http.HandleFunc("GET /", handleForm)
	http.HandleFunc("POST /shorten", shortener.HandleShorten)
	http.HandleFunc("GET /{url}/", shortener.HandleRedirect)

	fmt.Printf("Server is running on http://localhost:8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}
