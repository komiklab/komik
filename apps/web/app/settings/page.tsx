"use client";
import { StringLiteral } from "typescript";
import { AppLayout } from "../../components/layout/AppLayout";
import SettingsHeader from "../../components/layout/settings/SettingsHeader";
import SettingsSidebarNav from "../../components/layout/settings/SettingsSidebarNav";
import { SettingsSection, useSettingsSectionStore } from "../../stores/settingssection";
import { IconClipboardList } from "@tabler/icons-react";
import { Stack, Box, Text, NavLink, Group, Badge, ThemeIcon } from "@mantine/core";
import AgentPanel from "../../components/layout/agent/agentPanel";
const navItems: Array<{
    value: SettingsSection
    label: string;
    description: String;
    icon: typeof IconClipboardList
}> = [
    {
        value: "agent",
        label: "Agent Settings",
        description: "Configure the AI agent",
        icon: IconClipboardList,
    }
]
export default function SettingsPage() {
  const currentSection = useSettingsSectionStore((state) => state.currentSection)
  const setCurrentSection = useSettingsSectionStore((state) => state.setCurrentSection)
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
    <AppLayout header={<SettingsHeader />} navbar={sidebarNav}>
    {currentSection === "agent" && <AgentPanel />}
    </AppLayout>
  );
}
