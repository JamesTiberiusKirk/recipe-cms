import plugin from 'tailwindcss/plugin';

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.html", "./**/*.templ", "./**/*.go",],
  darkMode: 'class',
  theme: {
    extend: {
    },

    screens: {
      'phone': '450px',
      'tablet': '640px',
      'laptop': '1024px',
      'desktop': '1280px',
    },
  },
  plugins: [
    plugin(function ({addVariant}) {
        addVariant('progress-unfilled', ['&::-webkit-progress-bar', '&']);
        addVariant('progress-filled', ['&::-webkit-progress-value', '&::-moz-progress-bar', '&']);
        addVariant('progress', [ '&::-webkit-progress-bar', '&::-webkit-progress-value', '&::-moz-progress-bar', '&']);
    })
 ],
}

