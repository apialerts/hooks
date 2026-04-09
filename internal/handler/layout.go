package handler

const layoutStart = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{if .PageTitle}}{{.PageTitle}} - {{end}}hooks.apialerts.com</title>
    <meta name="description" content="{{if .PageDescription}}{{.PageDescription}}{{else}}Free webhook testing tool. Generate a unique URL, inspect requests, and toggle between success and failure responses to test retry logic. No sign-up required.{{end}}">
    <meta name="author" content="API Alerts">
    <meta name="keywords" content="webhook tester, test webhooks, webhook debugger, webhook testing tool, mock webhook endpoint, webhook receiver, simulate webhook failures, webhook inspector, free webhook tool">
    <link rel="canonical" href="{{.CanonicalURL}}">

    <!-- Open Graph -->
    <meta property="og:type" content="website">
    <meta property="og:site_name" content="hooks.apialerts.com">
    <meta property="og:title" content="{{if .PageTitle}}{{.PageTitle}} - {{end}}hooks.apialerts.com">
    <meta property="og:description" content="{{if .PageDescription}}{{.PageDescription}}{{else}}Free webhook testing tool. Generate a unique URL, inspect requests, and toggle between success and failure responses to test retry logic. No sign-up required.{{end}}">
    <meta property="og:url" content="{{.CanonicalURL}}">
    <meta property="og:image" content="{{.BaseURL}}/static/meta.png">
    <meta property="og:image:width" content="1200">
    <meta property="og:image:height" content="630">

    <!-- Twitter -->
    <meta name="twitter:card" content="summary_large_image">
    <meta name="twitter:image" content="{{.BaseURL}}/static/meta.png">
    <meta name="twitter:title" content="{{if .PageTitle}}{{.PageTitle}} - {{end}}hooks.apialerts.com">
    <meta name="twitter:description" content="{{if .PageDescription}}{{.PageDescription}}{{else}}Free webhook testing tool. Generate a unique URL, inspect requests, and toggle between success and failure responses to test retry logic. No sign-up required.{{end}}">

    <link rel="icon" href="/static/favicon.ico">
    <link rel="icon" href="/static/favicon-dark.svg" type="image/svg+xml" media="(prefers-color-scheme: dark)">
    <link rel="icon" href="/static/favicon-light.svg" type="image/svg+xml" media="(prefers-color-scheme: light)">
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Mulish:wght@400;500;600;700;800&display=swap" rel="stylesheet">
    <script src="/static/htmx.min.js"></script>
    <script src="https://cdn.tailwindcss.com"></script>
    <script>
    tailwind.config = {
        darkMode: 'class',
        theme: {
            extend: {
                fontFamily: {
                    sans: ['Mulish', 'system-ui', 'sans-serif'],
                },
                colors: {
                    brand: '#e8772e',
                    dark: {
                        bg: '#1c1c1c',
                        surface: '#252525',
                        'surface-high': '#2e2e2e',
                        'surface-highest': '#3a3a3a',
                        border: '#383838',
                        text: '#e5e5e5',
                        'text-secondary': '#a0a0a0',
                        'text-muted': '#6e6e6e',
                    }
                }
            }
        }
    }
    </script>
    <script>
        const scheme = localStorage.getItem('scheme') ?? 'dark';
        if (scheme === 'dark') document.documentElement.classList.add('dark');
    </script>
    <link rel="stylesheet" href="/static/styles.css">
</head>
<body class="bg-gray-50 dark:bg-dark-bg min-h-screen flex flex-col font-sans text-gray-900 dark:text-dark-text antialiased overflow-x-hidden">
    <header class="border-b border-gray-200 dark:border-dark-border bg-white/80 dark:bg-dark-surface/80 backdrop-blur-sm sticky top-0 z-50">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 h-14 flex items-center justify-between">
            <a href="/" class="text-xl font-bold tracking-tight text-gray-900 dark:text-dark-text"><span class="text-brand">hooks.</span>apialerts.com</a>
            <div class="flex items-center gap-1">
                <a href="https://github.com/apialerts/hooks" target="_blank" aria-label="GitHub"
                   class="p-2 rounded-full text-gray-400 dark:text-dark-text-muted hover:text-gray-600 dark:hover:text-dark-text-secondary hover:bg-gray-100 dark:hover:bg-dark-surface-high transition-colors">
                    <svg class="w-[18px] h-[18px]" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"/></svg>
                </a>
                <button id="theme-toggle" aria-label="Toggle theme"
                        class="p-2 rounded-full text-gray-400 dark:text-dark-text-muted hover:text-gray-600 dark:hover:text-dark-text-secondary hover:bg-gray-100 dark:hover:bg-dark-surface-high transition-colors">
                    <svg id="icon-moon" xmlns="http://www.w3.org/2000/svg" class="w-[18px] h-[18px]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z" />
                    </svg>
                    <svg id="icon-sun" xmlns="http://www.w3.org/2000/svg" class="w-[18px] h-[18px] hidden" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364-6.364l-.707.707M6.343 17.657l-.707.707M17.657 17.657l-.707-.707M6.343 6.343l-.707-.707M12 7a5 5 0 100 10A5 5 0 0012 7z" />
                    </svg>
                </button>
            </div>
        </div>
    </header>
    <script>
        (function() {
            var btn = document.getElementById('theme-toggle');
            var moon = document.getElementById('icon-moon');
            var sun = document.getElementById('icon-sun');
            function apply(scheme) {
                if (scheme === 'dark') {
                    document.documentElement.classList.add('dark');
                    moon.classList.remove('hidden');
                    sun.classList.add('hidden');
                } else {
                    document.documentElement.classList.remove('dark');
                    sun.classList.remove('hidden');
                    moon.classList.add('hidden');
                }
                localStorage.setItem('scheme', scheme);
            }
            apply(localStorage.getItem('scheme') ?? 'dark');
            btn.addEventListener('click', function() {
                var isDark = document.documentElement.classList.contains('dark');
                apply(isDark ? 'light' : 'dark');
            });
        })();
    </script>
`

const layoutEnd = `
    <footer class="mt-auto border-t border-gray-200 dark:border-dark-border">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 py-8">
            <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
                <div class="flex items-center gap-4">
                    <span class="text-xs font-semibold text-gray-500 dark:text-dark-text-muted">A free tool by <a href="https://apialerts.com" target="_blank" class="text-brand hover:underline">API Alerts</a></span>
                    <div class="flex items-center gap-2">
                        <a href="https://github.com/apialerts" target="_blank" rel="noopener" aria-label="GitHub" class="text-gray-400 dark:text-dark-text-muted hover:text-gray-600 dark:hover:text-dark-text-secondary transition-colors">
                            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"/></svg>
                        </a>
                        <a href="https://x.com/api_alerts" target="_blank" rel="noopener" aria-label="X" class="text-gray-400 dark:text-dark-text-muted hover:text-gray-600 dark:hover:text-dark-text-secondary transition-colors">
                            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M18.901 1.153h3.68l-8.04 9.19L24 22.846h-7.406l-5.8-7.584-6.638 7.584H.474l8.6-9.83L0 1.154h7.594l5.243 6.932ZM17.61 20.644h2.039L6.486 3.24H4.298Z"/></svg>
                        </a>
                        <a href="https://www.linkedin.com/company/apialerts" target="_blank" rel="noopener" aria-label="LinkedIn" class="text-gray-400 dark:text-dark-text-muted hover:text-gray-600 dark:hover:text-dark-text-secondary transition-colors">
                            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433a2.062 2.062 0 01-2.063-2.065 2.064 2.064 0 112.063 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>
                        </a>
                        <a href="https://dev.to/apialerts" target="_blank" rel="noopener" aria-label="DEV" class="text-gray-400 dark:text-dark-text-muted hover:text-gray-600 dark:hover:text-dark-text-secondary transition-colors">
                            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M7.42 10.05c-.18-.16-.46-.23-.84-.23H6v4.36h.58c.37 0 .65-.08.84-.23.2-.16.3-.46.3-.9v-2.1c0-.44-.1-.74-.3-.9zm13.37-7.03H3.21C1.72 3.02.46 4.25.43 5.74v12.52c.03 1.49 1.29 2.72 2.78 2.72h17.58c1.49 0 2.75-1.23 2.78-2.72V5.74c-.03-1.49-1.29-2.72-2.78-2.72zm-9.85 10.9c0 1.2-.62 2.05-2.08 2.05H6.59V8.03h2.27c1.45 0 2.08.85 2.08 2.05v3.84zm4.56 2.05H13V8.03h2.5v1.2h-1.3v2.07h1.3v1.2h-1.3v2.27h1.3v1.2zm4.87-4.45c0 .63-.08 1.08-.25 1.36-.24.39-.62.58-1.15.58-.36 0-.65-.13-.89-.38l-.06.32H17.2V8.03h1.3v2.54c.24-.27.53-.4.88-.4.53 0 .91.2 1.14.58.17.28.25.73.25 1.36v1.41zm-4.32-2.09v.01h-.01l.01-.01z"/></svg>
                        </a>
                    </div>
                </div>
                <div class="flex items-center gap-4 text-xs text-gray-400 dark:text-dark-text-muted">
                    <a href="/privacy" class="hover:text-gray-600 dark:hover:text-dark-text-secondary transition-colors">Privacy</a>
                    <a href="https://github.com/apialerts/hooks" target="_blank" class="hover:text-gray-600 dark:hover:text-dark-text-secondary transition-colors">Source</a>
                    <a href="https://apialerts.com" target="_blank" class="hover:text-gray-600 dark:hover:text-dark-text-secondary transition-colors">apialerts.com</a>
                </div>
            </div>
        </div>
    </footer>
</body>
</html>`
