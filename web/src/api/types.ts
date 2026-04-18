/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export enum ModelsTransactionType {
  TransactionTypeExpense = "expense",
  TransactionTypeIncome = "income",
  TransactionTypeTransfer = "transfer",
}

export enum ModelsCategoryType {
  CategoryTypeExpense = "expense",
  CategoryTypeIncome = "income",
  CategoryTypeTransfer = "transfer",
}

export enum ModelsAccountType {
  AccountTypeExpense = "expense",
  AccountTypeIncome = "income",
  AccountTypeTransfer = "transfer",
}

export interface HandlerHealthResponse {
  status?: string;
  timestamp?: string;
  version?: string;
}

export interface ModelsAccount {
  balance?: number;
  created_at?: string;
  deleted_at?: string;
  id?: string;
  name?: string;
  type?: ModelsAccountType;
  updated_at?: string;
}

export interface ModelsAccountBreakdown {
  account_id?: string;
  account_name?: string;
  total?: number;
  type?: string;
}

export interface ModelsAccountPagedResponse {
  items?: ModelsAccount[];
  page_number?: number;
  page_size?: number;
  total_items?: number;
  total_pages?: number;
}

export interface ModelsAccountResponse {
  data?: ModelsAccount;
}

export interface ModelsCategory {
  created_at?: string;
  deleted_at?: string;
  id?: string;
  name?: string;
  type?: ModelsCategoryType;
  updated_at?: string;
}

export interface ModelsCategoryBreakdown {
  category_id?: string;
  category_name?: string;
  total?: number;
  type?: string;
}

export interface ModelsCategoryPagedResponse {
  items?: ModelsCategory[];
  page_number?: number;
  page_size?: number;
  total_items?: number;
  total_pages?: number;
}

export interface ModelsCategoryResponse {
  data?: ModelsCategory;
}

export interface ModelsCreateAccountRequest {
  balance?: number;
  name: string;
  type: "expense" | "income" | "transfer";
}

export interface ModelsCreateCategoryRequest {
  name: string;
  type: "expense" | "income" | "transfer";
}

export interface ModelsCreateDraftRequest {
  account_id: string;
  amount: number;
  category_id: string;
  currency: string;
  notes?: string;
  source: "manual" | "ingestion";
  title: string;
  transaction_type: "expense" | "income" | "transfer";
  transfer_account_id?: string;
}

export interface ModelsCreateTransactionRequest {
  account_id: string;
  amount: number;
  category_id: string;
  currency: string;
  notes?: string;
  title: string;
  transaction_type: "expense" | "income" | "transfer";
  transfer_account_id?: string;
}

export interface ModelsDraft {
  account_id?: string;
  amount?: number;
  category_id?: string;
  confirmed_at?: string;
  created_at?: string;
  currency?: string;
  currency_rate?: number;
  deleted_at?: string;
  id?: string;
  notes?: string;
  rejected_at?: string;
  source?: string;
  status?: string;
  title?: string;
  transaction_type?: string;
  transfer_account_id?: string;
  updated_at?: string;
}

export interface ModelsDraftPagedResponse {
  items?: ModelsDraft[];
  page_number?: number;
  page_size?: number;
  total_items?: number;
  total_pages?: number;
}

export interface ModelsDraftResponse {
  data?: ModelsDraft;
}

export interface ModelsErrorResponse {
  error?: string;
}

export interface ModelsReportByAccount {
  accounts?: ModelsAccountBreakdown[];
  period_end?: string;
  period_start?: string;
}

export interface ModelsReportByCategory {
  categories?: ModelsCategoryBreakdown[];
  period_end?: string;
  period_start?: string;
}

export interface ModelsReportSummary {
  base_currency?: string;
  period_end?: string;
  period_start?: string;
  total_balance?: number;
  total_expenses?: number;
  total_income?: number;
  total_transfers?: number;
}

export interface ModelsReportSummaryResponse {
  data?: ModelsReportSummary;
}

export interface ModelsReportTrends {
  data_points?: ModelsTrendDataPoint[];
  period_end?: string;
  period_start?: string;
}

export interface ModelsTransaction {
  account_id?: string;
  amount?: number;
  category_id?: string;
  created_at?: string;
  currency?: string;
  currency_rate?: number;
  deleted_at?: string;
  id?: string;
  notes?: string;
  title?: string;
  transaction_type?: ModelsTransactionType;
  transfer_account_id?: string;
  updated_at?: string;
}

export interface ModelsTransactionPagedResponse {
  items?: ModelsTransaction[];
  page_number?: number;
  page_size?: number;
  total_items?: number;
  total_pages?: number;
}

export interface ModelsTransactionResponse {
  data?: ModelsTransaction;
}

export interface ModelsTrendDataPoint {
  date?: string;
  total?: number;
  type?: string;
}

export interface ModelsUpdateAccountRequest {
  balance?: number;
  name?: string;
  type?: ModelsAccountType;
}

export interface ModelsUpdateCategoryRequest {
  name?: string;
  type?: ModelsCategoryType;
}

export interface ModelsUpdateDraftRequest {
  account_id?: string;
  amount?: number;
  category_id?: string;
  currency?: string;
  notes?: string;
  title?: string;
  transaction_type?: string;
  transfer_account_id?: string;
}

export interface ModelsUpdateTransactionRequest {
  account_id?: string;
  amount?: number;
  category_id?: string;
  currency?: string;
  notes?: string;
  title?: string;
  transaction_type?: ModelsTransactionType;
  transfer_account_id?: string;
}

export type QueryParamsType = Record<string | number, any>;
export type ResponseFormat = keyof Omit<Body, "body" | "bodyUsed">;

export interface FullRequestParams extends Omit<RequestInit, "body"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseFormat;
  /** request body */
  body?: unknown;
  /** base url */
  baseUrl?: string;
  /** request cancellation token */
  cancelToken?: CancelToken;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string;
  baseApiParams?: Omit<RequestParams, "baseUrl" | "cancelToken" | "signal">;
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<RequestParams | void> | RequestParams | void;
  customFetch?: typeof fetch;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown> extends Response {
  data: D;
  error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
  Json = "application/json",
  JsonApi = "application/vnd.api+json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = "";
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private abortControllers = new Map<CancelToken, AbortController>();
  private customFetch = (...fetchParams: Parameters<typeof fetch>) => fetch(...fetchParams);

  private baseApiParams: RequestParams = {
    credentials: "same-origin",
    headers: {},
    redirect: "follow",
    referrerPolicy: "no-referrer",
  };

  constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
    Object.assign(this, apiConfig);
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected encodeQueryParam(key: string, value: any) {
    const encodedKey = encodeURIComponent(key);
    return `${encodedKey}=${encodeURIComponent(typeof value === "number" ? value : `${value}`)}`;
  }

  protected addQueryParam(query: QueryParamsType, key: string) {
    return this.encodeQueryParam(key, query[key]);
  }

  protected addArrayQueryParam(query: QueryParamsType, key: string) {
    const value = query[key];
    return value.map((v: any) => this.encodeQueryParam(key, v)).join("&");
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {};
    const keys = Object.keys(query).filter((key) => "undefined" !== typeof query[key]);
    return keys
      .map((key) =>
        Array.isArray(query[key])
          ? this.addArrayQueryParam(query, key)
          : this.addQueryParam(query, key),
      )
      .join("&");
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery);
    return queryString ? `?${queryString}` : "";
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string")
        ? JSON.stringify(input)
        : input,
    [ContentType.JsonApi]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string")
        ? JSON.stringify(input)
        : input,
    [ContentType.Text]: (input: any) =>
      input !== null && typeof input !== "string" ? JSON.stringify(input) : input,
    [ContentType.FormData]: (input: any) => {
      if (input instanceof FormData) {
        return input;
      }

      return Object.keys(input || {}).reduce((formData, key) => {
        const property = input[key];
        formData.append(
          key,
          property instanceof Blob
            ? property
            : typeof property === "object" && property !== null
              ? JSON.stringify(property)
              : `${property}`,
        );
        return formData;
      }, new FormData());
    },
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  };

  protected mergeRequestParams(params1: RequestParams, params2?: RequestParams): RequestParams {
    return {
      ...this.baseApiParams,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.baseApiParams.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected createAbortSignal = (cancelToken: CancelToken): AbortSignal | undefined => {
    if (this.abortControllers.has(cancelToken)) {
      const abortController = this.abortControllers.get(cancelToken);
      if (abortController) {
        return abortController.signal;
      }
      return void 0;
    }

    const abortController = new AbortController();
    this.abortControllers.set(cancelToken, abortController);
    return abortController.signal;
  };

  public abortRequest = (cancelToken: CancelToken) => {
    const abortController = this.abortControllers.get(cancelToken);

    if (abortController) {
      abortController.abort();
      this.abortControllers.delete(cancelToken);
    }
  };

  public request = async <T = any, E = any>({
    body,
    secure,
    path,
    type,
    query,
    format,
    baseUrl,
    cancelToken,
    ...params
  }: FullRequestParams): Promise<HttpResponse<T, E>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.baseApiParams.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const queryString = query && this.toQueryString(query);
    const payloadFormatter = this.contentFormatters[type || ContentType.Json];
    const responseFormat = format || requestParams.format;

    return this.customFetch(
      `${baseUrl || this.baseUrl || ""}${path}${queryString ? `?${queryString}` : ""}`,
      {
        ...requestParams,
        headers: {
          ...(requestParams.headers || {}),
          ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
        },
        signal: (cancelToken ? this.createAbortSignal(cancelToken) : requestParams.signal) || null,
        body: typeof body === "undefined" || body === null ? null : payloadFormatter(body),
      },
    ).then(async (response) => {
      const r = response as HttpResponse<T, E>;
      r.data = null as unknown as T;
      r.error = null as unknown as E;

      const responseToParse = responseFormat ? response.clone() : response;
      const data = !responseFormat
        ? r
        : await responseToParse[responseFormat]()
            .then((data) => {
              if (r.ok) {
                r.data = data;
              } else {
                r.error = data;
              }
              return r;
            })
            .catch((e) => {
              r.error = e;
              return r;
            });

      if (cancelToken) {
        this.abortControllers.delete(cancelToken);
      }

      if (!response.ok) throw data;
      return data;
    });
  };
}

/**
 * @title No title
 * @contact
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  accounts = {
    /**
     * @description Get a paginated list of accounts with optional filtering
     *
     * @tags accounts
     * @name AccountsList
     * @summary List all accounts
     * @request GET:/accounts
     */
    accountsList: (
      query?: {
        /** Search query */
        q?: string;
        /** Sort by field */
        sort_by?: string;
        /** Sort order (asc/desc) */
        sort_order?: string;
        /**
         * Page number
         * @default 1
         */
        page?: number;
        /**
         * Page size
         * @default 10
         */
        page_size?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<ModelsAccountPagedResponse, ModelsErrorResponse>({
        path: `/accounts`,
        method: "GET",
        query: query,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Create a new account with name, type and optional balance
     *
     * @tags accounts
     * @name AccountsCreate
     * @summary Create a new account
     * @request POST:/accounts
     */
    accountsCreate: (request: ModelsCreateAccountRequest, params: RequestParams = {}) =>
      this.request<ModelsAccountResponse, ModelsErrorResponse>({
        path: `/accounts`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get a single account by its UUID
     *
     * @tags accounts
     * @name AccountsDetail
     * @summary Get account by ID
     * @request GET:/accounts/{id}
     */
    accountsDetail: (id: string, params: RequestParams = {}) =>
      this.request<ModelsAccountResponse, ModelsErrorResponse>({
        path: `/accounts/${id}`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Update an existing account by its UUID
     *
     * @tags accounts
     * @name AccountsUpdate
     * @summary Update an account
     * @request PUT:/accounts/{id}
     */
    accountsUpdate: (id: string, request: ModelsUpdateAccountRequest, params: RequestParams = {}) =>
      this.request<ModelsAccountResponse, ModelsErrorResponse>({
        path: `/accounts/${id}`,
        method: "PUT",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Soft delete an account by its UUID
     *
     * @tags accounts
     * @name AccountsDelete
     * @summary Delete an account
     * @request DELETE:/accounts/{id}
     */
    accountsDelete: (id: string, params: RequestParams = {}) =>
      this.request<ModelsAccountResponse, ModelsErrorResponse>({
        path: `/accounts/${id}`,
        method: "DELETE",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  categories = {
    /**
     * @description Get a paginated list of categories with optional filtering
     *
     * @tags categories
     * @name CategoriesList
     * @summary List all categories
     * @request GET:/categories
     */
    categoriesList: (
      query?: {
        /** Search query */
        q?: string;
        /** Sort by field */
        sort_by?: string;
        /** Sort order (asc/desc) */
        sort_order?: string;
        /**
         * Page number
         * @default 1
         */
        page?: number;
        /**
         * Page size
         * @default 10
         */
        page_size?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<ModelsCategoryPagedResponse, ModelsErrorResponse>({
        path: `/categories`,
        method: "GET",
        query: query,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Create a new category with name and type
     *
     * @tags categories
     * @name CategoriesCreate
     * @summary Create a new category
     * @request POST:/categories
     */
    categoriesCreate: (request: ModelsCreateCategoryRequest, params: RequestParams = {}) =>
      this.request<ModelsCategoryResponse, ModelsErrorResponse>({
        path: `/categories`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get a single category by its UUID
     *
     * @tags categories
     * @name CategoriesDetail
     * @summary Get category by ID
     * @request GET:/categories/{id}
     */
    categoriesDetail: (id: string, params: RequestParams = {}) =>
      this.request<ModelsCategoryResponse, ModelsErrorResponse>({
        path: `/categories/${id}`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Update an existing category by its UUID
     *
     * @tags categories
     * @name CategoriesUpdate
     * @summary Update a category
     * @request PUT:/categories/{id}
     */
    categoriesUpdate: (
      id: string,
      request: ModelsUpdateCategoryRequest,
      params: RequestParams = {},
    ) =>
      this.request<ModelsCategoryResponse, ModelsErrorResponse>({
        path: `/categories/${id}`,
        method: "PUT",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Soft delete a category by its UUID
     *
     * @tags categories
     * @name CategoriesDelete
     * @summary Delete a category
     * @request DELETE:/categories/{id}
     */
    categoriesDelete: (id: string, params: RequestParams = {}) =>
      this.request<ModelsCategoryResponse, ModelsErrorResponse>({
        path: `/categories/${id}`,
        method: "DELETE",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  drafts = {
    /**
     * @description Get a paginated list of drafts with optional filtering by source and status
     *
     * @tags drafts
     * @name DraftsList
     * @summary List all drafts
     * @request GET:/drafts
     */
    draftsList: (
      query?: {
        /** Filter by source (manual, ingestion) */
        source?: string;
        /** Filter by status (pending, confirmed, rejected) */
        status?: string;
        /**
         * Page size
         * @default 10
         */
        page_size?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<ModelsDraftPagedResponse, ModelsErrorResponse>({
        path: `/drafts`,
        method: "GET",
        query: query,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Create a new draft with transaction details
     *
     * @tags drafts
     * @name DraftsCreate
     * @summary Create a new draft
     * @request POST:/drafts
     */
    draftsCreate: (request: ModelsCreateDraftRequest, params: RequestParams = {}) =>
      this.request<ModelsDraftResponse, ModelsErrorResponse>({
        path: `/drafts`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get a single draft by its UUID
     *
     * @tags drafts
     * @name DraftsDetail
     * @summary Get draft by ID
     * @request GET:/drafts/{id}
     */
    draftsDetail: (id: string, params: RequestParams = {}) =>
      this.request<ModelsDraftResponse, ModelsErrorResponse>({
        path: `/drafts/${id}`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Hard delete a rejected draft
     *
     * @tags drafts
     * @name DraftsDelete
     * @summary Delete a draft
     * @request DELETE:/drafts/{id}
     */
    draftsDelete: (id: string, params: RequestParams = {}) =>
      this.request<ModelsErrorResponse, ModelsErrorResponse>({
        path: `/drafts/${id}`,
        method: "DELETE",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Update an existing draft by its UUID
     *
     * @tags drafts
     * @name DraftsPartialUpdate
     * @summary Update a draft
     * @request PATCH:/drafts/{id}
     */
    draftsPartialUpdate: (
      id: string,
      request: ModelsUpdateDraftRequest,
      params: RequestParams = {},
    ) =>
      this.request<ModelsDraftResponse, ModelsErrorResponse>({
        path: `/drafts/${id}`,
        method: "PATCH",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Promote a draft to a committed transaction
     *
     * @tags drafts
     * @name ConfirmCreate
     * @summary Confirm a draft
     * @request POST:/drafts/{id}/confirm
     */
    confirmCreate: (id: string, params: RequestParams = {}) =>
      this.request<ModelsTransactionResponse, ModelsErrorResponse>({
        path: `/drafts/${id}/confirm`,
        method: "POST",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Discard a draft - status set to rejected, no transaction created
     *
     * @tags drafts
     * @name RejectCreate
     * @summary Reject a draft
     * @request POST:/drafts/{id}/reject
     */
    rejectCreate: (id: string, params: RequestParams = {}) =>
      this.request<ModelsErrorResponse, ModelsErrorResponse>({
        path: `/drafts/${id}/reject`,
        method: "POST",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  health = {
    /**
     * @description Returns the health status of the API
     *
     * @tags health
     * @name HealthList
     * @summary Health check
     * @request GET:/health
     */
    healthList: (params: RequestParams = {}) =>
      this.request<HandlerHealthResponse, any>({
        path: `/health`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  reports = {
    /**
     * @description Get aggregated spending data grouped by account
     *
     * @tags reports
     * @name ByAccountList
     * @summary Get spending breakdown by account
     * @request GET:/reports/by-account
     */
    byAccountList: (
      query: {
        /** Start date (YYYY-MM-DD) */
        start_date: string;
        /** End date (YYYY-MM-DD) */
        end_date: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ModelsReportByAccount, ModelsErrorResponse>({
        path: `/reports/by-account`,
        method: "GET",
        query: query,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get aggregated spending data grouped by category
     *
     * @tags reports
     * @name ByCategoryList
     * @summary Get spending breakdown by category
     * @request GET:/reports/by-category
     */
    byCategoryList: (
      query: {
        /** Start date (YYYY-MM-DD) */
        start_date: string;
        /** End date (YYYY-MM-DD) */
        end_date: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ModelsReportByCategory, ModelsErrorResponse>({
        path: `/reports/by-category`,
        method: "GET",
        query: query,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get aggregated financial data including totals by type and category breakdown
     *
     * @tags reports
     * @name SummaryList
     * @summary Get financial report summary
     * @request GET:/reports/summary
     */
    summaryList: (
      query: {
        /** Start date (YYYY-MM-DD) */
        start_date: string;
        /** End date (YYYY-MM-DD) */
        end_date: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ModelsReportSummaryResponse, ModelsErrorResponse>({
        path: `/reports/summary`,
        method: "GET",
        query: query,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get time-series data for charting
     *
     * @tags reports
     * @name TrendsList
     * @summary Get time-series trend data
     * @request GET:/reports/trends
     */
    trendsList: (
      query: {
        /** Start date (YYYY-MM-DD) */
        start_date: string;
        /** End date (YYYY-MM-DD) */
        end_date: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ModelsReportTrends, ModelsErrorResponse>({
        path: `/reports/trends`,
        method: "GET",
        query: query,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  transactions = {
    /**
     * @description Get a paginated list of transactions with optional filtering
     *
     * @tags transactions
     * @name TransactionsList
     * @summary List all transactions
     * @request GET:/transactions
     */
    transactionsList: (
      query?: {
        /** Search query */
        q?: string;
        /** Filter by account ID */
        account_id?: string;
        /** Filter by category ID */
        category_id?: string;
        /** Filter by transaction type (expense, income, transfer) */
        transaction_type?: string;
        /** Sort by field (title, transacted_at, created_at, amount) */
        sort_by?: string;
        /** Sort order (asc/desc) */
        sort_order?: string;
        /**
         * Page number
         * @default 1
         */
        page?: number;
        /**
         * Page size
         * @default 10
         */
        page_size?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<ModelsTransactionPagedResponse, ModelsErrorResponse>({
        path: `/transactions`,
        method: "GET",
        query: query,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Create a new transaction with amount, type, and account/category references
     *
     * @tags transactions
     * @name TransactionsCreate
     * @summary Create a new transaction
     * @request POST:/transactions
     */
    transactionsCreate: (request: ModelsCreateTransactionRequest, params: RequestParams = {}) =>
      this.request<ModelsTransactionResponse, ModelsErrorResponse>({
        path: `/transactions`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get a single transaction by its UUID
     *
     * @tags transactions
     * @name TransactionsDetail
     * @summary Get transaction by ID
     * @request GET:/transactions/{id}
     */
    transactionsDetail: (id: string, params: RequestParams = {}) =>
      this.request<ModelsTransactionResponse, ModelsErrorResponse>({
        path: `/transactions/${id}`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Update an existing transaction by its UUID
     *
     * @tags transactions
     * @name TransactionsUpdate
     * @summary Update a transaction
     * @request PUT:/transactions/{id}
     */
    transactionsUpdate: (
      id: string,
      request: ModelsUpdateTransactionRequest,
      params: RequestParams = {},
    ) =>
      this.request<ModelsTransactionResponse, ModelsErrorResponse>({
        path: `/transactions/${id}`,
        method: "PUT",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Soft delete a transaction by its UUID
     *
     * @tags transactions
     * @name TransactionsDelete
     * @summary Delete a transaction
     * @request DELETE:/transactions/{id}
     */
    transactionsDelete: (id: string, params: RequestParams = {}) =>
      this.request<ModelsTransactionResponse, ModelsErrorResponse>({
        path: `/transactions/${id}`,
        method: "DELETE",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
}
