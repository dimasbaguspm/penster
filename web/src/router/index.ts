import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      component: () => import("@/layouts/default-layout.vue"),
      children: [
        {
          path: "",
          name: "dashboard",
          component: () => import("@/features/dashboard/views/dashboard-view.vue"),
        },
        {
          path: "accounts",
          name: "accounts",
          component: () => import("@/features/accounts/views/accounts-list-view.vue"),
        },
        {
          path: "accounts/:id",
          name: "account-detail",
          component: () => import("@/features/accounts/views/account-detail-view.vue"),
        },
        {
          path: "transactions",
          name: "transactions",
          component: () => import("@/features/transactions/views/transactions-list-view.vue"),
        },
        {
          path: "transactions/:id",
          name: "transaction-detail",
          component: () => import("@/features/transactions/views/transaction-detail-view.vue"),
        },
        {
          path: "drafts",
          name: "drafts",
          component: () => import("@/features/drafts/views/drafts-view.vue"),
        },
        {
          path: "reports",
          name: "reports",
          component: () => import("@/features/reports/views/reports-view.vue"),
        },
      ],
    },
  ],
});

export default router;
