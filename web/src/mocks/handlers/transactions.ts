import { mockTransactions } from "@/mocks/data/fixtures";
import { http, HttpResponse } from "msw";

export const transactionsHandlers = [
  http.get("/transactions", () => {
    return HttpResponse.json({
      items: mockTransactions,
      page_number: 1,
      page_size: 10,
      total_items: mockTransactions.length,
      total_pages: 1,
    });
  }),
  http.post("/transactions", async ({ request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    const newTransaction = {
      id: `txn-${Date.now()}`,
      ...body,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };
    return HttpResponse.json({ data: newTransaction }, { status: 201 });
  }),
  http.get("/transactions/:id", ({ params }) => {
    const transaction = mockTransactions.find((t) => t.id === params.id);
    if (!transaction) {
      return HttpResponse.json({ error: "Transaction not found" }, { status: 404 });
    }
    return HttpResponse.json({ data: transaction });
  }),
  http.put("/transactions/:id", async ({ params, request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    const existing = mockTransactions.find((t) => t.id === params.id);
    if (!existing) {
      return HttpResponse.json({ error: "Transaction not found" }, { status: 404 });
    }
    return HttpResponse.json({
      data: { ...existing, ...body, updated_at: new Date().toISOString() },
    });
  }),
  http.delete("/transactions/:id", ({ params }) => {
    const exists = mockTransactions.find((t) => t.id === params.id);
    if (!exists) {
      return HttpResponse.json({ error: "Transaction not found" }, { status: 404 });
    }
    return HttpResponse.json({ data: null });
  }),
];
