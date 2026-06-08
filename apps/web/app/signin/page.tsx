'use client'
import { useState } from 'react';
import {
  Box,
  Button,
  Center,
  Divider,
  Group,
  Paper,
  PasswordInput,
  Stack,
  Text,
  TextInput,
  ThemeIcon,
  Transition,
  Anchor,
  Badge,
  ActionIcon,
  Tooltip,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { notifications } from '@mantine/notifications';
import {
  IconBox,
  IconShieldCheck,
  IconArrowLeft,
  IconLogin,
  IconLock,
} from '@tabler/icons-react';

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------
type View = 'oidc' | 'admin';

interface AdminFormValues {
  username: string;
  password: string;
}

// ---------------------------------------------------------------------------
// Sub-components
// ---------------------------------------------------------------------------

/** Primary OIDC panel */
function OidcPanel({
  onShowAdmin,
  onOidcLogin,
  loading,
}: {
  onShowAdmin: () => void;
  onOidcLogin: () => void;
  loading: boolean;
}) {
  return (
    <Stack gap="lg">
      <Stack gap={4}>
        <Text fw={500} fz={20} c="dark">
          Welcome back
        </Text>
        <Text fz="sm" c="dimmed" lh={1.6}>
          Sign in with your organisation's identity provider to continue.
        </Text>
      </Stack>

      <Button
        fullWidth
        variant="default"
        size="md"
        leftSection={<IconShieldCheck size={18} />}
        loading={loading}
        onClick={onOidcLogin}
        styles={{
          root: {
            fontWeight: 500,
            fontSize: 14,
          },
        }}
      >
        Log in with OIDC provider
      </Button>

      <Divider label="or" labelPosition="center" />

      <Center>
        <Anchor
          component="button"
          fz="sm"
          c="dimmed"
          underline="always"
          onClick={onShowAdmin}
          style={{ textUnderlineOffset: 3 }}
        >
          Use admin login
        </Anchor>
      </Center>
    </Stack>
  );
}

/** Admin credentials panel */
function AdminPanel({
  onBack,
  onSubmit,
  loading,
}: {
  onBack: () => void;
  onSubmit: (values: AdminFormValues) => void;
  loading: boolean;
}) {
  const form = useForm<AdminFormValues>({
    initialValues: { username: '', password: '' },
    validate: {
      username: (v) =>
        v.trim().length === 0 ? 'Username is required' : null,
      password: (v) =>
        v.length === 0 ? 'Password is required' : null,
    },
  });

  return (
    <Stack gap="md">
      {/* Back navigation */}
      <Group gap={6}>
        <Tooltip label="Back to sign in" position="right" withArrow>
          <ActionIcon
            variant="subtle"
            color="gray"
            size="sm"
            onClick={onBack}
            aria-label="Back to OIDC sign in"
          >
            <IconArrowLeft size={15} />
          </ActionIcon>
        </Tooltip>
        <Text
          fz="xs"
          c="dimmed"
          style={{ cursor: 'pointer' }}
          onClick={onBack}
        >
          Back to sign in
        </Text>
      </Group>

      {/* Admin badge + heading */}
      <Stack gap={8}>
        <Badge
          color="yellow"
          variant="light"
          size="sm"
          leftSection={<IconLock size={11} />}
          style={{ width: 'fit-content' }}
        >
          Admin access
        </Badge>
        <Text fw={500} fz={20} c="dark">
          Admin login
        </Text>
        <Text fz="sm" c="dimmed">
          Restricted to system administrators only.
        </Text>
      </Stack>

      {/* Form */}
      <form onSubmit={form.onSubmit(onSubmit)} noValidate>
        <Stack gap="sm">
          <TextInput
            label="Username"
            placeholder="admin@yourapp.com"
            autoComplete="username"
            size="sm"
            {...form.getInputProps('username')}
          />

          <PasswordInput
            label="Password"
            placeholder="••••••••"
            autoComplete="current-password"
            size="sm"
            {...form.getInputProps('password')}
          />

          <Group justify="flex-end" mt={-4}>
            <Anchor fz="xs" c="dimmed" href="#" underline="hover">
              Forgot password?
            </Anchor>
          </Group>

          <Button
            type="submit"
            fullWidth
            variant="default"
            size="md"
            leftSection={<IconLogin size={16} />}
            loading={loading}
            mt="xs"
            styles={{ root: { fontWeight: 500, fontSize: 14 } }}
          >
            Sign in as admin
          </Button>
        </Stack>
      </form>
    </Stack>
  );
}

// ---------------------------------------------------------------------------
// Main component
// ---------------------------------------------------------------------------

export function SignInPage() {
  const [view, setView] = useState<View>('oidc');
  const [oidcLoading, setOidcLoading] = useState(false);
  const [adminLoading, setAdminLoading] = useState(false);

  const handleOidcLogin = () => {
    setOidcLoading(true);
    // TODO: replace with your real OIDC redirect
    setTimeout(() => {
      notifications.show({
        title: 'Redirecting',
        message: 'Taking you to your identity provider…',
        color: 'blue',
      });
      setOidcLoading(false);
    }, 1500);
  };

  const handleAdminSubmit = (values: AdminFormValues) => {
    setAdminLoading(true);
    // TODO: replace with your real admin auth call
    setTimeout(() => {
      notifications.show({
        title: 'Welcome, admin',
        message: `Signed in as ${values.username}`,
        color: 'green',
      });
      setAdminLoading(false);
    }, 1500);
  };

  return (
    <Box
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: 'var(--mantine-color-gray-0)',
        padding: '1rem',
      }}
    >
      <Paper
        withBorder
        shadow="xs"
        p="xl"
        radius="md"
        style={{ width: '100%', maxWidth: 400 }}
      >
        {/* Logo / App name */}
        <Group gap={10} mb="xl">
          <ThemeIcon variant="light" color="gray" size={32} radius="md">
            <IconBox size={17} />
          </ThemeIcon>
          <Text fw={500} fz={15}>
            Komik
          </Text>
        </Group>

        {/* OIDC view */}
        <Transition
          mounted={view === 'oidc'}
          transition="fade"
          duration={180}
          timingFunction="ease"
        >
          {(styles) => (
            <div style={styles}>
              <OidcPanel
                onShowAdmin={() => setView('admin')}
                onOidcLogin={handleOidcLogin}
                loading={oidcLoading}
              />
            </div>
          )}
        </Transition>

        {/* Admin view */}
        <Transition
          mounted={view === 'admin'}
          transition="fade"
          duration={180}
          timingFunction="ease"
        >
          {(styles) => (
            <div style={styles}>
              <AdminPanel
                onBack={() => setView('oidc')}
                onSubmit={handleAdminSubmit}
                loading={adminLoading}
              />
            </div>
          )}
        </Transition>
      </Paper>
    </Box>
  );
}

export default SignInPage;