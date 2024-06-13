/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: "selector",
  content: ["./views/**/*.{go,js,templ,html}"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/forms")],
};
