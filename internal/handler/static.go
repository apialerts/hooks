package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

var privacyTmpl = template.Must(template.New("privacy").Parse(layoutStart + `
<main class="max-w-2xl mx-auto px-4 py-16">
    <h1 class="text-3xl font-bold dark:text-dark-text mb-8">Privacy Policy</h1>
    <div class="prose text-gray-600 dark:text-dark-text-secondary space-y-4">
        <p>This webhook testing tool is provided free of charge by <a href="https://apialerts.com" class="text-orange-500 hover:underline">API Alerts</a>.</p>
        <h2 class="text-xl font-semibold dark:text-dark-text mt-6">Data We Collect</h2>
        <p>When you create an endpoint, we store:</p>
        <ul class="list-disc pl-6">
            <li>A randomly generated endpoint ID</li>
            <li>Your IP address (for rate limiting)</li>
            <li>Incoming webhook request data (headers, body, metadata)</li>
        </ul>
        <h2 class="text-xl font-semibold dark:text-dark-text mt-6">Data Retention</h2>
        <p>All endpoint data is <strong>automatically deleted after 14 days of inactivity</strong>. There is no way to recover deleted data. We do not create backups of endpoint data.</p>
        <h2 class="text-xl font-semibold dark:text-dark-text mt-6">No Accounts</h2>
        <p>This tool does not require sign-up or authentication. No personal information is collected beyond your IP address for rate limiting purposes.</p>
        <h2 class="text-xl font-semibold dark:text-dark-text mt-6">No Tracking</h2>
        <p>We do not use cookies, analytics, or any third-party tracking on this tool.</p>
    </div>
</main>
` + layoutEnd))

var notFoundTmpl = template.Must(template.New("notfound").Parse(layoutStart + `
<main class="max-w-2xl mx-auto px-4 py-16 text-center">
    <h1 class="text-4xl font-bold dark:text-dark-text mb-4">Endpoint Not Found</h1>
    <p class="text-lg text-gray-600 dark:text-dark-text-secondary mb-8">This endpoint has expired or doesn't exist. Endpoints are automatically deleted after 14 days of inactivity.</p>
    <a href="/" class="bg-orange-500 hover:bg-orange-600 text-white font-semibold py-3 px-8 rounded-lg text-lg transition-colors inline-block">
        Create New Endpoint
    </a>
</main>
` + layoutEnd))

func (h *Handler) Privacy(w http.ResponseWriter, r *http.Request) {
	privacyTmpl.Execute(w, map[string]interface{}{
		"BaseURL": h.baseURL,
	})
}

func (h *Handler) RobotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, `User-agent: *
Allow: /$
Allow: /privacy
Disallow: /

Sitemap: `+h.baseURL+`/sitemap.xml
`)
}

func (h *Handler) SitemapXml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    <url><loc>%s/</loc></url>
    <url><loc>%s/privacy</loc></url>
</urlset>`, h.baseURL, h.baseURL)
}
