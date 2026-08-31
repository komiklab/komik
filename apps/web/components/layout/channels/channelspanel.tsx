"use client";

import { useState } from "react";
import {
  Accordion,
  ActionIcon,
  Alert,
  Button,
  Code,
  Group,
  Modal,
  Paper,
  PasswordInput,
  Select,
  Stack,
  Text,
  TextInput,
  Title,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { useForm } from "@mantine/form";
import { IconAlertCircle, IconPlus, IconSend } from "@tabler/icons-react";
import {
  postChannelChannelIdSend,
  useGetChannel,
  usePostChannel,
} from "../../../api/komik";
import type { ChannelRequest } from "../../../api/schemas";

type ChannelFormValues = {
  typeOf: "slack" | null;
  name: string;
  slackChannelName: string;
  botToken: string;
};

type ListedChannel = {
  id?: string;
  name?: string;
  typeOf?: string;
  type?: string;
  [key: string]: unknown;
};

const initialValues: ChannelFormValues = {
  typeOf: null,
  name: "",
  slackChannelName: "",
  botToken: "",
};

function getErrorMessage(error: unknown, fallback: string) {
  return error instanceof Error && error.message ? error.message : fallback;
}

export default function ChannelsPanel() {
  const { data: channelList, isLoading, isError, refetch } = useGetChannel();
  const { mutate: createChannel, isPending, error: createError, reset } =
    usePostChannel();
  const [opened, { open, close }] = useDisclosure(false);
  const [testingChannelId, setTestingChannelId] = useState<string | null>(null);
  const [testError, setTestError] = useState<string | null>(null);
  const channels = (channelList?.data?.channels ?? []) as ListedChannel[];

  const handleClose = () => {
    reset();
    close();
  };

  const handleCreate = (values: ChannelFormValues) => {
    if (values.typeOf !== "slack") return;

    const payload: ChannelRequest = {
      name: values.name.trim(),
      typeOf: values.typeOf,
      payload: {
        channelName: values.slackChannelName.trim(),
        botToken: values.botToken,
      },
    };

    createChannel(
      { data: payload },
      {
        onSuccess: () => {
          handleClose();
          refetch();
        },
      },
    );
  };

  const handleTest = async (channelId?: string) => {
    if (!channelId) {
      setTestError("This channel does not have an ID to test.");
      return;
    }

    setTestError(null);
    setTestingChannelId(channelId);

    try {
      await postChannelChannelIdSend(channelId);
    } catch (error) {
      setTestError(getErrorMessage(error, "Failed to send the test message."));
    } finally {
      setTestingChannelId(null);
    }
  };

  return (
    <>
      <Stack gap="md">
        <Title order={2}>Channels</Title>
        {isLoading && <Alert color="blue">Loading channels...</Alert>}
        {isError && <Alert color="red">Failed to load channels.</Alert>}
        {!isLoading && !isError && channels.length === 0 && (
          <Text c="dimmed">No channels found.</Text>
        )}

        <Accordion variant="separated">
          {channels.map((channel) => {
            const channelId = channel.id;
            const type = channel.typeOf ?? channel.type;

            return (
              <Accordion.Item
                key={channelId ?? channel.name}
                value={channelId ?? channel.name ?? "channel"}
              >
                <Accordion.Control>
                  <Stack gap={2}>
                    <Text fw={600}>{channel.name}</Text>
                    {type && <Text size="sm" c="dimmed">{type}</Text>}
                  </Stack>
                </Accordion.Control>
                <Accordion.Panel>
                  <Stack gap="sm">
                    <Code block>{JSON.stringify(channel, null, 2)}</Code>
                    <Group justify="flex-end">
                      <Button
                        leftSection={<IconSend size={16} />}
                        loading={testingChannelId === (channelId ?? "")}
                        size="xs"
                        variant="light"
                        onClick={() => handleTest(channelId)}
                      >
                        Test
                      </Button>
                    </Group>
                  </Stack>
                </Accordion.Panel>
              </Accordion.Item>
            );
          })}
        </Accordion>

        {testError && (
          <Alert icon={<IconAlertCircle size={16} />} color="red" title="Test failed">
            {testError}
          </Alert>
        )}
      </Stack>

      <ActionIcon
        size={60}
        radius="xl"
        variant="filled"
        color="blue"
        onClick={open}
        style={{ position: "fixed", bottom: 24, right: 24, zIndex: 1000 }}
      >
        <IconPlus size={28} />
      </ActionIcon>

      <Modal opened={opened} onClose={handleClose} title="Create Channel">
        <ChannelCreateForm error={createError} isLoading={isPending} onSubmit={handleCreate} />
      </Modal>
    </>
  );
}

function ChannelCreateForm({
  onSubmit,
  isLoading,
  error,
}: {
  onSubmit: (values: ChannelFormValues) => void;
  isLoading: boolean;
  error: unknown;
}) {
  const form = useForm<ChannelFormValues>({
    initialValues,
    validate: {
      typeOf: (value) => (value ? null : "Channel type is required"),
      name: (value) => (value.trim() ? null : "Channel name is required"),
      slackChannelName: (value, values) =>
        values.typeOf === "slack" && !value.trim()
          ? "Slack channel name is required"
          : null,
      botToken: (value, values) =>
        values.typeOf === "slack" && !value.trim() ? "Bot token is required" : null,
    },
  });

  return (
    <Paper withBorder shadow="sm" p="lg" radius="md">
      <form onSubmit={form.onSubmit(onSubmit)}>
        <Stack gap="lg">
          {Boolean(error) && (
            <Alert icon={<IconAlertCircle size={16} />} color="red" title="Error">
              {getErrorMessage(error, "An error occurred while creating the channel.")}
            </Alert>
          )}
          <Select
            label="Channel type"
            placeholder="Select a channel type"
            data={[{ value: "slack", label: "Slack" }]}
            {...form.getInputProps("typeOf")}
            required
          />

          {form.values.typeOf === "slack" && (
            <>
              <TextInput
                label="Channel name"
                placeholder="e.g. #alerts"
                {...form.getInputProps("slackChannelName")}
                required
              />
              <PasswordInput
                label="Bot token"
                placeholder="xoxb-..."
                {...form.getInputProps("botToken")}
                required
              />
              <TextInput
                label="Name"
                description="A name to identify this channel in Komik."
                placeholder="e.g. Production alerts"
                {...form.getInputProps("name")}
                required
              />
            </>
          )}

          <Group justify="flex-end">
            <Button type="submit" loading={isLoading}>Create</Button>
          </Group>
        </Stack>
      </form>
    </Paper>
  );
}
