import {
  TextInput,
  PasswordInput,
  Button,
  Paper,
  Title,
  Container,
} from "@mantine/core";
import { useForm } from "@mantine/form";

interface AdminFormProps {
  onSubmit: (values: AdminFormValues) => void;
  submitLabel: string;
}

export type AdminFormValues = {
  email: string;
  password: string;
};

export default function AdminForm({ onSubmit, submitLabel }: AdminFormProps) {
  const form = useForm({
    mode: "uncontrolled",
    initialValues: {
      email: "",
      password: "",
    },
    validate: {
      email: (value) => (/^\S+@\S+$/.test(value) ? null : "Invalid email"),
      password: (value) =>
        value.length < 6 ? "Password must be at least 6 characters" : null,
    },
  });

  return (
    <Container size={420} my={40}>
      <Title ta="center" order={2}>
        Create Administrator
      </Title>
      <Paper withBorder shadow="md" p={30} mt={30} radius="md">
        <form onSubmit={form.onSubmit(onSubmit)}>
          <TextInput
            label="Email"
            placeholder="you@mantine.dev"
            key={form.key("email")}
            {...form.getInputProps("email")}
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
          <Button fullWidth mt="xl" type="submit" >
            {submitLabel}
          </Button>
        </form>
      </Paper>
    </Container>
  );
}
