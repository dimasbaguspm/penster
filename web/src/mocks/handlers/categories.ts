import { mockCategories } from "@/mocks/data/fixtures";
import { http, HttpResponse } from "msw";

export const categoriesHandlers = [
  http.get("/categories", () => {
    return HttpResponse.json({
      items: mockCategories,
      page_number: 1,
      page_size: 10,
      total_items: mockCategories.length,
      total_pages: 1,
    });
  }),
  http.post("/categories", async ({ request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    const newCategory = {
      id: `cat-${Date.now()}`,
      ...body,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };
    return HttpResponse.json({ data: newCategory }, { status: 201 });
  }),
  http.get("/categories/:id", ({ params }) => {
    const category = mockCategories.find((c) => c.id === params.id);
    if (!category) {
      return HttpResponse.json({ error: "Category not found" }, { status: 404 });
    }
    return HttpResponse.json({ data: category });
  }),
  http.put("/categories/:id", async ({ params, request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    const existing = mockCategories.find((c) => c.id === params.id);
    if (!existing) {
      return HttpResponse.json({ error: "Category not found" }, { status: 404 });
    }
    return HttpResponse.json({
      data: { ...existing, ...body, updated_at: new Date().toISOString() },
    });
  }),
  http.delete("/categories/:id", ({ params }) => {
    const exists = mockCategories.find((c) => c.id === params.id);
    if (!exists) {
      return HttpResponse.json({ error: "Category not found" }, { status: 404 });
    }
    return HttpResponse.json({ data: null });
  }),
];
