/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./index.html", "./src/**/*.{js,jsx}"],
  mode: "jit",
  theme: {
    extend: {
      colors: {
        // ── design tokens (แก้ธีมที่นี่ที่เดียว) ──
        ink: "#18181b",      // ตัวอักษรหลัก (เกือบดำ)
        body: "#3f3f46",     // ตัวอักษรเนื้อหา
        muted: "#71717a",    // ตัวอักษรรอง
        accent: "#4f46e5",   // สีเน้น (indigo)
        accentDark: "#4338ca",
        line: "#e4e4e7",     // เส้นขอบ
        surface: "#fafafa",  // พื้นหลังโซน
        // เก็บไว้เผื่อโค้ดเก่าอ้างถึง
        primary: "#00040f",
        myblogbg: "#18181b",
      },
      fontFamily: {
        sans: ["Inter", "system-ui", "sans-serif"],
      },
      maxWidth: {
        prose: "42rem", // ความกว้างเหมาะกับการอ่าน
      },
    },
    screens: {
      xs: "480px",
      ss: "620px",
      sm: "768px",
      md: "1060px",
      lg: "1200px",
      xl: "1700px",
    },
  },
  plugins: [],
};
