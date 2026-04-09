/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./internal/handler/**/*.go'],
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
        },
      },
    },
  },
  plugins: [],
}
