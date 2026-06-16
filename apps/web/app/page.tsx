"use client";
import { NavLink, Stack, Text } from "@mantine/core";
import { RouteGuard } from "../providers/RouteGuard";
import { useState } from "react";
import { IconClipboardList } from "@tabler/icons-react";
import { AppLayout } from "../components/layout/AppLayout";
import { AuditLogDashboard } from "../components/layout/auditlog/AuditLogDashboard";
type section = "audit-log"|"inbox";

export default function HomePage() {
  const [currentSection, setCurrentSection] = useState<section>("audit-log");
  const sidebarnav = (
    <Stack>
      <NavLink
        label="Audit Log"
        leftSection={<IconClipboardList />}
        active={currentSection === "audit-log"}
        onClick={() => {
          setCurrentSection("audit-log");
        }}
      />
      <NavLink
        label="Inbox"
        leftSection={<IconClipboardList />}
        active={currentSection === "inbox"}
        onClick={() => {
          setCurrentSection("inbox");
        }}
      />
    </Stack>
  );
  return (
    // <RouteGuard>
    <AppLayout navbar={sidebarnav}>
      {currentSection === "audit-log" && <AuditLogDashboard />}
      {currentSection === "inbox" && <Text>Inbox</Text>}
    </AppLayout>

    // </RouteGuard>
  );
}
