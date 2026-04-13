module.exports = {
  content: ["./templates/**/*.templ"],
  theme: {
    extend: {
      colors: {
        bg: { DEFAULT: '#0a0a0f', surface: '#141420', border: '#1e1e2e' },
        accent: { DEFAULT: '#f59e0b', light: '#fbbf24', subtle: 'rgba(245, 158, 11, 0.15)' },
        'accent-2': { DEFAULT: '#ec4899', light: '#f472b6', subtle: 'rgba(236, 72, 153, 0.15)' },
        text: { DEFAULT: '#f0f0f5', muted: '#9090a0' },
      },
      fontFamily: {
        heading: ['"Space Grotesk"', 'sans-serif'],
        body: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['"JetBrains Mono"', 'monospace'],
      }
    }
  },
  plugins: [],
}
