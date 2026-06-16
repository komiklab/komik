"use client";
import { useEffect } from "react";
import { useRouter } from "next/navigation";
import AdminForm from "../../components/form/adminForm";
import type { AdminFormValues } from "../../components/form/adminForm";
import { useGetAdmin, usePostAdmin } from "../../api/komik";
import { Loader, Stack, Text, ThemeIcon } from "@mantine/core";
import { IconUserShield } from "@tabler/icons-react";
import { AuthCard } from "../../components/layout/authCard";

export default function Bootstrap() {
  const router = useRouter();
  const {
    data: adminData,
    isLoading: isCheckingAdmin,
    isError: isAdminCheckError,
  } = useGetAdmin();
  const { mutate: createAdmin, isPending: isCreatingAdmin } = usePostAdmin({
    mutation: {
      onSuccess: () => {
        router.replace("/signin");
      },
      onError: (error: any) => {
        console.log("Error creating admin", error);
        router.replace("/error");
      },
    },
  });

  useEffect(() => {
    if (isAdminCheckError) {
      router.replace("/error");
      return;
    }

    if (
      !isCheckingAdmin &&
      adminData?.status === 200 &&
      adminData.data.exists
    ) {
      router.replace("/");
    }
  }, [adminData, isAdminCheckError, isCheckingAdmin, router]);

  if (!isCheckingAdmin && adminData?.status === 200 && adminData.data.exists) {
    return null;
  }
  if (isCheckingAdmin) {
    return (
      <AuthCard appName="KomikLab">
        <Stack align="center" py="lg">
          <Loader color="teal" />
          <Text c="dimmed" fz="sm">
            Checking administrator setup...
          </Text>
        </Stack>
      </AuthCard>
    );
  }
  function createAdminCred(values: AdminFormValues) {
    createAdmin({
      data: {
        username: values.username,
        password: values.password,
      },
    });
  }
  return (
    <AuthCard appName="KomikLab">
      <Stack gap="lg">
        <Stack gap="sm">
          <ThemeIcon
            size={44}
            radius="md"
            variant="light"
            color="teal"
          >
            <IconUserShield size={22} />
          </ThemeIcon>
          <div>
            <Text fw={800} fz={24}>
              Create administrator
            </Text>
            <Text c="dimmed" fz="sm" mt={4} lh={1.6}>
              Set up the first local administrator account for this KomikLab instance.
            </Text>
          </div>
        </Stack>
        <AdminForm
          submitLabel={isCreatingAdmin ? "Creating..." : "Create"}
          onSubmit={(values: AdminFormValues) => createAdminCred(values)}
          loading={isCreatingAdmin}
        />
      </Stack>
    </AuthCard>
  );
}
