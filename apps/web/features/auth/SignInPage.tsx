"use client";
import { useState } from "react";
import { Anchor, Button, Center, Divider, Stack, Text } from "@mantine/core";
import { IconShieldCheck } from "@tabler/icons-react";
import { AuthCard } from "../../components/layout/authCard";
import { AdminPanel } from "./AdminPanel";
// import { AdminFormValues } from '../../components/form/AdminLoginForm';
// import { AdminPanel } from './AdminPanel';
import { getGetAuthOidcLoginUrl } from "../../api/komik";
import { useEffect } from "react";

type View = "oidc" | "admin";

export function SignInPage() {
  const [view, setView] = useState<View>("oidc");
  const [adminLoading, setAdminLoading] = useState(false);
  const [isOidcLoading, setIsOidcLoading] = useState(false);

  function RedirectOIDC() {
    setIsOidcLoading(true);
    window.location.href = getGetAuthOidcLoginUrl();
  }

  return (
    <AuthCard appName="Komik">
      {view === "oidc" ? (
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
            loading={isOidcLoading}
            leftSection={<IconShieldCheck size={18} />}
            onClick={() => RedirectOIDC()}
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
              onClick={() => setView("admin")}
              style={{ textUnderlineOffset: 3 }}
            >
              Use admin login
            </Anchor>
          </Center>
        </Stack>
      ) : (
        <AdminPanel onBack={() => setView("oidc")} loading={adminLoading} />
      )}
    </AuthCard>
  );
}

export default SignInPage;
