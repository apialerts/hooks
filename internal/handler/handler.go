package handler

import (
	"github.com/apialerts/hooks/internal/db"
)

var reservedPaths = map[string]bool{
	"privacy":     true,
	"endpoints":   true,
	"health":      true,
	"static":      true,
	"robots.txt":  true,
	"sitemap.xml": true,
	"favicon.ico": true,
}

type ResponsePreset struct {
	Status      int
	Delay       int
	Label       string
	DefaultBody string
}

var responsePresets = []ResponsePreset{
	{Status: 200, Delay: 0, Label: "200 OK", DefaultBody: `{"status": "ok"}`},
	{Status: 201, Delay: 0, Label: "201 Created", DefaultBody: `{"status": "created", "id": "example-id"}`},
	{Status: 400, Delay: 0, Label: "400 Bad Request", DefaultBody: `{"error": "bad_request", "message": "Invalid payload"}`},
	{Status: 401, Delay: 0, Label: "401 Unauthorized", DefaultBody: `{"error": "unauthorized", "message": "Authentication required"}`},
	{Status: 403, Delay: 0, Label: "403 Forbidden", DefaultBody: `{"error": "forbidden", "message": "Insufficient permissions"}`},
	{Status: 404, Delay: 0, Label: "404 Not Found", DefaultBody: `{"error": "not_found", "message": "Resource not found"}`},
	{Status: 500, Delay: 0, Label: "500 Internal Server Error", DefaultBody: `Internal Server Error`},
	{Status: 503, Delay: 0, Label: "503 Service Unavailable", DefaultBody: `Service Temporarily Unavailable`},
	{Status: 200, Delay: 35000, Label: "Timeout (35s)", DefaultBody: ``},
}

func defaultBodyForStatus(status, delay int) string {
	for _, p := range responsePresets {
		if p.Status == status && p.Delay == delay {
			return p.DefaultBody
		}
	}
	return `{"status": "ok"}`
}

type Handler struct {
	db      *db.Queries
	baseURL string
}

func New(db *db.Queries, baseURL string) *Handler {
	return &Handler{db: db, baseURL: baseURL}
}
