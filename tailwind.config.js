import plugin from 'tailwindcss/plugin';

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.html", "./**/*.templ", "./**/*.go",],
  darkMode: 'class',
  theme: {
    extend: {
    },

    // screens: {
      // 'phone': {'min': '0px', 'max': '660px'},
      // 'tablet': {'min': '661px', 'max': '1023px'},
      // 'laptop': {'min': '1024px', 'max': '1279px'},
      // 'dektop': {'min': '1280px'},
      //
      // 'tablet<': {'min': '661px'},
      // 'laptop<': {'min': '1024px'},
    // },
  },
  plugins: [
    plugin(function ({addVariant}) {
        addVariant('progress-unfilled', ['&::-webkit-progress-bar', '&']);
        addVariant('progress-filled', ['&::-webkit-progress-value', '&::-moz-progress-bar', '&']);
        addVariant('progress', [ '&::-webkit-progress-bar', '&::-webkit-progress-value', '&::-moz-progress-bar', '&']);
    })
 ],
}

