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
	"isoTime": func(t time.Time) string {
		return t.UTC().Format(time.RFC3339)
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
			return "bg-green-500/10 text-green-600 dark:text-green-400"
		case "POST":
			return "bg-blue-500/10 text-blue-600 dark:text-blue-400"
		case "PUT":
			return "bg-yellow-500/10 text-yellow-600 dark:text-yellow-400"
		case "PATCH":
			return "bg-orange-500/10 text-orange-600 dark:text-orange-400"
		case "DELETE":
			return "bg-red-500/10 text-red-600 dark:text-red-400"
		default:
			return "bg-gray-500/10 text-gray-600 dark:text-dark-text-secondary"
		}
	},
}).Parse(`
{{if .Requests}}
<div class="divide-y divide-gray-100 dark:divide-dark-border">
    {{range .Requests}}
    <div class="request-item flex items-center gap-3 px-3 py-2.5 cursor-pointer hover:bg-gray-50 dark:hover:bg-dark-surface-high transition-colors group"
         data-seq="{{.Seq}}"
         hx-get="/{{$.EndpointSlug}}/requests/{{.Seq}}"
         hx-target="#request-detail"
         hx-swap="innerHTML"
         onclick="if(this.style.borderLeft){return false}document.querySelectorAll('.request-item').forEach(function(el){el.style.borderLeft='';el.style.backgroundColor=''});this.style.borderLeft='2px solid #FF7A00';this.style.backgroundColor='rgba(255,122,0,0.08)'">
        <span class="px-2 py-0.5 rounded-md text-[11px] font-bold w-14 text-center tracking-wide {{methodColor .Method}}">{{.Method}}</span>
        <time class="text-xs font-mono text-gray-400 dark:text-dark-text-muted local-time" datetime="{{isoTime .ReceivedAt}}">{{isoTime .ReceivedAt}}</time>
        {{if .ResponseStatus}}<span class="text-xs font-mono font-bold {{statusColor .ResponseStatus}}">{{derefStatus .ResponseStatus}}</span>{{end}}
        <span class="text-[11px] text-gray-400 dark:text-dark-text-muted ml-auto">{{formatSize .BodySize}}</span>
        <svg class="w-3.5 h-3.5 text-gray-300 dark:text-dark-border group-hover:text-gray-400 dark:group-hover:text-dark-text-muted transition-colors flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7"/></svg>
    </div>
    {{end}}
</div>
{{else}}
<div class="h-full flex flex-col items-center justify-center p-10 text-center">
    <svg class="w-8 h-8 text-gray-300 dark:text-dark-border mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"/></svg>
    <p class="text-sm font-medium text-gray-400 dark:text-dark-text-muted mb-1">No requests yet</p>
    <p class="text-xs text-gray-400 dark:text-dark-text-muted">Send a webhook to your endpoint URL</p>
</div>
{{end}}
<script>
document.querySelectorAll('.local-time').forEach(function(el) {
    var d = new Date(el.getAttribute('datetime'));
    el.textContent = d.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit', second: '2-digit'});
});
</script>
`))

var requestDetailTmpl = template.Must(template.New("request_detail").Funcs(template.FuncMap{
	"isoTime": func(t time.Time) string {
		return t.UTC().Format(time.RFC3339)
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
			return "bg-green-500/10 text-green-600 dark:text-green-400"
		case "POST":
			return "bg-blue-500/10 text-blue-600 dark:text-blue-400"
		case "PUT":
			return "bg-yellow-500/10 text-yellow-600 dark:text-yellow-400"
		case "PATCH":
			return "bg-orange-500/10 text-orange-600 dark:text-orange-400"
		case "DELETE":
			return "bg-red-500/10 text-red-600 dark:text-red-400"
		default:
			return "bg-gray-500/10 text-gray-600 dark:text-dark-text-secondary"
		}
	},
	"formatTransaction": func(baseURL, slug string, r *db.Request) string {
		var b strings.Builder

		fullURL := fmt.Sprintf("%s/%s%s", baseURL, slug, r.Path)
		if r.QueryParams != "" {
			fullURL += "?" + r.QueryParams
		}

		b.WriteString(fmt.Sprintf("URL: %s\n", fullURL))
		b.WriteString(fmt.Sprintf("Method: %s\n", r.Method))
		if r.ResponseStatus != nil {
			b.WriteString(fmt.Sprintf("Response: %d\n", *r.ResponseStatus))
		}
		b.WriteString(fmt.Sprintf("Source: %s\n", r.SourceIP))
		b.WriteString(fmt.Sprintf("Received: %s\n", r.ReceivedAt.UTC().Format(time.RFC3339)))
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
<div class="rounded-xl bg-white dark:bg-dark-surface border border-gray-200 dark:border-dark-border h-full flex flex-col">
    <div class="flex items-center justify-between px-4 py-2.5 border-b border-gray-100 dark:border-dark-border">
        <div class="flex items-center gap-2.5 min-w-0">
            <span class="px-2 py-0.5 rounded-md text-[11px] font-bold flex-shrink-0 tracking-wide {{methodColor .Request.Method}}">{{.Request.Method}}</span>
            <span class="text-xs font-mono text-gray-500 dark:text-dark-text-secondary truncate">/{{.EndpointSlug}}{{.Request.Path}}{{if .Request.QueryParams}}?{{.Request.QueryParams}}{{end}}</span>
            {{if .Request.ResponseStatus}}<span class="text-xs font-bold font-mono flex-shrink-0 {{statusColor .Request.ResponseStatus}}">{{derefStatus .Request.ResponseStatus}}</span>{{end}}
        </div>
        <div class="flex items-center gap-1 flex-shrink-0">
            <button onclick="var b=this;navigator.clipboard.writeText(document.getElementById('transaction-text').textContent);var o=b.innerHTML;b.innerHTML='<svg class=\'w-3 h-3 text-green-500\' fill=\'none\' viewBox=\'0 0 24 24\' stroke=\'currentColor\' stroke-width=\'2\'><path stroke-linecap=\'round\' stroke-linejoin=\'round\' d=\'M5 13l4 4L19 7\'/></svg><span>Copied</span>';setTimeout(function(){b.innerHTML=o},1500)"
                    class="inline-flex items-center gap-1.5 text-xs text-gray-500 dark:text-dark-text-muted hover:text-gray-700 dark:hover:text-dark-text px-2.5 py-1.5 border border-gray-200 dark:border-dark-border rounded-lg transition-colors hover:bg-gray-50 dark:hover:bg-dark-surface-high font-medium">
                <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3"/></svg>
                <span>Copy</span>
            </button>
            <button onclick="downloadTransaction()"
                    class="inline-flex items-center gap-1.5 text-xs text-gray-500 dark:text-dark-text-muted hover:text-gray-700 dark:hover:text-dark-text px-2.5 py-1.5 border border-gray-200 dark:border-dark-border rounded-lg transition-colors hover:bg-gray-50 dark:hover:bg-dark-surface-high font-medium">
                <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/></svg>
                Download
            </button>
        </div>
    </div>
    <pre id="transaction-text" class="flex-1 overflow-y-auto p-4 text-xs font-mono text-gray-700 dark:text-dark-text-secondary whitespace-pre-wrap leading-relaxed">{{formatTransaction .BaseURL .EndpointSlug .Request}}</pre>
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
(function() {
    var pre = document.getElementById('transaction-text');
    if (!pre) return;
    pre.textContent = pre.textContent.replace(/Received: (\d{4}-\d{2}-\d{2}T[^\n]+)/, function(_, iso) {
        return 'Received: ' + new Date(iso).toLocaleString();
    });
})();
</script>
`))

func (h *Handler) PollRequests(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "id")
	ctx := r.Context()

	endpoint, err := h.db.GetEndpoint(ctx, slug)
	if err != nil {
		http.Error(w, "endpoint not found", http.StatusNotFound)
		return
	}

	requests, err := h.db.ListRequests(ctx, endpoint.ID, 50)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-Request-Count", strconv.Itoa(endpoint.RequestCount))

	render(w, requestListTmpl, map[string]interface{}{
		"EndpointSlug": slug,
		"Requests":     requests,
	})
}

func (h *Handler) RequestDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "id")
	seqStr := chi.URLParam(r, "reqId")
	seq, err := strconv.Atoi(seqStr)
	if err != nil {
		http.Error(w, "invalid request id", http.StatusBadRequest)
		return
	}

	endpoint, err := h.db.GetEndpoint(r.Context(), slug)
	if err != nil {
		http.Error(w, "endpoint not found", http.StatusNotFound)
		return
	}

	req, err := h.db.GetRequest(r.Context(), endpoint.ID, seq)
	if err != nil {
		http.Error(w, "request not found", http.StatusNotFound)
		return
	}

	render(w, requestDetailTmpl, map[string]interface{}{
		"BaseURL":      h.baseURL,
		"EndpointSlug": slug,
		"Request":      req,
	})
}
