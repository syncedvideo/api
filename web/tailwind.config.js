module.exports = {
  purge: [],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      colors: {
        brand: '#32e0c4',
        primary: {
          '50': '#f5fdfc',
          '100': '#ebfcf9',
          '200': '#ccf7f0',
          '300': '#adf3e7',
          '400': '#70e9d6',
          '500': '#32e0c4',
          '600': '#2dcab0',
          '700': '#26a893',
          '800': '#1e8676',
          '900': '#196e60'
        },
        gray: {
          '50': '#f4f4f4',
          '100': '#e9e9e9',
          '200': '#c8c8c8',
          '300': '#a6a6a6',
          '400': '#646464',
          '500': '#212121',
          '600': '#1e1e1e',
          '700': '#191919',
          '800': '#141414',
          '900': '#101010'
        }
      }
    }
  },
  variants: {
    extend: {
      opacity: ['disabled']
    }
  },
  plugins: []
}
