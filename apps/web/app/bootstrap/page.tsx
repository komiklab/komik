"use client"
import { TextInput, PasswordInput, Button, Paper, Title, Container } from '@mantine/core';
import { useForm } from '@mantine/form';
import { useAdminStore } from '../../stores/admin';
import { useRouter } from 'next/navigation';

export default function Bootstrap() {
  const router = useRouter()
  function createAdminCred() {
    console.log("creating admin cred")
    useAdminStore.setState({ doesAdminExist: true })
    router.replace("/signin")
  }
  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      email: '',
      password: '',
    },
    validate: {
      email: (value) => (/^\S+@\S+$/.test(value) ? null : 'Invalid email'),
      password: (value) => (value.length < 6 ? 'Password must be at least 6 characters' : null),
    },
  });

  return (
    <Container size={420} my={40}>
      <Title ta="center" order={2}>
        Create Administrator
      </Title>
      <Paper withBorder shadow="md" p={30} mt={30} radius="md">
        <form onSubmit={form.onSubmit((values) => console.log(values))}>
          <TextInput
            label="Email"
            placeholder="you@mantine.dev"
            key={form.key('email')}
            {...form.getInputProps('email')}
            required
          />
          <PasswordInput
            label="Password"
            placeholder="Your password"
            mt="md"
            key={form.key('password')}
            {...form.getInputProps('password')}
            required
          />
          <Button fullWidth mt="xl" type="submit" onClick={createAdminCred}>
            Create
          </Button>
        </form>
      </Paper>
    </Container>
  );
}
