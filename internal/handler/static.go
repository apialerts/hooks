package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

var privacyTmpl = template.Must(template.New("privacy").Parse(layoutStart + `
<main class="max-w-2xl mx-auto px-4 sm:px-6 py-16">
    <h1 class="text-3xl font-extrabold dark:text-dark-text mb-2">Privacy Policy</h1>
    <p class="text-sm text-gray-500 dark:text-dark-text-muted mb-10">Last updated: March 2025</p>
    <div class="space-y-8 text-gray-600 dark:text-dark-text-secondary leading-relaxed">
        <p>This webhook testing tool is provided free of charge by <a href="https://apialerts.com" class="text-brand hover:underline font-medium">API Alerts</a>.</p>

        <div>
            <h2 class="text-lg font-bold dark:text-dark-text mb-3">Data We Collect</h2>
            <p class="mb-3">When you create an endpoint, we store:</p>
            <ul class="list-disc pl-5 space-y-1.5 text-sm">
                <li>A randomly generated endpoint ID</li>
                <li>Your IP address (for rate limiting only)</li>
                <li>Incoming webhook request data (headers, body, metadata)</li>
            </ul>
        </div>

        <div>
            <h2 class="text-lg font-bold dark:text-dark-text mb-3">Data Retention</h2>
            <p>All endpoint data is <strong class="dark:text-dark-text">automatically deleted after 7 days of inactivity</strong>. There is no way to recover deleted data. We do not create backups of endpoint data.</p>
        </div>

        <div>
            <h2 class="text-lg font-bold dark:text-dark-text mb-3">No Accounts</h2>
            <p>This tool does not require sign-up or authentication. No personal information is collected beyond your IP address for rate limiting purposes.</p>
        </div>

        <div>
            <h2 class="text-lg font-bold dark:text-dark-text mb-3">No Tracking</h2>
            <p>We do not use cookies, analytics, or any third-party tracking on this tool.</p>
        </div>
    </div>
</main>
` + layoutEnd))

var notFoundTmpl = template.Must(template.New("notfound").Parse(layoutStart + `
<main class="max-w-lg mx-auto px-4 sm:px-6 py-24 text-center">
    <div class="w-14 h-14 rounded-full bg-gray-100 dark:bg-dark-surface-high flex items-center justify-center mx-auto mb-6">
        <svg class="w-7 h-7 text-gray-400 dark:text-dark-text-muted" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
    </div>
    <h1 class="text-2xl font-extrabold dark:text-dark-text mb-3">Endpoint Not Found</h1>
    <p class="text-gray-600 dark:text-dark-text-secondary mb-8 leading-relaxed">This endpoint has expired or doesn't exist. Endpoints are automatically deleted after 7 days of inactivity.</p>
    <a href="/" class="inline-flex items-center gap-2 bg-brand hover:brightness-110 text-black font-semibold py-3 px-8 rounded-full text-base transition-all shadow-sm hover:shadow-md">
        Create New Endpoint
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M13 7l5 5m0 0l-5 5m5-5H6"/></svg>
    </a>
</main>
` + layoutEnd))

func (h *Handler) Privacy(w http.ResponseWriter, r *http.Request) {
	render(w, privacyTmpl, map[string]interface{}{
		"BaseURL":         h.baseURL,
		"PageTitle":       "Privacy Policy",
		"PageDescription": "Privacy policy for hooks.apialerts.com. No cookies, no analytics, no tracking. All data auto-deleted after 7 days of inactivity.",
		"CanonicalURL":    h.baseURL + "/privacy",
	})
}

func (h *Handler) RobotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, `User-agent: *
Allow: /$
Allow: /privacy
Allow: /sitemap.xml
Disallow: /

Sitemap: `+h.baseURL+`/sitemap.xml
`)
}

func (h *Handler) SitemapXml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    <url>
        <loc>%s/</loc>
        <changefreq>weekly</changefreq>
        <priority>1.0</priority>
    </url>
    <url>
        <loc>%s/privacy</loc>
        <changefreq>monthly</changefreq>
        <priority>0.3</priority>
    </url>
</urlset>`, h.baseURL, h.baseURL)
}
