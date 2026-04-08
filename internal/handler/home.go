package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

var homeTmpl = template.Must(template.New("home").Parse(layoutStart + `
<script type="application/ld+json">
{
    "@context": "https://schema.org",
    "@type": "WebApplication",
    "name": "Webhook Tester",
    "url": "https://hooks.apialerts.com",
    "description": "Free webhook testing tool. Generate a unique URL, inspect incoming requests, and toggle between success, error, and timeout responses to test retry logic.",
    "applicationCategory": "DeveloperApplication",
    "operatingSystem": "Any",
    "offers": {
        "@type": "Offer",
        "price": "0",
        "priceCurrency": "USD"
    },
    "author": {
        "@type": "Organization",
        "name": "API Alerts",
        "url": "https://apialerts.com"
    }
}
</script>
<main class="max-w-2xl mx-auto px-4 sm:px-6 py-16 sm:py-24">
    <h1 class="text-3xl sm:text-4xl font-extrabold tracking-tight dark:text-dark-text mb-4">Webhook Tester</h1>
    <p class="text-gray-600 dark:text-dark-text-secondary leading-relaxed mb-8">
        Get a URL. Send webhooks to it. See what arrives. Control what comes back.<br>
        Test retry logic by toggling between 2xx, 4xx, 5xx, and timeouts.
    </p>

    <form method="POST" action="/endpoints" class="mb-16">
        <button type="submit" class="bg-brand hover:brightness-110 text-black font-semibold py-2.5 px-6 rounded-full text-sm transition-all">
            Create Endpoint
        </button>
        <span class="ml-3 text-xs text-gray-400 dark:text-dark-text-muted">Free, no sign-up. All data deleted after 7 days.</span>
    </form>

    <!-- curl example -->
    <div class="mb-16">
        <pre class="text-sm font-mono bg-white dark:bg-dark-surface border border-gray-200 dark:border-dark-border rounded-xl p-5 overflow-x-auto text-gray-700 dark:text-dark-text-secondary leading-relaxed"><span class="text-gray-400 dark:text-dark-text-muted select-none">$ </span><span class="dark:text-dark-text text-gray-900">curl -X POST https://hooks.apialerts.com/your-endpoint \
  -H "Content-Type: application/json" \
  -d '{"id": 42, "event": "order.completed"}'</span>

<span class="text-gray-400 dark:text-dark-text-muted"># Response depends on your configured mode:</span>
<span class="text-green-600 dark:text-green-400">200</span> {"status": "ok"}
<span class="text-red-600 dark:text-red-400">500</span> {"error": "internal_error", "message": "Something went wrong"}
<span class="text-yellow-600 dark:text-yellow-400">408</span> <span class="text-gray-400 dark:text-dark-text-muted">(30s timeout &mdash; no response)</span></pre>
    </div>

    <!-- Why -->
    <div class="mb-16">
        <h2 class="text-lg font-bold dark:text-dark-text mb-3">Why this exists</h2>
        <p class="text-gray-600 dark:text-dark-text-secondary leading-relaxed mb-3">
            We built this to test webhook delivery for <a href="https://apialerts.com" target="_blank" class="text-brand hover:underline">API Alerts, a notification delivery platform</a>. We needed to know: does our sender actually retry on 5xx? Does it give up on 4xx? Does it handle a 30-second timeout without crashing?
        </p>
        <p class="text-gray-600 dark:text-dark-text-secondary leading-relaxed">
            Existing tools show the request. We needed to control the response. So we built Hooks, and then open-sourced it.
        </p>
    </div>

    <!-- What you can do -->
    <div class="mb-16">
        <h2 class="text-lg font-bold dark:text-dark-text mb-3">What you can do</h2>
        <ul class="space-y-2 text-gray-600 dark:text-dark-text-secondary">
            <li class="flex gap-3 leading-relaxed"><span class="text-brand font-bold select-none">&rsaquo;</span> Inspect headers, body, query params, and source IP for every request</li>
            <li class="flex gap-3 leading-relaxed"><span class="text-brand font-bold select-none">&rsaquo;</span> Toggle the response: 200, 201, 400, 401, 403, 404, 500, 503, or a 30s timeout</li>
            <li class="flex gap-3 leading-relaxed"><span class="text-brand font-bold select-none">&rsaquo;</span> Edit the response body &mdash; each preset has a sensible JSON default</li>
            <li class="flex gap-3 leading-relaxed"><span class="text-brand font-bold select-none">&rsaquo;</span> Send a test request from the browser without leaving the page</li>
            <li class="flex gap-3 leading-relaxed"><span class="text-brand font-bold select-none">&rsaquo;</span> Copy or download the full request/response transaction as plain text</li>
        </ul>
    </div>

    <!-- Details -->
    <div class="mb-16 grid grid-cols-1 sm:grid-cols-2 gap-x-12 gap-y-8 text-sm">
        <div>
            <h3 class="font-bold dark:text-dark-text mb-2">Limits</h3>
            <ul class="space-y-1 text-gray-500 dark:text-dark-text-muted">
                <li>5 endpoints per IP</li>
                <li>50 requests stored per endpoint</li>
                <li>60 requests/min rate limit</li>
                <li>256KB max payload</li>
                <li>7-day expiry from last activity</li>
            </ul>
        </div>
        <div>
            <h3 class="font-bold dark:text-dark-text mb-2">Privacy</h3>
            <ul class="space-y-1 text-gray-500 dark:text-dark-text-muted">
                <li>No sign-up or accounts</li>
                <li>No cookies or analytics</li>
                <li>No third-party tracking</li>
                <li>Data auto-deleted after 7 days</li>
                <li>We don't back up endpoint data</li>
            </ul>
        </div>
    </div>

    <!-- Open source -->
    <div class="text-sm text-gray-500 dark:text-dark-text-muted border-t border-gray-200 dark:border-dark-border pt-8 space-y-3">
        <p>
            MIT licensed. Single Go binary, Postgres, HTMX. No JavaScript frameworks.
            <a href="https://github.com/apialerts/hooks" target="_blank" class="text-brand hover:underline">View on GitHub</a> or self-host with <code class="text-xs bg-gray-100 dark:bg-dark-surface-high px-1.5 py-0.5 rounded">docker compose up</code>.
        </p>
        <p>
            Need to send notifications to Slack, Discord, push, and more from a single API? Check out <a href="https://apialerts.com" target="_blank" class="text-brand hover:underline">API Alerts</a>.
        </p>
    </div>
</main>
` + layoutEnd))

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	homeTmpl.Execute(w, map[string]interface{}{
		"BaseURL":         h.baseURL,
		"PageTitle":       "Free Webhook Tester",
		"PageDescription": "Free webhook testing tool. Generate a unique URL, inspect incoming requests, and toggle between success, error, and timeout responses to test your retry logic. Open source, no sign-up required.",
		"CanonicalURL":    h.baseURL + "/",
	})
}

var limitTmpl = template.Must(template.New("limit").Parse(layoutStart + `
<main class="max-w-lg mx-auto px-4 sm:px-6 py-24">
    <h1 class="text-2xl font-extrabold dark:text-dark-text mb-3">Endpoint limit reached</h1>
    <p class="text-gray-600 dark:text-dark-text-secondary mb-6 leading-relaxed">You already have 5 endpoints, which is the maximum per IP address. Delete one to create a new one.</p>
    <div class="space-y-2">
        {{range .Endpoints}}
        <a href="/{{.ID}}" class="flex items-center justify-between px-4 py-3 rounded-xl bg-white dark:bg-dark-surface border border-gray-200 dark:border-dark-border hover:border-gray-300 dark:hover:border-dark-text-muted transition-colors group">
            <code class="text-sm font-mono text-gray-700 dark:text-dark-text">{{.ID}}</code>
            <svg class="w-4 h-4 text-gray-300 dark:text-dark-border group-hover:text-gray-400 dark:group-hover:text-dark-text-muted transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7"/></svg>
        </a>
        {{end}}
    </div>
</main>
` + layoutEnd))

func (h *Handler) CreateEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ip := r.RemoteAddr

	count, err := h.db.CountEndpointsByIP(ctx, ip)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if count >= 5 {
		endpoints, _ := h.db.ListEndpointsByIP(ctx, ip)
		w.WriteHeader(http.StatusTooManyRequests)
		limitTmpl.Execute(w, map[string]interface{}{
			"BaseURL":         h.baseURL,
			"Endpoints":       endpoints,
			"PageTitle":       "Endpoint Limit Reached",
			"PageDescription": "",
			"CanonicalURL":    h.baseURL + "/",
		})
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
