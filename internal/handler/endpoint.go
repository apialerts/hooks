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
<main class="max-w-7xl mx-auto w-full px-4 sm:px-6 py-4 flex flex-col flex-1 min-h-0">
    <!-- Endpoint URL -->
    <p class="text-xs font-medium text-gray-600 dark:text-dark-text-secondary mb-1.5">Your webhook endpoint. Expires 7 days after last activity.</p>
    <div class="inline-flex items-center gap-2 mb-1.5 self-start">
        <div class="inline-flex items-center gap-3 rounded-lg bg-white dark:bg-dark-surface border border-gray-200 dark:border-dark-border px-3 py-2">
            <code class="text-sm font-mono text-gray-700 dark:text-dark-text">{{.EndpointURL}}</code>
            <button onclick="copyText('{{.EndpointURL}}', this)"
                    class="text-gray-400 dark:text-dark-text-muted hover:text-gray-600 dark:hover:text-dark-text-secondary transition-colors flex-shrink-0" title="Copy URL">
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/></svg>
            </button>
        </div>
        <button hx-delete="/{{.EndpointID}}/delete"
                hx-confirm="Delete this endpoint and all its data? This cannot be undone."
                class="text-gray-400 dark:text-dark-text-muted hover:text-red-500 dark:hover:text-red-500 p-2 rounded-lg hover:bg-red-50 dark:hover:bg-red-500/10 transition-colors flex-shrink-0"
                title="Delete endpoint">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
        </button>
    </div>
    <div class="flex items-center gap-2 mb-4 flex-wrap">
        <p class="text-[11px] text-gray-400 dark:text-dark-text-muted">Send POST, PUT, PATCH, or DELETE requests here.</p>
        <button onclick="navigator.clipboard.writeText('curl -X POST {{.EndpointURL}} -H \'Content-Type: application/json\' -d \'{\x22event\x22: \x22test\x22}\'');this.textContent='Copied!';setTimeout(function(){document.getElementById('curl-hint').textContent='Copy curl example'},1500)" id="curl-hint"
                class="text-[11px] text-brand hover:underline cursor-pointer">Copy curl example</button>
    </div>
    {{if gt (len .Siblings) 1}}
    <div class="flex items-center gap-2 mb-4 flex-wrap">
        {{range .Siblings}}
        {{if eq .ID $.EndpointID}}
        <span class="text-xs font-mono font-semibold text-brand px-2 py-1 rounded-lg bg-brand/10">{{.ID}}</span>
        {{else}}
        <a href="/{{.ID}}" class="text-xs font-mono text-gray-500 dark:text-dark-text-muted hover:text-gray-700 dark:hover:text-dark-text px-2 py-1 rounded-lg hover:bg-gray-100 dark:hover:bg-dark-surface-high transition-colors">{{.ID}}</a>
        {{end}}
        {{end}}
    </div>
    {{else}}
    <div class="mb-4"></div>
    {{end}}

    <!-- Response config -->
    <p class="text-xs font-medium text-gray-600 dark:text-dark-text-secondary mb-1.5">Choose how this endpoint responds to incoming webhooks</p>
    <div class="rounded-xl bg-white dark:bg-dark-surface border border-gray-200 dark:border-dark-border mb-4">
        <div class="flex flex-wrap items-center gap-2 px-4 py-2.5 border-b border-gray-100 dark:border-dark-border">
            <span class="text-xs font-medium text-gray-500 dark:text-dark-text-muted">Response</span>
            <select id="response-mode" onchange="onPresetChange()"
                    class="text-xs font-semibold dark:text-dark-text bg-gray-50 dark:bg-dark-surface-high rounded-lg pl-2.5 pr-8 py-1.5 border border-gray-200 dark:border-dark-border focus:outline-none cursor-pointer flex-1 min-w-[140px]">
                {{range .Presets}}
                <option value="{{.Status}}-{{.Delay}}" data-body="{{.DefaultBody}}"
                        class="bg-white dark:bg-dark-surface"
                        {{if and (eq .Status $.CurrentStatus) (eq .Delay $.CurrentDelay)}}selected{{end}}>
                    {{.Label}}
                </option>
                {{end}}
            </select>
            <div class="flex items-center gap-2 ml-auto">
                <button onclick="testEndpoint()"
                        class="text-xs font-semibold text-gray-900 dark:text-dark-text px-3 py-1 rounded-full border border-gray-300 dark:border-dark-text-muted hover:bg-gray-50 dark:hover:bg-dark-surface-high transition-colors whitespace-nowrap">
                    Send Test
                </button>
                <button onclick="saveConfig()"
                        class="bg-brand hover:brightness-110 text-black text-xs font-semibold px-3.5 py-1 rounded-full transition-all whitespace-nowrap">
                    Save
                </button>
            </div>
        </div>
        <textarea id="response-body" rows="4"
                  class="w-full text-xs font-mono text-gray-700 dark:text-dark-text bg-transparent px-4 py-3 resize-none focus:outline-none leading-relaxed"
                  placeholder='{"status": "ok"}'
                  spellcheck="false">{{.CurrentBody}}</textarea>
    </div>

    <script>
        var currentStatus = '{{.CurrentStatus}}';
        var currentDelay = '{{.CurrentDelay}}';

        function copyText(text, btn) {
            navigator.clipboard.writeText(text);
            var orig = btn.innerHTML;
            btn.innerHTML = '<svg class="w-4 h-4 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/></svg>';
            setTimeout(function() { btn.innerHTML = orig; }, 1500);
        }

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

        function formatBody() {
            var ta = document.getElementById('response-body');
            try {
                var obj = JSON.parse(ta.value);
                ta.value = JSON.stringify(obj, null, 2);
            } catch(e) {}
        }
        formatBody();

        function onPresetChange() {
            var select = document.getElementById('response-mode');
            var option = select.options[select.selectedIndex];
            var parts = option.value.split('-');
            currentStatus = parts[0];
            currentDelay = parts[1];
            document.getElementById('response-body').value = option.getAttribute('data-body');
            formatBody();
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

        document.addEventListener('htmx:afterSwap', function(e) {
            if (e.detail.target.id === 'request-list') {
                resetCountdown();
                var count = e.detail.xhr.getResponseHeader('X-Request-Count');
                if (count !== null) {
                    document.getElementById('request-count').textContent = count;
                }
            }
            if (e.detail.target.id === 'request-detail') {
                if (window.innerWidth < 1024) {
                    document.getElementById('request-detail').scrollIntoView({behavior: 'smooth', block: 'start'});
                }
            }
        });
    </script>

    <!-- Two-column layout -->
    <div class="flex flex-col lg:flex-row gap-4 flex-1 min-h-0">
        <!-- Request list sidebar -->
        <div class="lg:w-80 xl:w-96 flex-shrink-0 flex flex-col min-h-0">
            <div class="flex items-center justify-between mb-2">
                <div class="flex items-center gap-2">
                    <h3 class="font-semibold dark:text-dark-text text-sm">Requests</h3>
                    <span class="text-xs font-mono bg-gray-100 dark:bg-dark-surface-high text-gray-500 dark:text-dark-text-muted px-1.5 py-0.5 rounded" id="request-count">{{.RequestCount}}</span>
                </div>
                <button onclick="refreshRequests()"
                        class="inline-flex items-center gap-1 text-xs text-brand hover:text-brand font-medium transition-colors">
                    <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/></svg>
                    Refresh
                </button>
            </div>
            <div class="relative mb-2">
                <div class="h-0.5 bg-gray-100 dark:bg-dark-surface-high rounded-full overflow-hidden">
                    <div id="countdown-bar" class="h-full bg-brand rounded-full animate-countdown"></div>
                </div>
            </div>
            <div id="request-list" class="flex-1 overflow-y-auto rounded-xl bg-white dark:bg-dark-surface border border-gray-200 dark:border-dark-border min-h-[300px] max-h-[50vh] lg:max-h-none"
                 hx-get="/{{.EndpointID}}/requests"
                 hx-trigger="load, poll, every 30s"
                 hx-swap="innerHTML">
                <div class="p-8 text-center text-gray-400 dark:text-dark-text-muted text-sm">Loading...</div>
            </div>
        </div>

        <!-- Request detail panel -->
        <div id="request-detail" class="flex-1 min-h-0 overflow-y-auto">
            <div class="h-full min-h-[200px] rounded-xl bg-white dark:bg-dark-surface border border-gray-200 dark:border-dark-border flex items-center justify-center">
                <div class="text-center">
                    <svg class="w-8 h-8 text-gray-300 dark:text-dark-border mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/></svg>
                    <p class="text-gray-400 dark:text-dark-text-muted text-sm">Select a request to view details</p>
                </div>
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
			"BaseURL":         h.baseURL,
			"PageTitle":       "Endpoint Not Found",
			"PageDescription": "",
			"CanonicalURL":    h.baseURL + "/",
		})
		return
	}

	currentBody := endpoint.ResponseBody
	if currentBody == "" {
		currentBody = defaultBodyForStatus(endpoint.ResponseStatus, endpoint.ResponseDelay)
	}

	siblings, _ := h.db.ListEndpointsByIP(ctx, r.RemoteAddr)

	endpointTmpl.Execute(w, map[string]interface{}{
		"BaseURL":         h.baseURL,
		"EndpointID":      endpoint.ID,
		"EndpointURL":     fmt.Sprintf("%s/%s", h.baseURL, endpoint.ID),
		"CurrentStatus":   endpoint.ResponseStatus,
		"CurrentDelay":    endpoint.ResponseDelay,
		"CurrentBody":     currentBody,
		"RequestCount":    endpoint.RequestCount,
		"Presets":         responsePresets,
		"Siblings":        siblings,
		"PageTitle":       "",
		"PageDescription": "",
		"CanonicalURL":    h.baseURL + "/",
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
	ctx := r.Context()

	if err := h.db.DeleteEndpoint(ctx, id); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Redirect to another endpoint if one exists, otherwise home
	redirect := "/"
	siblings, _ := h.db.ListEndpointsByIP(ctx, r.RemoteAddr)
	for _, s := range siblings {
		if s.ID != id {
			redirect = "/" + s.ID
			break
		}
	}

	w.Header().Set("HX-Redirect", redirect)
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
