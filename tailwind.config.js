/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './resources/**/*.{vue,js,ts,jsx,tsx}',
    './resources/views/root.html',
  ],
  theme: {
    darkMode: 'class',
    extend: {
      borderRadius: {
        lg: '12px',
        md: '8px',
        sm: '4px',
      }
    },
  },
  plugins: [
  ],
}
