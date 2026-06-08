/**
 * Custom Axios HTTP client for the Komik API.
 *
 * Echo's CSRF middleware (middleware.CSRFWithConfig) sets a `_csrf` cookie and
 * expects the token to be sent back as the `X-CSRF-Token` request header on
 * every mutating (non-GET/HEAD/OPTIONS) request.
 *
 * This module:
 *  1. Creates a pre-configured Axios instance (`AXIOS_INSTANCE`) with:
 *     - `withCredentials: true` so the browser sends and receives cookies on
 *       cross-origin requests to the API.
 *     - A request interceptor that reads `_csrf` from `document.cookie` and
 *       injects it as the `X-CSRF-Token` header.
 *  2. Exports `customInstance`, the function signature expected by Orval's
 *     `override.mutator` option. The generated code passes a full URL string +
 *     RequestInit options; `customInstance` returns the full
 *     `{ data, status, headers }` envelope so callers always have the HTTP
 *     status code available for decision-making.
 */

import axios from "axios";
import type { AxiosRequestConfig } from "axios";

// ---------------------------------------------------------------------------
// Axios instance
// ---------------------------------------------------------------------------

export const AXIOS_INSTANCE = axios.create({
  baseURL: "http://localhost:65080/api/v1",
  withCredentials: true, // send _csrf cookie on cross-origin requests
});

// ---------------------------------------------------------------------------
// CSRF token store & interceptors
// ---------------------------------------------------------------------------

/**
 * In-memory store for the CSRF token.
 *
 * The backend (GetAdmin handler) returns the current CSRF token in the
 * `X-CSRF-Token` response header. The response interceptor below captures it
 * here so the request interceptor can attach it to every subsequent mutating
 * request — no `document.cookie` cross-origin parsing needed.
 */
let csrfToken: string | null = null;

const ECHO_SEC_FETCH_SITE_CSRF_TOKEN = "_echo_csrf_using_sec_fetch_site_";

/** Response interceptor — capture real CSRF tokens whenever the backend sends them. */
AXIOS_INSTANCE.interceptors.response.use((response) => {
  const token = response.headers["x-csrf-token"];
  if (token && token !== ECHO_SEC_FETCH_SITE_CSRF_TOKEN) {
    csrfToken = token;
  }
  return response;
});

/** Request interceptor — attach the stored CSRF token to every request. */
AXIOS_INSTANCE.interceptors.request.use((config) => {
  if (csrfToken) {
    config.headers["X-CSRF-Token"] = csrfToken;
  }
  return config;
});

// ---------------------------------------------------------------------------
// Orval custom mutator
// ---------------------------------------------------------------------------

/**
 * The response envelope returned by every generated API operation.
 *
 * Orval's fetch client convention wraps responses as `{ data, status, headers }`
 * so callers always have the HTTP status code available for conditional logic.
 */
export type ApiResponse<TData> = {
  data: TData;
  status: number;
  headers: Record<string, string>;
};

/**
 * The function Orval calls for every generated API operation.
 *
 * The generated code (fetch client + custom mutator) passes the full URL
 * string as the first argument and a plain `RequestInit`-like object as the
 * second argument, e.g.:
 *
 *   customInstance<getAdminResponse>('http://localhost:65080/api/v1/admin', { method: 'GET' })
 *
 * We normalise both arguments into a proper `AxiosRequestConfig`, make the
 * request via `AXIOS_INSTANCE`, and re-wrap the Axios response into the
 * `{ data, status, headers }` envelope the generated types declare.
 */

/** Normalise a Fetch-style `HeadersInit` into a plain record for Axios. */
function normalizeHeaders(
  headers?: HeadersInit,
): Record<string, string> | undefined {
  if (!headers) return undefined;
  if (headers instanceof Headers) {
    const result: Record<string, string> = {};
    headers.forEach((value, key) => {
      result[key] = value;
    });
    return result;
  }
  if (Array.isArray(headers)) {
    return Object.fromEntries(headers);
  }
  return headers as Record<string, string>;
}

// `customInstance` is parameterised over the *inner* data type `TData`.
// The full envelope `ApiResponse<TData>` is the actual return value, but we
// cast to `TData` so orval's generated `Promise<getAdminResponse>` signature
// is satisfied (orval passes the envelope type as `T`).
export const customInstance = <TData>(
  // Orval passes the full URL string as the first arg when using the fetch
  // client + custom mutator combination.
  urlOrConfig: string | AxiosRequestConfig,
  // The second arg comes from the generated code as RequestInit.
  options?: RequestInit,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
): Promise<any> => {
  const source = axios.CancelToken.source();

  // Normalise headers from Fetch-style HeadersInit → plain record for Axios.
  const axiosOptions: AxiosRequestConfig = options
    ? {
        ...options,
        headers: normalizeHeaders(options.headers),
        // Axios uses `data` for the request body; Fetch uses `body`.
        data: options.body,
        // Fetch allows signal: null, Axios expects undefined.
        signal: options.signal ?? undefined,
      }
    : {};

  // Lift a plain URL string into an AxiosRequestConfig.
  const resolvedConfig: AxiosRequestConfig =
    typeof urlOrConfig === "string"
      ? { url: urlOrConfig, ...axiosOptions }
      : { ...urlOrConfig, ...axiosOptions };

  const promise = AXIOS_INSTANCE<TData>({
    ...resolvedConfig,
    cancelToken: source.token,
  }).then((response): ApiResponse<TData> => ({
    // Re-wrap into the { data, status, headers } envelope that all
    // orval-generated types declare. This gives every caller access to the
    // HTTP status code for conditional logic without any extra plumbing.
    data: response.data,
    status: response.status,
    headers: normalizeHeaders(response.headers as HeadersInit) ?? {},
  }));

  // Attach a cancel helper so React Query can abort in-flight requests when a
  // component unmounts or a new request supersedes the previous one.
  (promise as ReturnType<typeof promise> & { cancel?: () => void }).cancel =
    () => {
      source.cancel("Query was cancelled by React Query");
    };

  return promise;
};
