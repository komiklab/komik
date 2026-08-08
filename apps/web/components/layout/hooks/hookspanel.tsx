"use client";

import {
  Stack,
  Title,
  Alert,
  Text,
  Accordion,
  Code,
  ActionIcon,
  Modal,
  Paper,
  TextInput,
  Group,
  Button,
} from "@mantine/core";
import { IconAlertCircle, IconPlus } from "@tabler/icons-react";
import { useGetHook, usePostHook } from "../../../api/komik";
import { useDisclosure } from "@mantine/hooks";
import { useForm } from "@mantine/form";
import { isValueExpired } from "next/dist/client/components/segment-cache/cache-map";

export default function HooksPanel() {
  const { data: hooksList, isLoading, isError, refetch } = useGetHook();
  const [opened, { open, close }] = useDisclosure(false);
  const hooks = hooksList?.data?.hooks;
  const {
    mutate: createHook,
    isPending,
    error: createError,
    isError: isCreateError,
    reset: resetCreate,
  } = usePostHook();
  const handleClose = () => {
    resetCreate();
    close();
  };
    const handleCreate = (payload: any) => {
    createHook({ data: payload }, {
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
        <Title order={2}>Hooks</Title>
        {isLoading && (
          <Alert icon={<IconAlertCircle size={16} />} color="blue">
            Loading hooks...
          </Alert>
        )}
        {isError && (
          <Alert icon={<IconAlertCircle size={16} />} color="red">
            Failed to load hooks.
          </Alert>
        )}
        {!isLoading && !isError && hooks.length === 0 && (
          <Text c="dimmed">No hooks found.</Text>
        )}

        <Accordion variant="separated">
          {hooks?.map((hook) => (
            <Accordion.Item key={hook.id} value={hook.id}>
              <Accordion.Control>
                <Stack gap={2}>
                  <Text fw={600}>{hook.name}</Text>

                  <Text size="sm" c="dimmed">
                    {hook.description}
                  </Text>
                </Stack>
              </Accordion.Control>

              <Accordion.Panel>
                <Code block>{JSON.stringify(hook, null, 2)}</Code>
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
      <Modal opened={opened} onClose={handleClose} title="Create Agent">
        <HooksCreateForm
          onSubmit={handleCreate}
          isLoading={isPending}
          error={isCreateError ? createError : null}
        />
      </Modal>
    </>
  );
}
const initialValues = {
  name: '',
}
function HooksCreateForm({ onSubmit, isLoading, error }: any) {
  const form = useForm({
    initialValues,
    validate: {
      name: (value:string) => (value.trim().length === 0? 'Name is required':null)
    }
  })
  const handleSubmit = (values) => {
    const payload = {
      name: values.name,
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
          <TextInput
            label="Name"
            placeholder="Enter hook name"
            {...form.getInputProps('name')}
            required
          />
          <Group justify="flex-end" mt="md">
            <Button type="submit" loading={isLoading}>
              Create
            </Button>
          </Group>
          </Stack>
          </form>
      </Paper>
  )
}
