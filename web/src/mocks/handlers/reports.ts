import {
  mockReportByAccount,
  mockReportByCategory,
  mockReportSummary,
  mockReportTrends,
} from "@/mocks/data/fixtures";
import { http, HttpResponse } from "msw";

export const reportsHandlers = [
  http.get("/reports/summary", () => {
    return HttpResponse.json({ data: mockReportSummary });
  }),
  http.get("/reports/by-category", () => {
    return HttpResponse.json({ data: mockReportByCategory });
  }),
  http.get("/reports/by-account", () => {
    return HttpResponse.json({ data: mockReportByAccount });
  }),
  http.get("/reports/trends", () => {
    return HttpResponse.json({ data: mockReportTrends });
  }),
];
