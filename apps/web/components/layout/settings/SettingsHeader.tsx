import { Box, Group, Text, ThemeIcon } from "@mantine/core";
import { IconSparkles } from "@tabler/icons-react";
import UserMenu from "../../ui/UserMenu";

export default function SettingsHeader() {
  return (
    <Group h="100%" px="lg" justify="space-between" wrap="nowrap">
      <Group gap="sm" wrap="nowrap">
        <ThemeIcon size={38} radius="md" variant="gradient" gradient={{ from: "teal", to: "blue" }}>
          <IconSparkles size={20} />
        </ThemeIcon>
        <Box>
          <Group gap="xs">
            <Text fw={800} fz="lg" lh={1}>
              KomikLab Settings
            </Text>
          </Group>
          <Text fz="xs" c="dimmed" mt={3}>
            Settings
          </Text>
        </Box>
      </Group>
    </Group>
  );
}
