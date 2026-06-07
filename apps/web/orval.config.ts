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
    },
  },
});