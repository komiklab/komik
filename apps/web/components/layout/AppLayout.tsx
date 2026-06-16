"use client";

import { AppShell, Burger } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import AppHeader from "./AppHeader";
interface AppLayoutProps {
  children: React.ReactNode;
  navbar?: React.ReactNode;
}

export function AppLayout({ children, navbar }: AppLayoutProps) {
  return (

    <AppShell
      padding="md"
      header={{ height: 60 }}
      navbar={navbar ? 
        {
        width: 300,
        breakpoint: 'sm',
      } : undefined}
    >
      <AppShell.Header>
        <AppHeader/>
      </AppShell.Header>

      <AppShell.Navbar>{navbar}</AppShell.Navbar>

      <AppShell.Main>{children}</AppShell.Main>
    </AppShell>
      
  );
}
