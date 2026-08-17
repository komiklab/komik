"use client";

import { Alert, Button, Group, Paper, Stack, TextInput, JsonInput } from "@mantine/core";
import { useForm } from "@mantine/form";
import { IconAlertCircle } from "@tabler/icons-react";

export default function SendMessageForm({
  onSubmit,
  isLoading,
  error,
}: {
  onSubmit: (payload: any) => void;
  isLoading: boolean;
  error: Error | null;
}) {
  const form = useForm({
    initialValues: {
      reference: "",
      message: "",
    },
    validate: {
      reference: (value: string) => (value.trim().length === 0 ? "Reference is required" : null),
      message: (value: string) => {
        if (value.trim().length === 0) return "Message is required";
        try {
          JSON.parse(value);
          return null;
        } catch (e) {
          return "Invalid JSON format";
        }
      },
    },
  });
  const handleSubmit = (values: any) => {
    onSubmit({
      ...values,
      message: JSON.parse(values.message),
    });
  };
  return (
    <Paper withBorder shadow="sm" p="lg" radius="md" maw={560} mx="auto">
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <Stack gap="lg">
          {error && (
            <Alert icon={<IconAlertCircle size={16} />} color="red" title="Error">
              {error?.message || "An error occurred while creating the Hook."}
            </Alert>
          )}
          <TextInput
            label="Reference"
            placeholder="unique reference"
            {...form.getInputProps('reference')}
            required
          />
          <JsonInput
            label="Message"
            placeholder='{"key": "value"}'
            formatOnBlur
            autosize
            minRows={4}
            {...form.getInputProps('message')}
            required
          />
          <Group justify="flex-end" mt="md">
            <Button type="submit" loading={isLoading}>
              Send
            </Button>
          </Group>
        </Stack>
      </form>
    </Paper>
  );
}