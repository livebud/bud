const colors = require("tailwindcss/colors")

module.exports = {
  mode: "jit",
  corePlugins: {
    preflight: false,
  },
  theme: {
    extend: {
      colors: {
        orange: colors.orange,
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
