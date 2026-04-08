package handler

const layoutStart = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Webhook Tester - hooks.apialerts.com</title>
    <meta name="description" content="Free webhook testing tool. Generate a unique URL, inspect requests, and toggle between success and failure responses. Built by API Alerts.">
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
<body class="bg-gray-50 dark:bg-dark-bg min-h-screen flex flex-col font-sans text-gray-900 dark:text-dark-text antialiased">
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
                <span class="text-xs font-semibold text-gray-500 dark:text-dark-text-muted">A free tool by <a href="https://apialerts.com" target="_blank" class="text-brand hover:underline">API Alerts</a></span>
                <div class="flex items-center gap-4 text-xs text-gray-400 dark:text-dark-text-muted">
                    <a href="/privacy" class="hover:text-gray-600 dark:hover:text-dark-text-secondary transition-colors">Privacy</a>
                    <a href="https://github.com/apialerts/hooks" target="_blank" class="hover:text-gray-600 dark:hover:text-dark-text-secondary transition-colors">GitHub</a>
                    <a href="https://apialerts.com" target="_blank" class="hover:text-gray-600 dark:hover:text-dark-text-secondary transition-colors">apialerts.com</a>
                </div>
            </div>
        </div>
    </footer>
</body>
</html>`
