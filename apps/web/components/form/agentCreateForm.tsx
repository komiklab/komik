import { useState } from "react";
import {
  ActionIcon,
  Alert,
  Box,
  Button,
  Collapse,
  Divider,
  Group,
  Paper,
  Stack,
  Text,
  Textarea,
  TextInput,
  Checkbox,
  TagsInput,
} from "@mantine/core";
import {
  IconAlertCircle,
  IconPlus,
  IconTrash,
  IconChevronDown,
  IconChevronUp,
} from "@tabler/icons-react";
import { useForm } from "@mantine/form";

// ---- default values, mirrors the OpenAPI examples ----
const initialValues = {
  name: "",
  endpoint: "",
  description: "",
  triggered_by: [],
  parameters: [],
  tags: [],
};

// Small reusable component for editing parameters
function ParameterList({ label, addLabel, values, onAdd, onRemove, onChange }: any) {
  return (
    <Box>
      <Group justify="space-between" mb="xs">
        <Text size="sm" fw={500}>
          {label}
        </Text>
        <Button
          size="xs"
          variant="light"
          leftSection={<IconPlus size={14} />}
          onClick={onAdd}
        >
          {addLabel}
        </Button>
      </Group>
      {values.length === 0 && (
        <Text size="xs" c="dimmed" mb="xs">
          None added
        </Text>
      )}
      <Stack gap="xs">
        {values.map((row: any, index: number) => (
          <Group key={index} gap="xs" wrap="nowrap" align="center">
            <TextInput
              placeholder="Name"
              value={row.name}
              onChange={(e) => onChange(index, "name", e.currentTarget.value)}
              style={{ flex: 1 }}
              size="sm"
            />
            <TextInput
              placeholder="Type (e.g. string)"
              value={row.type}
              onChange={(e) => onChange(index, "type", e.currentTarget.value)}
              style={{ flex: 1 }}
              size="sm"
            />
            <TextInput
              placeholder="Description"
              value={row.description}
              onChange={(e) =>
                onChange(index, "description", e.currentTarget.value)
              }
              style={{ flex: 2 }}
              size="sm"
            />
            <Checkbox
              label="Required"
              checked={row.required}
              onChange={(e) =>
                onChange(index, "required", e.currentTarget.checked)
              }
              size="sm"
            />
            <ActionIcon
              color="red"
              variant="subtle"
              onClick={() => onRemove(index)}
            >
              <IconTrash size={16} />
            </ActionIcon>
          </Group>
        ))}
      </Stack>
    </Box>
  );
}

export function AgentCreateForm({ onSubmit, isLoading, error }: any) {
  const [advancedOpen, setAdvancedOpen] = useState(false);

  const form = useForm({
    initialValues,
    validate: {
      name: (value: string) => (value.trim().length === 0 ? "Name is required" : null),
      endpoint: (value: string) => {
        if (value.trim().length === 0) return "Endpoint is required";
        try {
          const url = new URL(value);
          if (url.protocol !== "http:" && url.protocol !== "https:") {
            return "Endpoint must be a valid HTTP/HTTPS URL";
          }
        } catch {
          return "Endpoint must be a valid URL";
        }
        return null;
      },
      description: (value: string) =>
        value.trim().length === 0 ? "Description is required" : null,
    },
  });

  // Array field helpers
  const parametersHelpers = {
    add: () =>
      form.insertListItem("parameters", {
        name: "",
        type: "",
        description: "",
        required: false,
      }),
    remove: (index: number) => form.removeListItem("parameters", index),
    change: (index: number, key: string, value: any) =>
      form.setFieldValue(`parameters.${index}.${key}`, value),
  };

  const handleSubmit = (values: typeof initialValues) => {
    // filter out empty values before submitting
    const payload = {
      ...values,
      triggered_by: values.triggered_by.filter((c: string) => c.trim() !== ""),
      tags: values.tags.filter((t: string) => t.trim() !== ""),
      parameters: values.parameters.filter(
        (p: any) => p.type.trim() !== "" || p.description.trim() !== ""
      ),
    };
    onSubmit ? onSubmit(payload) : console.log(payload);
  };

  return (
    <Paper withBorder shadow="sm" p="lg" radius="md" maw={800} mx="auto">
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <Stack gap="lg">
          {error && (
            <Alert
              icon={<IconAlertCircle size={16} />}
              color="red"
              title="Error"
            >
              {error?.message || "An error occurred while creating the agent."}
            </Alert>
          )}

          <Stack gap="sm">
            <TextInput
              label="Name"
              placeholder="my-agent"
              required
              {...form.getInputProps("name")}
            />
            <TextInput
              label="Endpoint"
              placeholder="http://localhost:8080"
              required
              {...form.getInputProps("endpoint")}
            />
          </Stack>
          <Textarea
            label="Description"
            placeholder="What this agent does"
            autosize
            minRows={2}
            required
            {...form.getInputProps("description")}
          />

          <Divider
            label={
              <Group
                gap={4}
                style={{ cursor: "pointer" }}
                onClick={() => setAdvancedOpen((o) => !o)}
              >
                <Text size="sm" fw={600}>
                  Triggered by & Parameters
                </Text>
                {advancedOpen ? (
                  <IconChevronUp size={14} />
                ) : (
                  <IconChevronDown size={14} />
                )}
              </Group>
            }
            labelPosition="left"
          />

          <Collapse expanded={advancedOpen}>
            <Stack gap="md">
              <TagsInput
                label="Triggered By"
                placeholder="Press Enter to add capability"
                required
                {...form.getInputProps("triggered_by")}
              />

              <ParameterList
                label="Parameters"
                addLabel="Add parameter"
                values={form.values.parameters}
                onAdd={parametersHelpers.add}
                onRemove={parametersHelpers.remove}
                onChange={parametersHelpers.change}
              />

              <TagsInput
                label="Tags"
                placeholder="Press Enter to add tag"
                {...form.getInputProps("tags")}
              />
            </Stack>
          </Collapse>

          <Group justify="flex-end">
            <Button type="submit" loading={isLoading}>
              Create agent
            </Button>
          </Group>
        </Stack>
      </form>
    </Paper>
  );
}
