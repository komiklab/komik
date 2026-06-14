import {
  TextInput,
  PasswordInput,
  Button,
  Paper,
  Title,
  Container,
  Text,
} from "@mantine/core";
import { useForm } from "@mantine/form";

interface AdminFormProps {
  onSubmit: (values: AdminFormValues) => void;
  submitLabel: string;
  error?: string | null;
  loading?: boolean;
}

export type AdminFormValues = {
  username: string;
  password: string;
};

export default function AdminForm({ onSubmit, submitLabel, error, loading }: AdminFormProps) {
  const form = useForm({
    mode: "uncontrolled",
    initialValues: {
      username: "",
      password: "",
    },
    validate: {
      username: (value) => (value.length > 0 ? null : "Username is required"),
      password: (value) =>
        value.length < 6 ? "Password must be at least 6 characters" : null,
    },
  });

  return (

        <form onSubmit={form.onSubmit(onSubmit)}>
          {error && (
            <Text c="red" size="sm" mb="sm" fw={500}>
              {error}
            </Text>
          )}
          <TextInput
            label="Username"
            placeholder="Your username"
            key={form.key("username")}
            {...form.getInputProps("username")}
            required
          />
          <PasswordInput
            label="Password"
            placeholder="Your password"
            mt="md"
            key={form.key("password")}
            {...form.getInputProps("password")}
            required
          />
          <Button fullWidth mt="xl" type="submit" loading={loading}>
            {submitLabel}
          </Button>
        </form>

  );
}
