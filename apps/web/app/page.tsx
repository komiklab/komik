"use client";

import {
  Badge,
  Box,
  Group,
  NavLink,
  Stack,
  Text,
  ThemeIcon,
} from "@mantine/core";
import {
  IconClipboardList,
  IconFishHook,
  IconInbox,
  IconTimelineEvent,
} from "@tabler/icons-react";
import { AppLayout } from "../components/layout/AppLayout";
import AgentPanel from "../components/layout/agent/agentPanel";
import { AuditLogDashboard } from "../components/layout/auditlog/AuditLogDashboard";
import HooksPanel from "../components/layout/hooks/hookspanel";
import { InboxPanel } from "../components/layout/inbox/InboxPanel";
import { useUiStore, type WorkspaceSection } from "../stores/ui";
import { RouteGuard } from "../providers/RouteGuard";

const navItems: Array<{
  value: WorkspaceSection;
  label: string;
  description: string;
  icon: typeof IconClipboardList;
}> = [
  {
    value: "audit-log",
    label: "Audit Log",
    description: "System events and security trail",
    icon: IconTimelineEvent,
  },
  {
    value: "inbox",
    label: "Inbox",
    description: "Incoming work queue preview",
    icon: IconInbox,
  },
  {
    value: "agent",
    label: "Agents",
    description: "Configure AI agents",
    icon: IconClipboardList,
  },
  {
    value: "hooks",
    label: "Hooks",
    description: "Configure webhooks",
    icon: IconFishHook,
  },
];

export default function HomePage() {
  const currentSection = useUiStore((state) => state.currentSection);
  const setCurrentSection = useUiStore((state) => state.setCurrentSection);

  const sidebarNav = (
    <Stack gap="lg">
      <Box>
        <Text tt="uppercase" fz={11} fw={800} c="dimmed" mb="xs">
          Workspace
        </Text>
        <Stack gap={6}>
          {navItems.map((item) => {
            const Icon = item.icon;

            return (
              <NavLink
                key={item.value}
                active={currentSection === item.value}
                label={
                  <Group justify="space-between" gap="xs" wrap="nowrap">
                    <Text fw={700}>{item.label}</Text>
                    {item.value === "inbox" ? (
                      <Badge size="xs" color="gray" variant="light">
                        Mock
                      </Badge>
                    ) : null}
                  </Group>
                }
                description={item.description}
                leftSection={
                  <ThemeIcon
                    size={34}
                    radius="md"
                    variant={currentSection === item.value ? "filled" : "light"}
                    color={currentSection === item.value ? "teal" : "gray"}
                  >
                    <Icon size={18} />
                  </ThemeIcon>
                }
                onClick={() => setCurrentSection(item.value)}
                styles={{
                  root: {
                    borderRadius: 8,
                    border:
                      currentSection === item.value
                        ? "1px solid rgba(32, 201, 151, 0.34)"
                        : "1px solid transparent",
                  },
                }}
              />
            );
          })}
        </Stack>
      </Box>
    </Stack>
  );

  return (
    <RouteGuard>
      <AppLayout navbar={sidebarNav}>
        {currentSection === "audit-log" && <AuditLogDashboard />}
        {currentSection === "inbox" && <InboxPanel />}
        {currentSection === "agent" && <AgentPanel />}
        {currentSection === "hooks" && <HooksPanel />}
      </AppLayout>
    </RouteGuard>
  );
}
