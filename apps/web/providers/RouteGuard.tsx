"use client";

import { useEffect } from "react";
import { useRouter, usePathname } from "next/navigation";
import { useGetAdmin, useGetAuthMe } from "../api/komik";

export function RouteGuard({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();

  const { data: adminData, isLoading: isAdminLoading, isError: isAdminError } = useGetAdmin();
  const adminExists = adminData?.status === 200 && adminData.data.exists;

  // Only fetch session data if the admin actually exists
  const { data: sessionData, isLoading: isSessionLoading, isError: isSessionError } = useGetAuthMe({
    query: {
      enabled: !!adminExists,
      retry: false, // Do not retry if unauthorized
    }
  });

  useEffect(() => {
    // 1. Wait for admin check to finish
    if (isAdminLoading) return;

    if (isAdminError) {
      if (pathname !== "/error") router.replace("/error");
      return;
    }

    // 2. Admin doesn't exist -> go to bootstrap
    if (!adminExists) {
      if (pathname !== "/bootstrap") router.replace("/bootstrap");
      return;
    }

    // 3. Admin exists. Now wait for session check to finish
    if (isSessionLoading) return;

    // 4. Admin exists, session is loaded -> Route accordingly
    if (sessionData) {
      // If they are logged in but on the signin page, push them to home
      if (pathname === "/signin") {
        router.replace("/");
      }
    } else if (isSessionError || !sessionData) {
      // Not logged in -> go to signin
      if (pathname !== "/signin") {
        router.replace("/signin");
      }
    }

  }, [isAdminLoading, isAdminError, adminExists, isSessionLoading, sessionData, isSessionError, router, pathname]);

  // Show loading state while checking either admin or active session
  if (isAdminLoading || (adminExists && isSessionLoading)) {
    return <div>Loading...</div>;
  }

  // Once all checks pass and routing is stable, render children
  return <>{children}</>;
}
