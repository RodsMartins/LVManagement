const colors = require('tailwindcss/colors')

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    'internal/templates/*.templ',
    'internal/templates/*/*.templ',
    'internal/templates/*/*/*.templ',
  ],
  theme: {
    container: {
      center: true,
      padding: {
        DEFAULT: "1rem",
        mobile: "2rem",
        tablet: "4rem",
        desktop: "5rem",
      },
    },
    extend: {
      backgroundColor: theme => ({
        primary: '#38a169',
        secondary: '#ffed4a',
        danger: '#e3342f',
      }),
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ]
}