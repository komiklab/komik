"use client";
import { Anchor, Box, Group, Paper, Text, ThemeIcon } from "@mantine/core";
import { IconSparkles } from "@tabler/icons-react";
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
        background:
          "radial-gradient(circle at 20% 10%, rgba(32, 201, 151, 0.22), transparent 32%), linear-gradient(135deg, var(--mantine-color-dark-9), #121b24 48%, #101521)",
        padding: "1.25rem",
      }}
    >
      <Paper
        withBorder
        shadow="xl"
        p={{ base: "lg", sm: "xl" }}
        radius="md"
        style={{
          width: "100%",
          maxWidth: 430,
          background: "rgba(12, 18, 27, 0.82)",
          borderColor: "rgba(255,255,255,0.1)",
          backdropFilter: "blur(18px)",
        }}
      >
        {/* Logo / App name */}
        <Group gap={10} mb="xl">
          <ThemeIcon
            variant="gradient"
            gradient={{ from: "teal", to: "blue" }}
            size={36}
            radius="md"
          >
            <IconSparkles size={18} />
          </ThemeIcon>
          <div>
            <Text fw={800} fz={16} lh={1}>
              {appName}
            </Text>
            <Text fz="xs" c="dimmed" mt={4}>
              Operations console
            </Text>
          </div>
        </Group>

        {/* Panel content injected here */}
        {children}

        {/* Footer */}
        <Text fz="xs" c="dimmed" ta="center" mt="lg" lh={1.6}>
          By signing in you agree to our{" "}
          <Anchor fz="xs" href="#" underline="hover" c="teal.3">
            Terms of Service
          </Anchor>{" "}
          and{" "}
          <Anchor fz="xs" href="#" underline="hover" c="teal.3">
            Privacy Policy
          </Anchor>
          .
        </Text>
      </Paper>
    </Box>
  );
}
