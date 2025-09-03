function getApiBaseUrl(): string {
  if (import.meta.env.VITE_API_BASE_URL) {
    return import.meta.env.VITE_API_BASE_URL;
  }

  const metaTag = document.querySelector('meta[name="backend-url"]');
  if (metaTag) {
    const content = metaTag.getAttribute("content");
    if (content && content !== "__BACKEND_BASE_URL__") {
      return content;
    }
  }

  const runtimeUrl = "__BACKEND_BASE_URL__";
  if (runtimeUrl !== "__BACKEND_BASE_URL__") {
    return runtimeUrl;
  }

  if (typeof window !== "undefined") {
    const { protocol, hostname } = window.location;
    if (hostname === "localhost" && window.location.port === "5173") {
      return "http://localhost:8181/api";
    }
    return `${protocol}//${hostname}:8181/api`;
  }

  return "http://localhost:8181/api";
}

export const env = {
  API_BASE_URL: getApiBaseUrl(),
} as const;
