"use client";
import { ActionIcon, Badge, Group, Stack, Text, Tooltip } from "@mantine/core";
import { IconArrowLeft, IconLock } from "@tabler/icons-react";
// import { AdminLoginForm, AdminFormValues } from '../../components/form/AdminLoginForm';
import { useState } from "react";
import AdminForm, { AdminFormValues } from "../../components/form/adminForm";
import { usePostAuthLogin } from "../../api/komik";
import { useRouter } from "next/navigation";

interface AdminPanelProps {
  onBack: () => void;
}

export function AdminPanel({ onBack }: AdminPanelProps) {
  const router = useRouter();
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const { mutate: AuthLoginAdmin, isPending: isAdminLoading } =
    usePostAuthLogin({
      mutation: {
        onSuccess: () => {
          setErrorMsg(null);
          router.replace("/");
        },
        onError: (error: any) => {
          const backenderror = error.response?.data
          console.log("data is " + JSON.stringify(backenderror))
          // if (error.status === 401) {
          //   setErrorMsg("Invalid email or password");
          // } else {
          //   setErrorMsg("An error occurred during sign in");
          // }
          setErrorMsg(backenderror.error_message)
          console.log("Error signing in", error);
        },
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
        <Text fz="xs" c="dimmed" style={{ cursor: "pointer" }} onClick={onBack}>
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
          style={{ width: "fit-content" }}
        >
          Admin access
        </Badge>
        <Text fw={800} fz={24}>
          Admin login
        </Text>
        <Text fz="sm" c="dimmed">
          Restricted to system administrators only.
        </Text>
      </Stack>

      {/* Delegate form rendering to the form component */}
      <AdminForm
        onSubmit={(values: AdminFormValues) => AuthLoginAdmin({ data: values })}
        submitLabel="Login"
        error={errorMsg}
        loading={isAdminLoading}
      />
    </Stack>
  );
}
