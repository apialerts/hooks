package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apialerts/hooks/internal/db"
	"github.com/go-chi/chi/v5"
)

var endpointTmpl = template.Must(template.New("endpoint").Parse(layoutStart + `
<main class="max-w-7xl mx-auto w-full px-4 py-4 flex flex-col flex-1 min-h-0">
    <div class="bg-white dark:bg-dark-surface rounded-lg shadow-sm border border-gray-200 dark:border-dark-border p-3 sm:p-4 mb-4 space-y-3">
        <div class="flex flex-col sm:flex-row sm:items-center gap-3">
            <div class="flex items-center gap-2 flex-1 min-w-0">
                <code class="text-sm font-mono bg-gray-100 dark:bg-dark-surface-high px-2 py-1.5 rounded truncate flex-1" id="endpoint-url">{{.EndpointURL}}</code>
                <button onclick="navigator.clipboard.writeText(document.getElementById('endpoint-url').textContent)" class="text-gray-500 hover:text-gray-700 dark:hover:text-dark-text px-2 py-1.5 border border-gray-300 dark:border-dark-border rounded text-xs flex-shrink-0">
                    Copy
                </button>
            </div>
            <div class="flex items-center gap-2 flex-shrink-0">
                <select id="response-mode" onchange="onPresetChange()"
                        class="text-sm bg-white dark:bg-dark-surface-high dark:text-dark-text border border-gray-300 dark:border-dark-border rounded px-2 py-1.5">
                    {{range .Presets}}
                    <option value="{{.Status}}-{{.Delay}}" data-body="{{.DefaultBody}}"
                            {{if and (eq .Status $.CurrentStatus) (eq .Delay $.CurrentDelay)}}selected{{end}}>
                        {{.Label}}
                    </option>
                    {{end}}
                </select>
                <button onclick="testEndpoint()"
                        class="bg-gray-100 dark:bg-dark-surface-high hover:bg-gray-200 dark:hover:bg-dark-surface-highest text-gray-700 dark:text-dark-text-secondary text-xs font-medium px-3 py-1.5 rounded transition-colors">
                    Test
                </button>
                <button hx-delete="/{{.EndpointID}}/delete"
                        hx-confirm="Delete this endpoint and all its data? This cannot be undone."
                        class="text-xs text-gray-400 dark:text-dark-text-muted hover:text-red-500 dark:hover:text-red-500 px-2 py-1.5 transition-colors">
                    Delete
                </button>
            </div>
        </div>
        <div>
            <div class="flex items-center justify-between mb-1">
                <span class="text-xs text-gray-500 dark:text-dark-text-muted">Response body</span>
                <button onclick="saveConfig()"
                        class="bg-orange-500 hover:bg-orange-600 text-white text-xs font-medium px-3 py-1 rounded transition-colors">
                    Save
                </button>
            </div>
            <textarea id="response-body" rows="3"
                      class="w-full text-xs font-mono bg-gray-50 dark:bg-dark-bg dark:text-dark-text border border-gray-200 dark:border-dark-border rounded px-3 py-2 resize-none"
                      placeholder='{"status": "ok"}'>{{.CurrentBody}}</textarea>
        </div>
    </div>

    <script>
        var currentStatus = '{{.CurrentStatus}}';
        var currentDelay = '{{.CurrentDelay}}';

        function resetCountdown() {
            var bar = document.getElementById('countdown-bar');
            bar.classList.remove('animate-countdown');
            void bar.offsetWidth;
            bar.classList.add('animate-countdown');
        }

        function refreshRequests() {
            htmx.trigger('#request-list', 'poll');
            resetCountdown();
        }

        function onPresetChange() {
            var select = document.getElementById('response-mode');
            var option = select.options[select.selectedIndex];
            var parts = option.value.split('-');
            currentStatus = parts[0];
            currentDelay = parts[1];
            document.getElementById('response-body').value = option.getAttribute('data-body');
            saveConfig();
        }

        function saveConfig() {
            var body = document.getElementById('response-body').value;
            htmx.ajax('PUT', '/{{.EndpointID}}/config', {
                values: {status: currentStatus, delay: currentDelay, body: body},
                swap: 'none'
            });
        }

        function testEndpoint() {
            fetch('/{{.EndpointID}}/test', {
                method: 'POST'
            }).then(function() {
                refreshRequests();
            });
        }

        // Reset countdown after every automatic poll too
        document.addEventListener('htmx:afterSwap', function(e) {
            if (e.detail.target.id === 'request-list') {
                resetCountdown();
                // Update request count from response header
                var count = e.detail.xhr.getResponseHeader('X-Request-Count');
                if (count !== null) {
                    document.getElementById('request-count').textContent = count;
                }
            }
        });
    </script>

    <div class="flex flex-col lg:flex-row gap-4 flex-1 min-h-0">
        <div class="lg:w-80 xl:w-96 flex-shrink-0 flex flex-col min-h-0">
            <div class="flex items-center justify-between mb-2">
                <h3 class="font-semibold dark:text-dark-text text-sm">Requests (<span id="request-count">{{.RequestCount}}</span>)</h3>
                <button onclick="refreshRequests()"
                        class="text-sm text-orange-500 hover:text-orange-600 font-medium">
                    Refresh Now
                </button>
            </div>
            <div class="relative mb-2">
                <div class="h-1 bg-gray-200 dark:bg-dark-surface-high rounded-full overflow-hidden">
                    <div id="countdown-bar" class="h-full bg-orange-500 rounded-full animate-countdown"></div>
                </div>
            </div>
            <div id="request-list" class="flex-1 overflow-y-auto bg-white dark:bg-dark-surface rounded-lg shadow-sm border border-gray-200 dark:border-dark-border"
                 hx-get="/{{.EndpointID}}/requests"
                 hx-trigger="load, poll, every 30s"
                 hx-swap="innerHTML">
                <div class="p-8 text-center text-gray-400">Loading...</div>
            </div>
        </div>

        <div id="request-detail" class="flex-1 min-h-0 overflow-y-auto">
            <div class="h-full bg-white dark:bg-dark-surface rounded-lg shadow-sm border border-gray-200 dark:border-dark-border flex items-center justify-center">
                <p class="text-gray-400 dark:text-dark-text-muted text-sm">Select a request to view details</p>
            </div>
        </div>
    </div>
</main>

<style>
@keyframes countdown {
    from { width: 100%; }
    to { width: 0%; }
}
.animate-countdown {
    animation: countdown 30s linear infinite;
}
</style>
` + layoutEnd))

func (h *Handler) HandleEndpoint(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if strings.Contains(r.Header.Get("Accept"), "text/html") && r.Method == http.MethodGet {
		h.serveEndpointUI(w, r, id)
		return
	}

	h.receiveWebhook(w, r, id)
}

func (h *Handler) serveEndpointUI(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	endpoint, err := h.db.GetEndpoint(ctx, id)
	if err != nil {
		notFoundTmpl.Execute(w, map[string]interface{}{
			"BaseURL": h.baseURL,
		})
		return
	}

	currentBody := endpoint.ResponseBody
	if currentBody == "" {
		currentBody = defaultBodyForStatus(endpoint.ResponseStatus, endpoint.ResponseDelay)
	}

	endpointTmpl.Execute(w, map[string]interface{}{
		"BaseURL":       h.baseURL,
		"EndpointID":    endpoint.ID,
		"EndpointURL":   fmt.Sprintf("%s/%s", h.baseURL, endpoint.ID),
		"CurrentStatus": endpoint.ResponseStatus,
		"CurrentDelay":  endpoint.ResponseDelay,
		"CurrentBody":   currentBody,
		"RequestCount":  endpoint.RequestCount,
		"Presets":       responsePresets,
	})
}

func (h *Handler) TestEndpoint(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()

	endpoint, err := h.db.GetEndpoint(ctx, id)
	if err != nil {
		http.Error(w, "endpoint not found", http.StatusNotFound)
		return
	}

	responseBody := endpoint.ResponseBody
	if responseBody == "" {
		responseBody = defaultBodyForStatus(endpoint.ResponseStatus, endpoint.ResponseDelay)
	}
	respStatus := endpoint.ResponseStatus

	seq, err2 := h.db.TouchEndpoint(ctx, id)
	if err2 != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	testBody := `{"event": "test.webhook", "message": "Hello from hooks.apialerts.com!", "timestamp": "` + time.Now().UTC().Format(time.RFC3339) + `"}`
	testHeaders, _ := json.Marshal(map[string][]string{
		"Content-Type":    {"application/json"},
		"X-Test-Request":  {"true"},
		"X-Hook-Endpoint": {id},
	})

	req := &db.Request{
		EndpointID:     id,
		Seq:            seq,
		Method:         "POST",
		Path:           "/",
		QueryParams:    "",
		Headers:        testHeaders,
		Body:           testBody,
		BodySize:       len(testBody),
		ContentType:    "application/json",
		SourceIP:       "hooks.apialerts.com",
		ResponseStatus: &respStatus,
		ResponseBody:   responseBody,
	}

	h.db.CreateRequest(ctx, req)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteEndpoint(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.db.DeleteEndpoint(r.Context(), id); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	r.ParseForm()
	status, _ := strconv.Atoi(r.FormValue("status"))
	delay, _ := strconv.Atoi(r.FormValue("delay"))
	body := r.FormValue("body")

	if status == 0 {
		status = 200
	}

	// If body is empty, use the preset default
	if body == "" {
		body = defaultBodyForStatus(status, delay)
	}

	if err := h.db.UpdateEndpointConfig(r.Context(), id, status, delay, body); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
