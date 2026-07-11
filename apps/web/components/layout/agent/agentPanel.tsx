"use client";
import { useState } from "react";
import {
  Accordion,
  ActionIcon,
  Alert,
  Box,
  Button,
  Checkbox,
  Code,
  Collapse,
  Divider,
  Group,
  Modal,
  NumberInput,
  Paper,
  Select,
  Stack,
  Text,
  Textarea,
  TextInput,
  Title,
} from "@mantine/core";
import { IconAlertCircle, IconPlus } from "@tabler/icons-react";
import { useGetAgent, usePostAgent } from "../../../api/komik";
import { useDisclosure } from "@mantine/hooks";
import { useForm } from '@mantine/form';
import {  IconTrash, IconChevronDown, IconChevronUp } from '@tabler/icons-react';

export default function AgentPanel() {
  const { data: agentsList, isLoading, isError, refetch } = useGetAgent();
  const agents = agentsList?.data?.agents;
  const [opened, { open, close }] = useDisclosure(false);
  const { mutate: createAgent, isPending, error: createError, isError: isCreateError, reset: resetCreate } = usePostAgent();

  const handleClose = () => {
    resetCreate();
    close();
  };

  const handleCreate = (payload: any) => {
    createAgent({ data: payload }, {
      onSuccess: () => {
        close();
        refetch();
      },
      onError: (error) => {
        console.error("Failed to create agent", error);
      }
    });
  };

  return (
    <>
      <Stack gap="md">
        <Title order={2}>Agents</Title>

        {isLoading && (
          <Alert icon={<IconAlertCircle size={16} />} color="blue">
            Loading agents...
          </Alert>
        )}

        {isError && (
          <Alert icon={<IconAlertCircle size={16} />} color="red">
            Failed to load agents.
          </Alert>
        )}

        {!isLoading && !isError && agents.length === 0 && (
          <Text c="dimmed">No agents found.</Text>
        )}

        <Accordion variant="separated">
          {agents?.map((agent) => (
            <Accordion.Item key={agent.id} value={agent.id}>
              <Accordion.Control>
                <Stack gap={2}>
                  <Text fw={600}>{agent.name}</Text>

                  <Text size="sm" c="dimmed">
                    {agent.description}
                  </Text>
                </Stack>
              </Accordion.Control>

              <Accordion.Panel>
                <Code block>{JSON.stringify(agent, null, 2)}</Code>
              </Accordion.Panel>
            </Accordion.Item>
          ))}
        </Accordion>
      </Stack>

      {/* Floating Action Button */}
      <ActionIcon
        size={60}
        radius="xl"
        variant="filled"
        color="blue"
        onClick={open}
        style={{
          position: "fixed",
          bottom: 24,
          right: 24,
          zIndex: 1000,
        }}
      >
        <IconPlus size={28} />
      </ActionIcon>

      {/* Create Agent Modal */}
      <Modal opened={opened} onClose={handleClose} title="Create Agent">
        {/* Replace with your form */}
        <AgentCreateForm onSubmit={handleCreate} isLoading={isPending} error={isCreateError ? createError : null} />
      </Modal>
    </>
  );
}




// ---- default values, mirrors the OpenAPI examples ----
const initialValues = {
  name: '',
  image: '',
  description: '',
  imagePullSecret: '',
  resources: {
    memory: '256Mi',
    cpu: '512m',
    ephemeralStorage: '1Gi',
    timeoutSeconds: 3600,
  },
  env: [],
  secrets: [],
  annotations: [],
};

// Small reusable component for editing an array of { key1, key2 } rows
function KeyValueList({ label, addLabel, keys, values, onAdd, onRemove, onChange, placeholder }) {
  return (
    <Box>
      <Group justify="space-between" mb="xs">
        <Text size="sm" fw={500}>{label}</Text>
        <Button size="xs" variant="light" leftSection={<IconPlus size={14} />} onClick={onAdd}>
          {addLabel}
        </Button>
      </Group>
      {values.length === 0 && (
        <Text size="xs" c="dimmed" mb="xs">None added</Text>
      )}
      <Stack gap="xs">
        {values.map((row, index) => (
          <Group key={index} gap="xs" wrap="nowrap" align="flex-end">
            <TextInput
              placeholder={placeholder[0]}
              value={row[keys[0]]}
              onChange={(e) => onChange(index, keys[0], e.currentTarget.value)}
              style={{ flex: 1 }}
              size="sm"
            />
            <TextInput
              placeholder={placeholder[1]}
              value={row[keys[1]]}
              onChange={(e) => onChange(index, keys[1], e.currentTarget.value)}
              style={{ flex: 1 }}
              size="sm"
            />
            <ActionIcon color="red" variant="subtle" onClick={() => onRemove(index)}>
              <IconTrash size={16} />
            </ActionIcon>
          </Group>
        ))}
      </Stack>
    </Box>
  );
}

function AgentCreateForm({ onSubmit, isLoading, error }: any) {
  const [advancedOpen, setAdvancedOpen] = useState(false);

  const form = useForm({
    initialValues,
    validate: {
      name: (value) => (value.trim().length === 0 ? 'Name is required' : null),
      image: (value) => (value.trim().length === 0 ? 'Image is required' : null),
    },
  });

  // generic helpers for the three array fields (env, secrets, annotations)
  const makeArrayHelpers = (field, emptyRow) => ({
    add: () => form.insertListItem(field, { ...emptyRow }),
    remove: (index) => form.removeListItem(field, index),
    change: (index, key, value) => form.setFieldValue(`${field}.${index}.${key}`, value),
  });

  const envHelpers = makeArrayHelpers('env', { name: '', value: '' });
  const secretsHelpers = makeArrayHelpers('secrets', { name: '', secretKeyRef: '' });
  const annotationsHelpers = makeArrayHelpers('annotations', { key: '', value: '' });

  const handleSubmit = (values) => {
    // strip empty optional strings so we don't send blank fields to the API
    const payload = {
      name: values.name,
      image: values.image,
      ...(values.description && { description: values.description }),
      ...(values.imagePullSecret && { imagePullSecret: values.imagePullSecret }),
      resources: values.resources,
      ...(values.env.length && { env: values.env }),
      ...(values.secrets.length && { secrets: values.secrets }),
      ...(values.annotations.length && { annotations: values.annotations }),
    };
    onSubmit ? onSubmit(payload) : console.log(payload);
  };

  return (
    <Paper withBorder shadow="sm" p="lg" radius="md" maw={560} mx="auto">
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <Stack gap="lg">
          {error && (
            <Alert icon={<IconAlertCircle size={16} />} color="red" title="Error">
              {error?.message || "An error occurred while creating the agent."}
            </Alert>
          )}

          {/* <div>
            <Title order={3}>Create agent</Title>
            <Text size="sm" c="dimmed">Fields marked required must be filled in.</Text>
          </div> */}

          {/* ---------------- Required section ---------------- */}
          <Stack gap="sm">
            {/* <Text size="sm" fw={600}>Required</Text> */}
            <TextInput
              label="Name"
              placeholder="my-agent"
              required
              {...form.getInputProps('name')}
            />
            <TextInput
              label="Image"
              placeholder="ghcr.io/komiklab/komikagent:v1"
              required
              {...form.getInputProps('image')}
            />
          </Stack>

          <Divider
            label={
              <Group gap={4} style={{ cursor: 'pointer' }} onClick={() => setAdvancedOpen((o) => !o)}>
                <Text size="sm" fw={600}>Advanced settings</Text>
                {advancedOpen ? <IconChevronUp size={14} /> : <IconChevronDown size={14} />}
              </Group>
            }
            labelPosition="left"
          />

          {/* ---------------- Optional / default-value section ---------------- */}
          <Collapse expanded={advancedOpen}>
            <Stack gap="md">
              <Textarea
                label="Description"
                placeholder="What this agent does"
                autosize
                minRows={2}
                {...form.getInputProps('description')}
              />
              <TextInput
                label="Image pull secret"
                placeholder="my-secret"
                {...form.getInputProps('imagePullSecret')}
              />

              <Box>
                <Text size="sm" fw={500} mb="xs">Resources</Text>
                <Stack gap="xs">
                  <TextInput
                    label="Memory"
                    size="sm"
                    {...form.getInputProps('resources.memory')}
                  />
                  <TextInput
                    label="CPU"
                    size="sm"
                    {...form.getInputProps('resources.cpu')}
                  />
                  <TextInput
                    label="Ephemeral storage"
                    size="sm"
                    {...form.getInputProps('resources.ephemeralStorage')}
                  />
                  <NumberInput
                    label="Timeout (seconds)"
                    size="sm"
                    min={0}
                    {...form.getInputProps('resources.timeoutSeconds')}
                  />
                </Stack>
              </Box>

              <KeyValueList
                label="Environment variables"
                addLabel="Add env var"
                keys={['name', 'value']}
                placeholder={['NAME', 'value']}
                values={form.values.env}
                onAdd={envHelpers.add}
                onRemove={envHelpers.remove}
                onChange={envHelpers.change}
              />

              <KeyValueList
                label="Secrets"
                addLabel="Add secret"
                keys={['name', 'secretKeyRef']}
                placeholder={['name', 'secretKeyRef']}
                values={form.values.secrets}
                onAdd={secretsHelpers.add}
                onRemove={secretsHelpers.remove}
                onChange={secretsHelpers.change}
              />

              <KeyValueList
                label="Annotations"
                addLabel="Add annotation"
                keys={['key', 'value']}
                placeholder={['key', 'value']}
                values={form.values.annotations}
                onAdd={annotationsHelpers.add}
                onRemove={annotationsHelpers.remove}
                onChange={annotationsHelpers.change}
              />
            </Stack>
          </Collapse>

          <Group justify="flex-end">
            <Button type="submit" loading={isLoading}>Create agent</Button>
          </Group>
        </Stack>
      </form>
    </Paper>
  );
}