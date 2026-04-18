import { accountsHandlers } from "@/mocks/handlers/accounts";
import { categoriesHandlers } from "@/mocks/handlers/categories";
import { draftsHandlers } from "@/mocks/handlers/drafts";
import { healthHandlers } from "@/mocks/handlers/health";
import { reportsHandlers } from "@/mocks/handlers/reports";
import { transactionsHandlers } from "@/mocks/handlers/transactions";
import { setupWorker } from "msw/browser";

const handlers = [
  ...accountsHandlers,
  ...categoriesHandlers,
  ...transactionsHandlers,
  ...draftsHandlers,
  ...reportsHandlers,
  ...healthHandlers,
];

export const worker = setupWorker(...handlers);
