import { ref } from "vue";
import { useApi } from "@/composables/use-api";
import type {
  ModelsReportSummary,
  ModelsReportByCategory,
  ModelsReportByAccount,
  ModelsReportTrends,
} from "@/api/types";

export function useReports() {
  const { api, loading, error, wrap } = useApi();

  const summary = ref<ModelsReportSummary | null>(null);
  const byCategory = ref<ModelsReportByCategory | null>(null);
  const byAccount = ref<ModelsReportByAccount | null>(null);
  const trends = ref<ModelsReportTrends | null>(null);

  async function fetchSummary(startDate: string, endDate: string) {
    return wrap(async () => {
      const res = await api.reports.summaryList({ start_date: startDate, end_date: endDate });
      summary.value = res.data?.data || null;
    });
  }

  async function fetchByCategory(startDate: string, endDate: string) {
    return wrap(async () => {
      const res = await api.reports.byCategoryList({ start_date: startDate, end_date: endDate });
      byCategory.value = res.data || null;
    });
  }

  async function fetchByAccount(startDate: string, endDate: string) {
    return wrap(async () => {
      const res = await api.reports.byAccountList({ start_date: startDate, end_date: endDate });
      byAccount.value = res.data || null;
    });
  }

  async function fetchTrends(startDate: string, endDate: string) {
    return wrap(async () => {
      const res = await api.reports.trendsList({ start_date: startDate, end_date: endDate });
      trends.value = res.data || null;
    });
  }

  return {
    summary,
    byCategory,
    byAccount,
    trends,
    loading,
    error,
    fetchSummary,
    fetchByCategory,
    fetchByAccount,
    fetchTrends,
  };
}