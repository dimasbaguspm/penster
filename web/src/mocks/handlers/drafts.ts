import { mockDrafts } from "@/mocks/data/fixtures";
import { http, HttpResponse } from "msw";

export const draftsHandlers = [
  http.get("/drafts", () => {
    return HttpResponse.json({
      items: mockDrafts,
      page_number: 1,
      page_size: 10,
      total_items: mockDrafts.length,
      total_pages: 1,
    });
  }),
  http.post("/drafts", async ({ request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    const newDraft = {
      id: `dft-${Date.now()}`,
      ...body,
      status: "pending",
      source: "manual",
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };
    return HttpResponse.json({ data: newDraft }, { status: 201 });
  }),
  http.get("/drafts/:id", ({ params }) => {
    const draft = mockDrafts.find((d) => d.id === params.id);
    if (!draft) {
      return HttpResponse.json({ error: "Draft not found" }, { status: 404 });
    }
    return HttpResponse.json({ data: draft });
  }),
  http.patch("/drafts/:id", async ({ params, request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    const existing = mockDrafts.find((d) => d.id === params.id);
    if (!existing) {
      return HttpResponse.json({ error: "Draft not found" }, { status: 404 });
    }
    return HttpResponse.json({
      data: { ...existing, ...body, updated_at: new Date().toISOString() },
    });
  }),
  http.delete("/drafts/:id", ({ params }) => {
    const exists = mockDrafts.find((d) => d.id === params.id);
    if (!exists) {
      return HttpResponse.json({ error: "Draft not found" }, { status: 404 });
    }
    return HttpResponse.json({ data: null });
  }),
  http.post("/drafts/:id/confirm", ({ params }) => {
    const draft = mockDrafts.find((d) => d.id === params.id);
    if (!draft) {
      return HttpResponse.json({ error: "Draft not found" }, { status: 404 });
    }
    const confirmedTransaction = {
      id: `txn-${Date.now()}`,
      account_id: draft.account_id,
      category_id: draft.category_id,
      amount: draft.amount,
      currency: draft.currency,
      title: draft.title,
      transaction_type: draft.transaction_type,
      transfer_account_id: draft.transfer_account_id,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };
    return HttpResponse.json({ data: confirmedTransaction }, { status: 201 });
  }),
  http.post("/drafts/:id/reject", ({ params }) => {
    const draft = mockDrafts.find((d) => d.id === params.id);
    if (!draft) {
      return HttpResponse.json({ error: "Draft not found" }, { status: 404 });
    }
    return HttpResponse.json({ data: null });
  }),
];
