"use client";

import {
  Badge,
  Button,
  Group,
  Paper,
  Stack,
  Table,
  Text,
  ThemeIcon,
  Title,
} from "@mantine/core";
import { IconApi, IconInbox, IconRefresh } from "@tabler/icons-react";

const mockInboxRows = [
  {
    id: "mock-inbox-001",
    source: "Slack",
    subject: "New workspace request from product ops",
    owner: "Unassigned",
    priority: "Normal",
    receivedAt: "Awaiting API",
  },
];

export function InboxPanel() {
  return (
    <Stack gap="lg">
      <Group justify="space-between" align="flex-start">
        <Group gap="md">
          <ThemeIcon size={48} radius="md" variant="gradient" gradient={{ from: "teal", to: "cyan" }}>
            <IconInbox size={24} />
          </ThemeIcon>
          <div>
            <Group gap="xs">
              <Title order={1} fz={{ base: 26, md: 34 }}>
                Inbox
              </Title>
              <Badge color="gray" variant="light">
                Mock data
              </Badge>
            </Group>
            <Text c="dimmed" mt={4}>
              A focused queue for inbound work once the inbox endpoint is available.
            </Text>
          </div>
        </Group>

        <Button leftSection={<IconRefresh size={16} />} variant="light" disabled>
          Sync pending
        </Button>
      </Group>

      <Paper withBorder p="lg" bg="rgba(255,255,255,0.04)">
        <Group gap="md" align="flex-start">
          <ThemeIcon color="blue" variant="light" size={40} radius="md">
            <IconApi size={20} />
          </ThemeIcon>
          <div>
            <Text fw={700}>API contract needed</Text>
            <Text fz="sm" c="dimmed" lh={1.6}>
              Expected later: backend should expose an inbox list endpoint with pagination
              and Orval should generate a query hook. This view can then replace the
              single mock row with generated backend data.
            </Text>
          </div>
        </Group>
      </Paper>

      <Paper withBorder bg="rgba(255,255,255,0.035)">
        <Table.ScrollContainer minWidth={760}>
          <Table verticalSpacing="md" horizontalSpacing="lg" highlightOnHover>
            <Table.Thead>
              <Table.Tr>
                <Table.Th>Source</Table.Th>
                <Table.Th>Subject</Table.Th>
                <Table.Th>Owner</Table.Th>
                <Table.Th>Priority</Table.Th>
                <Table.Th>Received</Table.Th>
              </Table.Tr>
            </Table.Thead>
            <Table.Tbody>
              {mockInboxRows.map((row) => (
                <Table.Tr key={row.id}>
                  <Table.Td>
                    <Badge variant="light" color="violet">
                      {row.source}
                    </Badge>
                  </Table.Td>
                  <Table.Td>
                    <Text fw={700}>{row.subject}</Text>
                    <Text fz="xs" c="dimmed">
                      {row.id}
                    </Text>
                  </Table.Td>
                  <Table.Td>{row.owner}</Table.Td>
                  <Table.Td>
                    <Badge variant="outline" color="teal">
                      {row.priority}
                    </Badge>
                  </Table.Td>
                  <Table.Td c="dimmed">{row.receivedAt}</Table.Td>
                </Table.Tr>
              ))}
            </Table.Tbody>
          </Table>
        </Table.ScrollContainer>
      </Paper>
    </Stack>
  );
}
