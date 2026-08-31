export const theme = {
  colors: {
    cream: "#FFFBF5",
    creamSoft: "#FFF8F0",
    card: "#FFFFFF",
    border: "#E8E0D8",
    borderSoft: "#EDE8E3",
    charcoal: "#1E1E1E",
    charcoalMuted: "#6B6B6B",
    terracotta: "#E76F51",
    terracottaDark: "#D94F2B",
    sage: "#8BA888",
    sageMuted: "#6B7F6B",
    mustard: "#E9C46A",
    shadow: "rgba(0,0,0,0.08)",
  },
  radius: {
    sm: 12,
    md: 16,
    lg: 20,
    pill: 28,
    card: 20,
  },
  spacing: {
    xs: 8,
    sm: 12,
    md: 16,
    lg: 24,
    xl: 32,
  },
  font: {
    heading: "Playfair_700Bold",
    body: "Inter_400Regular",
    bodyMedium: "Inter_500Medium",
  },
  text: {
    title: 32,
    subtitle: 18,
    body: 15,
    caption: 13,
    label: 11,
  },
} as const;

export type Theme = typeof theme;
