"use client";

import { createTheme } from "@mantine/core";

export const theme = createTheme({
  primaryColor: "teal",
  fontFamily:
    "Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, Segoe UI, sans-serif",
  headings: {
    fontFamily:
      "Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, Segoe UI, sans-serif",
    fontWeight: "700",
  },
  defaultRadius: "md",
  colors: {
    ink: [
      "#f4f7fb",
      "#dbe4ef",
      "#b7c7dd",
      "#90a7c8",
      "#708db7",
      "#5e7daf",
      "#5274ab",
      "#426296",
      "#395782",
      "#2f4b70",
    ],
  },
  components: {
    Paper: {
      defaultProps: {
        radius: "md",
      },
    },
    Button: {
      defaultProps: {
        radius: "md",
      },
    },
    NavLink: {
      defaultProps: {
        variant: "subtle",
      },
    },
  },
});
