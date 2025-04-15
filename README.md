# ⚡ 2g-Wallpaper

A high-performance, concurrent wallpaper downloader written in Go, Avoids duplicate downloads, stops automatically after too many repeats, and securely loads in a `.env` file.

---

## 🚀 Features

- 🔄 Concurrent downloads (configurable)
- 🔍 Duplicate detection using image hashing (MD5)
- 🧠 Auto-stops after a configurable number of repeated images
- 🔐 Loads URL from `.env` 
- 📁 Saves all images locally as `./images/imagen_x.jpg`

---

## 🛠 Requirements

- Go 1.18 or higher
- External module for `.env` support:

```bash
go get github.com/joho/godotenv

