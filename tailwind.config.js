/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: "selector",
  content: ["./views/**/*.{go,js,templ,html}"],
  theme: {
    colors: {
      white: "#DADDE1",
      black: "#1E2030",
      gray: "#2F334D",
      blue: "#82AAFF",
      green: "#C3E88D",
      red: "#FF007C",
    },
    extend: {},
  },
  plugins: [require("@tailwindcss/forms")],
};
