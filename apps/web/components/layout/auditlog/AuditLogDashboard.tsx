"use client";

import { useMemo, useState } from "react";
import {
  ActionIcon,
  Badge,
  Code,
  Group,
  Loader,
  NumberFormatter,
  Paper,
  Stack,
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
import { MantineReactTable, useMantineReactTable } from "mantine-react-table";
import type { MRT_ColumnDef, MRT_PaginationState } from "mantine-react-table";
import { useGetAuditlog } from "../../../api/komik";
import type { AuditlogResponse } from "../../../api/schemas";

// ─── helpers ────────────────────────────────────────────────────────────────

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
  const n = severity?.toLowerCase();
  if (n === "critical" || n === "error") return "red";
  if (n === "warning" || n === "warn") return "yellow";
  if (n === "info") return "blue";
  return "teal";
}

function compactId(value?: string) {
  if (!value) return "—";
  if (value.length <= 16) return value;
  return `${value.slice(0, 8)}…${value.slice(-6)}`;
}

// ─── column definitions ─────────────────────────────────────────────────────

const columns: MRT_ColumnDef<AuditlogResponse>[] = [
  {
    accessorKey: "resource_type",
    header: "Event",
    size: 220,
    Cell: ({ row }) => {
      const item = row.original;
      return (
        <Stack gap={2}>
          <Group gap="xs" wrap="nowrap">
            <ThemeIcon
              size={28}
              radius="md"
              variant="light"
              color={severityColor(item.severity)}
            >
              <IconShieldCheck size={15} />
            </ThemeIcon>
            <Text fw={700} fz="sm">
              {item.event_type ?? "System event"}
            </Text>
          </Group>
          <Tooltip label={item.event_id ?? "No event id"} withinPortal>
            <Text fz="xs" c="dimmed" style={{ cursor: "default" }}>
              {compactId(item.event_id)}
            </Text>
          </Tooltip>
        </Stack>
      );
    },
  },
  {
    accessorKey: "severity",
    header: "Severity",
    size: 110,
    Cell: ({ cell }) => {
      const val = cell.getValue<string | undefined>();
      return (
        <Badge color={severityColor(val)} variant="light">
          {val ?? "normal"}
        </Badge>
      );
    },
  },
  {
    accessorKey: "initiator_type",
    header: "Initiator",
    size: 160,
    Cell: ({ row }) => {
      const item = row.original;
      return (
        <>
          <Text fw={600} fz="sm">
            {item.initiator_type ?? "Unknown"}
          </Text>
          <Text fz="xs" c="dimmed">
            {item.initiator_id ?? "—"}
          </Text>
        </>
      );
    },
  },
  {
    accessorKey: "occurred_at",
    header: "Occurred",
    size: 180,
    Cell: ({ cell }) => {
      const val = cell.getValue<string | undefined>();
      return (
        <Group gap={6} wrap="nowrap">
          <IconClock size={15} color="var(--mantine-color-dimmed)" />
          <Text fz="sm">{formatDate(val)}</Text>
        </Group>
      );
    },
  },
  {
    accessorKey: "correlation_id",
    header: "Correlation",
    size: 160,
    Cell: ({ cell }) => {
      const val = cell.getValue<string | undefined>();
      return (
        <Tooltip label={val ?? "No correlation id"} withinPortal>
          <Code>{compactId(val)}</Code>
        </Tooltip>
      );
    },
  },
  {
    accessorKey: "data",
    header: "Payload",
    size: 320,
    Cell: ({ cell }) => {
      const val = cell.getValue<string | undefined>();
      return (
        <Text fz="sm" lineClamp={2} maw={320}>
          {val ?? "No payload"}
        </Text>
      );
    },
  },
];

// ─── component ───────────────────────────────────────────────────────────────

export function AuditLogDashboard() {
  const [pagination, setPagination] = useState<MRT_PaginationState>({
    pageIndex: 0,
    pageSize: 20,
  });

  const offset = pagination.pageIndex * pagination.pageSize;

  const auditLogQuery = useGetAuditlog(
    { limit: pagination.pageSize, offset },
    {
      query: {
        placeholderData: (prev) => prev,
      },
    },
  );

  const response =
    auditLogQuery.data?.status === 200 ? auditLogQuery.data.data : undefined;
  const items: AuditlogResponse[] = response?.items ?? [];
  const metadata = response?.metadata;
  const total = metadata?.total ?? 0;

  const firstTimestamp = items[0]?.occurred_at;
  const latestEvent = firstTimestamp ? formatDate(firstTimestamp) : "No events yet";

  // keep referential stability for MRT
  const data = useMemo(() => items, [items]);

  const table = useMantineReactTable({
    columns,
    data,
    // server-side pagination
    manualPagination: true,
    rowCount: total,
    onPaginationChange: setPagination,
    state: {
      pagination,
      isLoading: auditLogQuery.isLoading,
      showAlertBanner: auditLogQuery.isError,
      showProgressBars: auditLogQuery.isFetching && !auditLogQuery.isLoading,
    },
    // MRT display options
    enableSorting: false,
    enableFilters: false,
    enableColumnActions: false,
    enableFullScreenToggle: false,
    enableDensityToggle: false,
    enableHiding: false,
    enableGlobalFilter: false,
    // styling
    mantineTableProps: {
      highlightOnHover: true,
      withTableBorder: false,
      withColumnBorders: false,
    },
    mantinePaperProps: {
      withBorder: true,
      bg: "rgba(255,255,255,0.035)",
    },
    mantineToolbarAlertBannerProps: auditLogQuery.isError
      ? {
          color: "red",
          children: "Audit log could not be loaded. Check the API server and authentication state.",
        }
      : undefined,
    // pagination labels
    localization: {
      rowsPerPage: "Rows per page",
    } as never,
    // page size options
    paginationDisplayMode: "pages",
    mantinePaginationProps: {
      rowsPerPageOptions: ["10", "20", "50"],
      showRowsPerPage: true,
      color: "teal",
    },
    // custom no-data state
    renderEmptyRowsFallback: () => (
      <Stack align="center" py={64} gap="xs">
        <ThemeIcon color="teal" variant="light" size={48} radius="md">
          <IconDatabaseSearch size={24} />
        </ThemeIcon>
        <Text fw={800}>No audit events found</Text>
        <Text c="dimmed" ta="center" fz="sm">
          This page is connected to the backend, but the current pagination
          window did not return any events.
        </Text>
      </Stack>
    ),
  });

  return (
    <Stack gap="lg">
      {/* ── Header ── */}
      <Group justify="space-between" align="flex-start">
        <Group gap="md">
          <ThemeIcon
            size={48}
            radius="md"
            variant="gradient"
            gradient={{ from: "teal", to: "blue" }}
          >
            <IconDatabaseSearch size={24} />
          </ThemeIcon>
          <div>
            <Title order={1} style={{ fontSize: "clamp(1.5rem, 3vw, 2.125rem)" }}>
              Audit Log
            </Title>
            <Text c="dimmed" mt={4}>
              Trace security-sensitive actions across KomikLab.
            </Text>
          </div>
        </Group>

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

      {/* ── Stat cards ── */}
      <Group grow align="stretch">
        <Paper withBorder p="lg" bg="rgba(255,255,255,0.045)">
          <Text fz="xs" c="dimmed" tt="uppercase" fw={800}>
            Total events
          </Text>
          <Text fz={30} fw={800} mt={4}>
            {auditLogQuery.isLoading ? (
              <Loader size="sm" color="teal" />
            ) : (
              <NumberFormatter value={total} thousandSeparator />
            )}
          </Text>
        </Paper>

        <Paper withBorder p="lg" bg="rgba(255,255,255,0.045)">
          <Text fz="xs" c="dimmed" tt="uppercase" fw={800}>
            Current page
          </Text>
          <Text fz={30} fw={800} mt={4}>
            {pagination.pageIndex + 1}
          </Text>
        </Paper>

        <Paper withBorder p="lg" bg="rgba(255,255,255,0.045)">
          <Text fz="xs" c="dimmed" tt="uppercase" fw={800}>
            Latest event
          </Text>
          <Text fz="sm" fw={700} mt={11}>
            {auditLogQuery.isLoading ? <Loader size="sm" color="teal" /> : latestEvent}
          </Text>
        </Paper>
      </Group>

      {/* ── Table ── */}
      <MantineReactTable table={table} />
    </Stack>
  );
}
