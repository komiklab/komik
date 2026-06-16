"use client";
import { useState } from "react";
import { Anchor, Button, Center, Divider, Stack, Text } from "@mantine/core";
import { IconShieldCheck } from "@tabler/icons-react";
import { AuthCard } from "../../components/layout/authCard";
import { AdminPanel } from "./AdminPanel";
// import { AdminFormValues } from '../../components/form/AdminLoginForm';
// import { AdminPanel } from './AdminPanel';
import { getGetAuthOidcLoginUrl } from "../../api/komik";

type View = "oidc" | "admin";

export function SignInPage() {
  const [view, setView] = useState<View>("oidc");
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
            <Text fw={800} fz={24}>
              Welcome back
            </Text>
            <Text fz="sm" c="dimmed" lh={1.6}>
              Sign in with your organisation's identity provider to continue.
            </Text>
          </Stack>

          <Button
            fullWidth
            variant="gradient"
            gradient={{ from: "teal", to: "blue" }}
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
              c="teal.3"
              underline="always"
              onClick={() => setView("admin")}
              style={{ textUnderlineOffset: 3 }}
            >
              Use admin login
            </Anchor>
          </Center>
        </Stack>
      ) : (
        <AdminPanel onBack={() => setView("oidc")} />
      )}
    </AuthCard>
  );
}

export default SignInPage;
