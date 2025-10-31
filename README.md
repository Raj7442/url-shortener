# URL Shortener Service (Golang + Docker)

A simple URL shortener service built with **Go (Golang)** using the `net/http` package.  
This project accepts a long URL, generates a short one, and redirects users when they visit the short link.

---

## 🚀 Features

✅ Shortens any valid URL via a REST API  
✅ Returns the **same short URL** for the same input (no duplicates)  
✅ Redirects users to the original URL when short link is opened  
✅ Tracks and returns **Top 3 most-shortened domains**  
✅ In-memory data storage (no database required)  
✅ [BONUS] Dockerized for easy deployment  

---

## 🧰 Tech Stack

- **Language:** Go (Golang)
- **Framework:** net/http
- **Storage:** In-memory maps
- **Containerization:** Docker
- **Testing:** Go’s built-in `testing` package
- **Version Control:** Git + GitHub

---

## 🧠 API Endpoints

| Method | Endpoint | Description |
|--------|-----------|-------------|
| `POST` | `/api/shorten` | Shorten a long URL |
| `GET` | `/api/metrics` | Get top 3 most-shortened domains |
| `GET` | `/{shortCode}` | Redirect to the original long URL |

---

## 🧩 Example Usage

### 1️⃣ Shorten a URL
```bash
curl -X POST http://localhost:8080/api/shorten \
     -H "Content-Type: application/json" \
     -d '{"url":"https://www.example.com"}'
