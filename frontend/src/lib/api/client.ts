import createClient from "openapi-fetch";
import type { paths } from "./v1";

export const api = createClient<paths>({ baseUrl: "/api" });

api.use({
  async onResponse({ response }) {
    if (!response.ok) {
      const body = await response.clone().text();
      if (!body) {
        return new Response(JSON.stringify({ detail: response.statusText || "Request failed" }), {
          status: response.status,
          headers: { "Content-Type": "application/json" },
        });
      }
    }
    return undefined;
  },
});
