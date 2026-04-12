module.exports = {
  content: ["./templates/**/*.templ"],
  theme: {
    extend: {
      colors: {
        bg: { DEFAULT: '#0a0a0f', surface: '#141420', border: '#1e1e2e' },
        accent: { DEFAULT: '#7c3aed', light: '#a78bfa' },
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
