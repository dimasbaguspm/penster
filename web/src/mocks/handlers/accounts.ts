import { mockAccounts } from "@/mocks/data/fixtures";
import { http, HttpResponse } from "msw";

export const accountsHandlers = [
  http.get("/accounts", () => {
    return HttpResponse.json({
      items: mockAccounts,
      page_number: 1,
      page_size: 10,
      total_items: mockAccounts.length,
      total_pages: 1,
    });
  }),
  http.post("/accounts", async ({ request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    const newAccount = {
      id: `acc-${Date.now()}`,
      ...body,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };
    return HttpResponse.json({ data: newAccount }, { status: 201 });
  }),
  http.get("/accounts/:id", ({ params }) => {
    const account = mockAccounts.find((a) => a.id === params.id);
    if (!account) {
      return HttpResponse.json({ error: "Account not found" }, { status: 404 });
    }
    return HttpResponse.json({ data: account });
  }),
  http.put("/accounts/:id", async ({ params, request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    const existing = mockAccounts.find((a) => a.id === params.id);
    if (!existing) {
      return HttpResponse.json({ error: "Account not found" }, { status: 404 });
    }
    return HttpResponse.json({
      data: { ...existing, ...body, updated_at: new Date().toISOString() },
    });
  }),
  http.delete("/accounts/:id", ({ params }) => {
    const exists = mockAccounts.find((a) => a.id === params.id);
    if (!exists) {
      return HttpResponse.json({ error: "Account not found" }, { status: 404 });
    }
    return HttpResponse.json({ data: null });
  }),
];
