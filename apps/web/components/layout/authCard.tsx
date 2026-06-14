"use client";
import { Anchor, Box, Group, Paper, Text, ThemeIcon } from "@mantine/core";
import { IconBox } from "@tabler/icons-react";
import { ReactNode } from "react";

interface AuthCardProps {
  /** App name shown next to the logo */
  appName?: string;
  /** Slot for the auth panel content */
  children: ReactNode;
}

export function AuthCard({ appName = "YourApp", children }: AuthCardProps) {
  return (
    <Box
      style={{
        minHeight: "100vh",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        backgroundColor: "var(--mantine-color-gray-0)",
        padding: "1rem",
      }}
    >
      <Paper
        withBorder
        shadow="xs"
        p="xl"
        radius="md"
        style={{ width: "100%", maxWidth: 400 }}
      >
        {/* Logo / App name */}
        <Group gap={10} mb="xl">
          <ThemeIcon variant="light" color="gray" size={32} radius="md">
            <IconBox size={17} />
          </ThemeIcon>
          <Text fw={500} fz={15}>
            {appName}
          </Text>
        </Group>

        {/* Panel content injected here */}
        {children}

        {/* Footer */}
        <Text fz="xs" c="dimmed" ta="center" mt="lg" lh={1.6}>
          By signing in you agree to our{" "}
          <Anchor fz="xs" href="#" underline="hover">
            Terms of Service
          </Anchor>{" "}
          and{" "}
          <Anchor fz="xs" href="#" underline="hover">
            Privacy Policy
          </Anchor>
          .
        </Text>
      </Paper>
    </Box>
  );
}
