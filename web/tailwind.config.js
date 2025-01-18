/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["../internal/view/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [
    require('daisyui'),
  ],
  daisyui: {
    themes: ["dark", "light", "retro", "cyberpunk", "valentine", "aqua", "business"],
  },
}

