import { mockHealth } from "@/mocks/data/fixtures";
import { http, HttpResponse } from "msw";

export const healthHandlers = [
  http.get("/health", () => {
    return HttpResponse.json(mockHealth);
  }),
];
