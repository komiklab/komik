import { defineConfig } from "orval";

export default defineConfig({
  komik: {
    input: {
      target: "../../libs/openapi/openapi.yaml",
    },
    output: {
      target: "api/komik.ts",
      schemas: "api/schemas",
      client: "react-query",
      baseUrl: "http://localhost:65080/api/v1",
      override: {
        mutator: {
          // All generated API functions will call `customInstance` from this
          // file, which injects the X-CSRF-Token header automatically.
          path: "./api/httpClient.ts",
          name: "customInstance",
        },
      },
    },
  },
});