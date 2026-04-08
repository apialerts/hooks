package handler

const layoutStart = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Webhook Tester - hooks.apialerts.com</title>
    <meta name="description" content="Free webhook testing tool. Generate a unique URL, inspect requests, and toggle between success and failure responses. Built by API Alerts.">
    <link rel="icon" href="/static/favicon.ico">
    <link rel="icon" href="/static/favicon.svg" type="image/svg+xml">
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Mulish:wght@400;500;600;700&display=swap" rel="stylesheet">
    <script src="/static/htmx.min.js"></script>
    <script src="https://cdn.tailwindcss.com"></script>
    <script>
    tailwind.config = {
        darkMode: 'class',
        theme: {
            extend: {
                colors: {
                    dark: {
                        bg: '#1a1a1a',
                        surface: '#232323',
                        'surface-low': '#1e1e1e',
                        'surface-container': '#272727',
                        'surface-high': '#2e2e2e',
                        'surface-highest': '#3a3a3a',
                        border: '#4a4641',
                        'border-subtle': '#333330',
                        text: '#e3e3e3',
                        'text-secondary': '#a8a8a8',
                        'text-muted': '#777777',
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
<body class="bg-gray-50 dark:bg-dark-bg min-h-screen flex flex-col font-['Mulish',sans-serif] text-gray-900 dark:text-dark-text transition-colors">
    <header class="bg-white dark:bg-dark-surface border-b border-gray-200 dark:border-dark-border">
        <div class="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
            <a href="/" class="text-lg font-bold text-gray-900 dark:text-dark-text">hooks<span class="text-orange-500">.apialerts.com</span></a>
            <div class="flex items-center gap-4">
                <button id="theme-toggle" aria-label="Toggle theme" class="p-2 rounded-lg text-gray-500 dark:text-dark-text-muted hover:bg-gray-100 dark:hover:bg-dark-surface-high transition-colors">
                    <svg id="icon-moon" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z" />
                    </svg>
                    <svg id="icon-sun" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 hidden" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364-6.364l-.707.707M6.343 17.657l-.707.707M17.657 17.657l-.707-.707M6.343 6.343l-.707-.707M12 7a5 5 0 100 10A5 5 0 0012 7z" />
                    </svg>
                </button>
                <a href="https://apialerts.com" target="_blank" class="text-sm text-gray-500 dark:text-dark-text-muted hover:text-orange-500 transition-colors">
                    Built by API Alerts
                </a>
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
    <footer class="mt-auto bg-white dark:bg-dark-surface border-t border-gray-200 dark:border-dark-border">
        <div class="max-w-7xl mx-auto px-4 py-6 text-center">
            <p class="text-gray-600 dark:text-dark-text-secondary mb-2">Need webhook delivery with retries, templates, and multi-channel routing?</p>
            <a href="https://apialerts.com" target="_blank" class="text-orange-500 hover:text-orange-600 font-semibold">
                Try API Alerts &rarr;
            </a>
            <div class="mt-4 text-xs text-gray-400 dark:text-dark-text-muted">
                <a href="/privacy" class="hover:text-gray-600 dark:hover:text-dark-text-secondary">Privacy</a>
                <span class="mx-2">&middot;</span>
                <a href="https://github.com/apialerts/hooks" target="_blank" class="hover:text-gray-600 dark:hover:text-dark-text-secondary">GitHub</a>
            </div>
        </div>
    </footer>
</body>
</html>`
