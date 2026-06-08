"use client";
import { useEffect } from "react";
import { useRouter } from "next/navigation";
import AdminForm from "../../components/form/adminForm";
import type { AdminFormValues } from "../../components/form/adminForm";
import { useGetAdmin, usePostAdmin } from "../../api/komik";

export default function Bootstrap() {
  const router = useRouter();
  const {
    data: adminData,
    isLoading: isCheckingAdmin,
    isError: isAdminCheckError,
    error
  } = useGetAdmin();
  const {
    mutate: createAdmin,
    isPending: isCreatingAdmin,
  } = usePostAdmin({
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

    if (!isCheckingAdmin && adminData?.status === 200 && adminData.data.exists) {
      router.replace("/signin");
    }
  }, [adminData, isAdminCheckError, isCheckingAdmin, router]);

  if (!isCheckingAdmin && adminData?.status === 200 && adminData.data.exists) {
    return null;
  }
  if (isCheckingAdmin) {
    return <div>checking</div>;
  }
  function createAdminCred(values: AdminFormValues) {
    createAdmin({
      data: {
        username: values.email,
        password: values.password,
      },
    });
  }
  return (
    <AdminForm
      submitLabel={isCreatingAdmin ? "Creating..." : "Create"}
      onSubmit={(values: AdminFormValues) => createAdminCred(values)}
    />
  );
}
