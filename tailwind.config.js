/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: "selector",
  content: ["./views/**/*.{go,js,templ,html}"],
  theme: {
    // colors: {
    //   white: "#DADDE1",
    //   black: "#1E2030",
    //   gray: "#2F334D",
    //   blue: "#82AAFF",
    //   green: "#C3E88D",
    //   red: "#FF007C",
    // },
    extend: {
      colors: {
        tokyonight: {
          background: "#1a1b26",
          foreground: "#c0caf5",
          current: "#c0caf5",
          comment: "#565f89",
          cyan: "#7dcfff",
          blue: "#7aa2f7",
          purple: "#bb9af7",
          orange: "#ff9e64",
          yellow: "#e0af68",
          green: "#9ece6a",
          magenta: "#ff007c",
          red: "#f7768e",
        },
      },
    },
  },
  plugins: [require("@tailwindcss/forms")],
};
