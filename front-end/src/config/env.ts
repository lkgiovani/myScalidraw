export const env = {
  API_BASE_URL:
    import.meta.env.VITE_API_BASE_URL || "http://localhost:8181/api",
} as const;
