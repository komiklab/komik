import { Badge, Box, Group, Text, ThemeIcon } from "@mantine/core";
import { IconSparkles } from "@tabler/icons-react";
import UserMenu from "../ui/UserMenu";

export default function AppHeader() {
  return (
    <Group h="100%" px="lg" justify="space-between" wrap="nowrap">
      <Group gap="sm" wrap="nowrap">
        <ThemeIcon size={38} radius="md" variant="gradient" gradient={{ from: "teal", to: "blue" }}>
          <IconSparkles size={20} />
        </ThemeIcon>
        <Box>
          <Group gap="xs">
            <Text fw={800} fz="lg" lh={1}>
              KomikLab
            </Text>
            <Badge size="sm" variant="light" color="teal">
              Admin
            </Badge>
          </Group>
          <Text fz="xs" c="dimmed" mt={3}>
            Operations console
          </Text>
        </Box>
      </Group>
      <UserMenu />
    </Group>
  );
}
