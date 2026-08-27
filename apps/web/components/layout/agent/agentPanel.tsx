"use client";
import { useState } from "react";
import {
  Accordion,
  ActionIcon,
  Alert,
  Code,
  Modal,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { IconAlertCircle, IconPlus } from "@tabler/icons-react";
import { useGetAgent, usePostAgent } from "../../../api/komik";
import { useDisclosure } from "@mantine/hooks";
import { AgentCreateForm } from "../../form/agentCreateForm";

export default function AgentPanel() {
  const { data: agentsList, isLoading, isError, refetch } = useGetAgent();
  const agents = agentsList?.data && "agents" in agentsList.data ? agentsList.data.agents ?? [] : [];
  const [opened, { open, close }] = useDisclosure(false);
  const {
    mutate: createAgent,
    isPending,
    error: createError,
    isError: isCreateError,
    reset: resetCreate,
  } = usePostAgent();

  const handleClose = () => {
    resetCreate();
    close();
  };

  const handleCreate = (payload: any) => {
    createAgent(
      { data: payload },
      {
        onSuccess: () => {
          close();
          refetch();
        },
        onError: (error) => {
          console.error("Failed to create agent", error);
        },
      },
    );
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
          {agents?.map((agent: any) => (
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
      <Modal opened={opened} onClose={handleClose} title="Create Agent" size="xl">
        {/* Replace with your form */}
        <AgentCreateForm
          onSubmit={handleCreate}
          isLoading={isPending}
          error={isCreateError ? createError : null}
        />
      </Modal>
    </>
  );
}


