"use client";

import { useState } from "react";
import {
  ActionIcon,
  Badge,
  Button,
  Code,
  Group,
  Loader,
  NumberFormatter,
  Pagination,
  Paper,
  Select,
  Stack,
  Table,
  Text,
  ThemeIcon,
  Title,
  Tooltip,
} from "@mantine/core";
import {
  IconAlertTriangle,
  IconClock,
  IconDatabaseSearch,
  IconRefresh,
  IconShieldCheck,
} from "@tabler/icons-react";
import { useGetAuditlog } from "../../../api/komik";
import type { AuditlogResponse } from "../../../api/schemas";

const PAGE_SIZE_OPTIONS = ["10", "20", "50"];

function formatDate(value?: string) {
  if (!value) return "Unknown";

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;

  return new Intl.DateTimeFormat("en", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(date);
}

function severityColor(severity?: string) {
  const normalized = severity?.toLowerCase();

  if (normalized === "critical" || normalized === "error") return "red";
  if (normalized === "warning" || normalized === "warn") return "yellow";
  if (normalized === "info") return "blue";

  return "teal";
}

function compactId(value?: string) {
  if (!value) return "Not available";
  if (value.length <= 16) return value;
  return `${value.slice(0, 8)}...${value.slice(-6)}`;
}

function AuditRows({ items }: { items: AuditlogResponse[] }) {
  return (
    <>
      {items.map((item, index) => (
        <Table.Tr key={item.event_id ?? `${item.occurred_at}-${index}`}>
          <Table.Td>
            <Stack gap={2}>
              <Group gap="xs" wrap="nowrap">
                <ThemeIcon size={28} radius="md" variant="light" color={severityColor(item.severity)}>
                  <IconShieldCheck size={15} />
                </ThemeIcon>
                <Text fw={700}>{item.resource_type ?? "System event"}</Text>
              </Group>
              <Tooltip label={item.event_id ?? "No event id"}>
                <Text fz="xs" c="dimmed">
                  {compactId(item.event_id)}
                </Text>
              </Tooltip>
            </Stack>
          </Table.Td>
          <Table.Td>
            <Badge color={severityColor(item.severity)} variant="light">
              {item.severity ?? "normal"}
            </Badge>
          </Table.Td>
          <Table.Td>
            <Text fw={600}>{item.initiator_type ?? "Unknown"}</Text>
            <Text fz="xs" c="dimmed">
              {item.initiator_id ?? "No initiator id"}
            </Text>
          </Table.Td>
          <Table.Td>
            <Group gap={6} wrap="nowrap">
              <IconClock size={15} color="var(--mantine-color-dimmed)" />
              <Text fz="sm">{formatDate(item.occurred_at)}</Text>
            </Group>
          </Table.Td>
          <Table.Td>
            <Tooltip label={item.correlation_id ?? "No correlation id"}>
              <Code>{compactId(item.correlation_id)}</Code>
            </Tooltip>
          </Table.Td>
          <Table.Td>
            <Text fz="sm" lineClamp={2} maw={360}>
              {item.data ?? "No payload"}
            </Text>
          </Table.Td>
        </Table.Tr>
      ))}
    </>
  );
}

export function AuditLogDashboard() {
  const [activePage, setActivePage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const offset = (activePage - 1) * pageSize;

  const auditLogQuery = useGetAuditlog(
    { limit: pageSize, offset },
    {
      query: {
        placeholderData: (previousData) => previousData,
      },
    },
  );

  const response = auditLogQuery.data?.status === 200 ? auditLogQuery.data.data : undefined;
  const items = response?.items ?? [];
  const metadata = response?.metadata;
  const total = metadata?.total ?? items.length;
  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  const firstTimestamp = items[0]?.occurred_at;
  const latestEvent = firstTimestamp ? formatDate(firstTimestamp) : "No events yet";

  return (
    <Stack gap="lg">
      <Group justify="space-between" align="flex-start">
        <Group gap="md">
          <ThemeIcon size={48} radius="md" variant="gradient" gradient={{ from: "teal", to: "blue" }}>
            <IconDatabaseSearch size={24} />
          </ThemeIcon>
          <div>
            <Title order={1} fz={{ base: 26, md: 34 }}>
              Audit Log
            </Title>
            <Text c="dimmed" mt={4}>
              Trace security-sensitive actions across KomikLab.
            </Text>
          </div>
        </Group>

        <Group gap="xs">
          <Select
            w={112}
            aria-label="Rows per page"
            data={PAGE_SIZE_OPTIONS}
            value={String(pageSize)}
            onChange={(value) => {
              setPageSize(Number(value ?? 10));
              setActivePage(1);
            }}
          />
          <Tooltip label="Refresh audit log">
            <ActionIcon
              size={36}
              variant="light"
              color="teal"
              loading={auditLogQuery.isFetching}
              onClick={() => auditLogQuery.refetch()}
            >
              <IconRefresh size={18} />
            </ActionIcon>
          </Tooltip>
        </Group>
      </Group>

      <Group grow align="stretch">
        <Paper withBorder p="lg" bg="rgba(255,255,255,0.045)">
          <Text fz="xs" c="dimmed" tt="uppercase" fw={800}>
            Total events
          </Text>
          <Text fz={30} fw={800} mt={4}>
            <NumberFormatter value={total} thousandSeparator />
          </Text>
        </Paper>
        <Paper withBorder p="lg" bg="rgba(255,255,255,0.045)">
          <Text fz="xs" c="dimmed" tt="uppercase" fw={800}>
            Current page
          </Text>
          <Text fz={30} fw={800} mt={4}>
            {activePage}
          </Text>
        </Paper>
        <Paper withBorder p="lg" bg="rgba(255,255,255,0.045)">
          <Text fz="xs" c="dimmed" tt="uppercase" fw={800}>
            Latest event
          </Text>
          <Text fz="sm" fw={700} mt={11}>
            {latestEvent}
          </Text>
        </Paper>
      </Group>

      <Paper withBorder bg="rgba(255,255,255,0.035)">
        {auditLogQuery.isLoading ? (
          <Stack align="center" py={72}>
            <Loader color="teal" />
            <Text c="dimmed">Loading audit trail...</Text>
          </Stack>
        ) : auditLogQuery.isError ? (
          <Stack align="center" py={72} px="lg">
            <ThemeIcon color="red" variant="light" size={48} radius="md">
              <IconAlertTriangle size={24} />
            </ThemeIcon>
            <Text fw={800}>Audit log could not be loaded</Text>
            <Text c="dimmed" ta="center">
              The generated audit log query returned an error. Check the API server
              and authentication state, then try again.
            </Text>
            <Button leftSection={<IconRefresh size={16} />} onClick={() => auditLogQuery.refetch()}>
              Retry
            </Button>
          </Stack>
        ) : items.length === 0 ? (
          <Stack align="center" py={72} px="lg">
            <ThemeIcon color="teal" variant="light" size={48} radius="md">
              <IconDatabaseSearch size={24} />
            </ThemeIcon>
            <Text fw={800}>No audit events found</Text>
            <Text c="dimmed" ta="center">
              This page is connected to the backend, but the current pagination
              window did not return any events.
            </Text>
          </Stack>
        ) : (
          <>
            <Table.ScrollContainer minWidth={1040}>
              <Table verticalSpacing="md" horizontalSpacing="lg" highlightOnHover>
                <Table.Thead>
                  <Table.Tr>
                    <Table.Th>Event</Table.Th>
                    <Table.Th>Severity</Table.Th>
                    <Table.Th>Initiator</Table.Th>
                    <Table.Th>Occurred</Table.Th>
                    <Table.Th>Correlation</Table.Th>
                    <Table.Th>Payload</Table.Th>
                  </Table.Tr>
                </Table.Thead>
                <Table.Tbody>
                  <AuditRows items={items} />
                </Table.Tbody>
              </Table>
            </Table.ScrollContainer>

            <Group justify="space-between" px="lg" py="md">
              <Text fz="sm" c="dimmed">
                Showing {offset + 1}-{Math.min(offset + items.length, total)} of{" "}
                <NumberFormatter value={total} thousandSeparator /> events
              </Text>
              <Pagination value={activePage} onChange={setActivePage} total={totalPages} color="teal" />
            </Group>
          </>
        )}
      </Paper>
    </Stack>
  );
}
