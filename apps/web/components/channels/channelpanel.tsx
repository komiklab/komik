"use client"
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

import { useGetChannel} from "../../api/komik";
import { IconAlertCircle, IconPlus } from "@tabler/icons-react";

export default function ChannelPanel() {
    const { data: channelList, isLoading, isError, refetch } = useGetChannel();
    const channels = channelList?.data && "channels" in channelList.data ? channelList.data.channels ?? [] : [];
    return (
    <>
      <Stack gap="md">
        <Title order={2}>Channels</Title>
        {isLoading && (
          <Alert icon={<IconAlertCircle size={16} />} color="blue">
            Loading channels...
          </Alert>
        )}
        {isError && (
          <Alert icon={<IconAlertCircle size={16} />} color="red">
            Failed to load channels.
          </Alert>
        )}
        {!isLoading && !isError && channels.length === 0 && (
          <Text c="dimmed">No channels found.</Text>
        )}
        <Accordion variant="separated">
          {channels?.map((channel) => (
            <Accordion.Item key={channel.id} value={channel.id}>
              <Accordion.Control>
                <Stack gap={2}>
                  <Text fw={600}>{channel.name}</Text>

                  <Text size="sm" c="dimmed">
                    {channel.typeOf}
                  </Text>
                </Stack>
              </Accordion.Control>

              <Accordion.Panel>
                <Stack gap="sm">
                  <Code block>{JSON.stringify(channel, null, 2)}</Code>
                  <Group justify="flex-end">
                    <Button
                      leftSection={<IconSend size={16} />}
                      variant="light"
                      size="xs"
                      onClick={() => handleOpenSend(hook.name, hook.secretKey)}
                    >
                      Send Message
                    </Button>
                  </Group>
                </Stack>
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
      {/* <Modal opened={opened} onClose={handleClose} title="Create Hook">
        <HooksCreateForm
          onSubmit={handleCreate}
          isLoading={isPending}
          error={isCreateError ? createError : null}
        />
      </Modal>
      <Modal opened={sendOpened} onClose={handleSendClose} title="Send Message">
        <SendMessageForm
          onSubmit={handleSendMessage}
          isLoading={isSendPending}
          error={sendError ? { message: sendError } : null}
        />
      </Modal> */}
    </>
  );
}