package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

var homeTmpl = template.Must(template.New("home").Parse(layoutStart + `
<main class="max-w-4xl mx-auto px-4 py-16 text-center">
    <h1 class="text-4xl font-bold dark:text-dark-text mb-4">Free Webhook Tester</h1>
    <p class="text-xl text-gray-600 dark:text-dark-text-secondary mb-8">
        Instantly generate a unique URL to receive, inspect, and debug webhook requests.
        Toggle between success and failure responses to test your retry logic.
    </p>
    <form method="POST" action="/endpoints">
        <button type="submit" class="bg-orange-500 hover:bg-orange-600 text-white font-semibold py-3 px-8 rounded-lg text-lg transition-colors">
            Create Endpoint
        </button>
    </form>
    <div class="mt-12 grid grid-cols-1 md:grid-cols-3 gap-8 text-left">
        <div class="p-6 bg-white dark:bg-dark-surface rounded-lg shadow-sm border border-gray-200 dark:border-dark-border">
            <h3 class="font-semibold dark:text-dark-text mb-2">Inspect Requests</h3>
            <p class="text-gray-600 dark:text-dark-text-secondary text-sm">See every header, body, and query parameter. Formatted JSON, real-time updates.</p>
        </div>
        <div class="p-6 bg-white dark:bg-dark-surface rounded-lg shadow-sm border border-gray-200 dark:border-dark-border">
            <h3 class="font-semibold dark:text-dark-text mb-2">Toggle Responses</h3>
            <p class="text-gray-600 dark:text-dark-text-secondary text-sm">Switch between 200, 400, 500, timeout, and more. Test your error handling instantly.</p>
        </div>
        <div class="p-6 bg-white dark:bg-dark-surface rounded-lg shadow-sm border border-gray-200 dark:border-dark-border">
            <h3 class="font-semibold dark:text-dark-text mb-2">No Sign-Up</h3>
            <p class="text-gray-600 dark:text-dark-text-secondary text-sm">Generate an endpoint instantly. No account needed. Endpoints persist for 14 days.</p>
        </div>
    </div>
</main>
` + layoutEnd))

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	homeTmpl.Execute(w, map[string]interface{}{
		"BaseURL": h.baseURL,
	})
}

func (h *Handler) CreateEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ip := r.RemoteAddr

	count, err := h.db.CountEndpointsByIP(ctx, ip)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if count >= 5 {
		http.Error(w, "maximum endpoints per IP reached (5)", http.StatusTooManyRequests)
		return
	}

	id := generateSlug()
	for reservedPaths[id] {
		id = generateSlug()
	}

	if err := h.db.CreateEndpoint(ctx, id, ip); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s", id), http.StatusSeeOther)
}
