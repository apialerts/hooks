package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/apialerts/hooks/internal/db"

	"github.com/go-chi/chi/v5"
)

var requestListTmpl = template.Must(template.New("request_list").Funcs(template.FuncMap{
	"formatTime": func(t time.Time) string {
		return t.Format("15:04:05")
	},
	"formatSize": func(size int) string {
		if size < 1024 {
			return fmt.Sprintf("%dB", size)
		}
		return fmt.Sprintf("%.1fKB", float64(size)/1024)
	},
	"statusColor": func(status *int) string {
		if status == nil {
			return "text-gray-500 dark:text-dark-text-muted"
		}
		s := *status
		if s >= 200 && s < 300 {
			return "text-green-600 dark:text-green-400"
		}
		if s >= 400 && s < 500 {
			return "text-yellow-600 dark:text-yellow-400"
		}
		if s >= 500 {
			return "text-red-600 dark:text-red-400"
		}
		return "text-gray-600 dark:text-dark-text-secondary"
	},
	"derefStatus": func(status *int) string {
		if status == nil {
			return ""
		}
		return strconv.Itoa(*status)
	},
	"methodColor": func(method string) string {
		switch method {
		case "GET":
			return "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400"
		case "POST":
			return "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400"
		case "PUT":
			return "bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400"
		case "PATCH":
			return "bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400"
		case "DELETE":
			return "bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400"
		default:
			return "bg-gray-100 text-gray-700 dark:bg-dark-surface-high dark:text-dark-text"
		}
	},
}).Parse(`
{{if .Requests}}
<div class="divide-y divide-gray-200 dark:divide-dark-border">
    {{range .Requests}}
    <div class="flex items-center gap-3 px-3 py-2.5 cursor-pointer hover:bg-gray-50 dark:hover:bg-dark-surface-high transition-colors"
         hx-get="/{{$.EndpointID}}/requests/{{.Seq}}"
         hx-target="#request-detail"
         hx-swap="innerHTML">
        <span class="px-2 py-0.5 rounded text-xs font-semibold w-16 text-center {{methodColor .Method}}">{{.Method}}</span>
        <span class="text-xs font-mono text-gray-500 dark:text-dark-text-secondary">{{formatTime .ReceivedAt}}</span>
        {{if .ResponseStatus}}<span class="text-xs font-mono font-semibold {{statusColor .ResponseStatus}}">{{derefStatus .ResponseStatus}}</span>{{end}}
        <span class="text-xs text-gray-400 dark:text-dark-text-muted ml-auto">{{formatSize .BodySize}}</span>
    </div>
    {{end}}
</div>
{{else}}
<div class="p-8 text-center text-gray-400 dark:text-dark-text-muted">
    <p class="text-lg mb-2">No requests yet</p>
    <p class="text-sm">Send a webhook to your endpoint URL to see it here</p>
</div>
{{end}}
`))

var requestDetailTmpl = template.Must(template.New("request_detail").Funcs(template.FuncMap{
	"formatTime": func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	},
	"formatHeaderRows": func(raw []byte) []map[string]string {
		var headers map[string][]string
		if err := json.Unmarshal(raw, &headers); err != nil {
			return nil
		}
		var rows []map[string]string
		for k, v := range headers {
			rows = append(rows, map[string]string{"Key": k, "Value": strings.Join(v, ", ")})
		}
		sort.Slice(rows, func(i, j int) bool { return rows[i]["Key"] < rows[j]["Key"] })
		return rows
	},
	"formatJSON": func(s string) string {
		var obj interface{}
		if err := json.Unmarshal([]byte(s), &obj); err != nil {
			return s
		}
		pretty, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			return s
		}
		return string(pretty)
	},
	"derefStatus": func(status *int) string {
		if status == nil {
			return ""
		}
		return strconv.Itoa(*status)
	},
	"statusColor": func(status *int) string {
		if status == nil {
			return "text-gray-500 dark:text-dark-text-muted"
		}
		s := *status
		if s >= 200 && s < 300 {
			return "text-green-600 dark:text-green-400"
		}
		if s >= 400 && s < 500 {
			return "text-yellow-600 dark:text-yellow-400"
		}
		if s >= 500 {
			return "text-red-600 dark:text-red-400"
		}
		return "text-gray-600 dark:text-dark-text-secondary"
	},
	"methodColor": func(method string) string {
		switch method {
		case "GET":
			return "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400"
		case "POST":
			return "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400"
		case "PUT":
			return "bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400"
		case "PATCH":
			return "bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400"
		case "DELETE":
			return "bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400"
		default:
			return "bg-gray-100 text-gray-700 dark:bg-dark-surface-high dark:text-dark-text"
		}
	},
	"formatTransaction": func(baseURL string, r *db.Request) string {
		var b strings.Builder

		fullURL := fmt.Sprintf("%s/%s%s", baseURL, r.EndpointID, r.Path)
		if r.QueryParams != "" {
			fullURL += "?" + r.QueryParams
		}

		b.WriteString(fmt.Sprintf("URL: %s\n", fullURL))
		b.WriteString(fmt.Sprintf("Method: %s\n", r.Method))
		if r.ResponseStatus != nil {
			b.WriteString(fmt.Sprintf("Response: %d\n", *r.ResponseStatus))
		}
		b.WriteString(fmt.Sprintf("Source: %s\n", r.SourceIP))
		b.WriteString(fmt.Sprintf("Received: %s\n", r.ReceivedAt.Format("Mon Jan 02 15:04:05 MST 2006")))
		b.WriteString(fmt.Sprintf("\nRequest size: %dB\n", r.BodySize))
		if r.ResponseBody != "" {
			b.WriteString(fmt.Sprintf("Response size: %dB\n", len(r.ResponseBody)))
		}

		b.WriteString("\n---------- Request ----------\n\n")

		var headers map[string][]string
		if err := json.Unmarshal(r.Headers, &headers); err == nil {
			keys := make([]string, 0, len(headers))
			for k := range headers {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				b.WriteString(fmt.Sprintf("%s: %s\n", k, strings.Join(headers[k], ", ")))
			}
		}

		b.WriteString("\n")
		if r.Body != "" {
			var obj interface{}
			if err := json.Unmarshal([]byte(r.Body), &obj); err == nil {
				pretty, _ := json.MarshalIndent(obj, "", "  ")
				b.WriteString(string(pretty))
			} else {
				b.WriteString(r.Body)
			}
		} else {
			b.WriteString("(body is empty)")
		}

		if r.ResponseStatus != nil {
			b.WriteString("\n\n---------- Response ----------\n\n")
			b.WriteString(fmt.Sprintf("Status: %d\n\n", *r.ResponseStatus))
			if r.ResponseBody != "" {
				var obj interface{}
				if err := json.Unmarshal([]byte(r.ResponseBody), &obj); err == nil {
					pretty, _ := json.MarshalIndent(obj, "", "  ")
					b.WriteString(string(pretty))
				} else {
					b.WriteString(r.ResponseBody)
				}
			} else {
				b.WriteString("(body is empty)")
			}
		}

		return b.String()
	},
}).Parse(`
<div class="bg-white dark:bg-dark-surface rounded-lg shadow-sm border border-gray-200 dark:border-dark-border h-full flex flex-col">
    <div class="flex items-center justify-between px-4 py-2 border-b border-gray-200 dark:border-dark-border">
        <div class="flex items-center gap-3 min-w-0">
            <span class="px-2 py-0.5 rounded text-xs font-semibold flex-shrink-0 {{methodColor .Request.Method}}">{{.Request.Method}}</span>
            <span class="text-xs font-mono text-gray-500 dark:text-dark-text-secondary truncate">{{.BaseURL}}/{{.Request.EndpointID}}{{.Request.Path}}{{if .Request.QueryParams}}?{{.Request.QueryParams}}{{end}}</span>
            {{if .Request.ResponseStatus}}<span class="text-xs font-semibold font-mono flex-shrink-0 {{statusColor .Request.ResponseStatus}}">{{derefStatus .Request.ResponseStatus}}</span>{{end}}
        </div>
        <div class="flex items-center gap-2 flex-shrink-0">
            <button onclick="navigator.clipboard.writeText(document.getElementById('transaction-text').textContent)"
                    class="text-xs text-gray-500 dark:text-dark-text-muted hover:text-gray-700 dark:hover:text-dark-text px-2 py-1 border border-gray-300 dark:border-dark-border rounded transition-colors">
                Copy
            </button>
            <button onclick="downloadTransaction()"
                    class="text-xs text-gray-500 dark:text-dark-text-muted hover:text-gray-700 dark:hover:text-dark-text px-2 py-1 border border-gray-300 dark:border-dark-border rounded transition-colors">
                Download
            </button>
        </div>
    </div>
    <pre id="transaction-text" class="flex-1 overflow-y-auto p-4 text-xs font-mono text-gray-800 dark:text-dark-text whitespace-pre-wrap">{{formatTransaction .BaseURL .Request}}</pre>
</div>
<script>
function downloadTransaction() {
    var text = document.getElementById('transaction-text').textContent;
    var blob = new Blob([text], {type: 'text/plain'});
    var a = document.createElement('a');
    a.href = URL.createObjectURL(blob);
    a.download = 'transaction-{{.Request.Seq}}.txt';
    a.click();
    URL.revokeObjectURL(a.href);
}
</script>
`))

func (h *Handler) PollRequests(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()

	requests, err := h.db.ListRequests(ctx, id, 50)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	count, _ := h.db.GetRequestCount(ctx, id)
	w.Header().Set("X-Request-Count", strconv.Itoa(count))

	requestListTmpl.Execute(w, map[string]interface{}{
		"EndpointID": id,
		"Requests":   requests,
	})
}

func (h *Handler) RequestDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	seqStr := chi.URLParam(r, "reqId")
	seq, err := strconv.Atoi(seqStr)
	if err != nil {
		http.Error(w, "invalid request id", http.StatusBadRequest)
		return
	}

	req, err := h.db.GetRequest(r.Context(), id, seq)
	if err != nil {
		http.Error(w, "request not found", http.StatusNotFound)
		return
	}

	requestDetailTmpl.Execute(w, map[string]interface{}{
		"BaseURL": h.baseURL,
		"Request": req,
	})
}
