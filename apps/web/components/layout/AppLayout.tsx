"use client";

import { AppShell, Box } from "@mantine/core";
import AppHeader from "./AppHeader";

interface AppLayoutProps {
  children: React.ReactNode;
  navbar?: React.ReactNode;
  header?: React.ReactNode;
}

export function AppLayout({ children, navbar, header }: AppLayoutProps) {
  return (
    <AppShell
      padding={0}
      header={{ height: 72 }}
      navbar={
        navbar
          ? {
              width: 288,
              breakpoint: "sm",
            }
          : undefined
      }
      styles={{
        root: {
          background:
            "linear-gradient(135deg, var(--mantine-color-dark-9) 0%, #13201f 42%, #101521 100%)",
        },
        header: {
          borderColor: "rgba(255,255,255,0.08)",
          background: "rgba(10, 16, 24, 0.82)",
          backdropFilter: "blur(16px)",
        },
        navbar: {
          borderColor: "rgba(255,255,255,0.08)",
          background: "rgba(10, 16, 24, 0.62)",
          backdropFilter: "blur(16px)",
        },
        main: {
          minHeight: "100vh",
        },
      }}
    >
      <AppShell.Header>
        {header ? header : <AppHeader />}
      </AppShell.Header>

      {navbar ? <AppShell.Navbar p="md">{navbar}</AppShell.Navbar> : null}

      <AppShell.Main>
        <Box px={{ base: "md", md: "xl" }} py={{ base: "md", md: "xl" }}>
          {children}
        </Box>
      </AppShell.Main>
    </AppShell>
  );
}
