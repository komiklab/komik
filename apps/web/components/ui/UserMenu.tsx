import { Avatar, Group, Menu, Text, UnstyledButton } from "@mantine/core";
import { IconLogout, IconSettings } from "@tabler/icons-react";
import { usePostAuthLogout } from "../../api/komik";

export default function UserMenu() {
  const { mutate: logout } = usePostAuthLogout({
    mutation: {
      onSuccess: () => {
        //redirect to /
        window.location.href = "/";
      },
      onError: (error) => {
        console.log("Error logging out:", error);
      },
    },
  });
  const handleLogout = async () => {
    await logout();
  };
  return (
    <Menu shadow="lg" width={220} position="bottom-end">
      <Menu.Target>
        <UnstyledButton
          style={{
            borderRadius: 999,
            padding: 4,
          }}
        >
          <Group gap="xs" wrap="nowrap">
            <Avatar radius="xl" color="teal">
              KL
            </Avatar>
            <Text visibleFrom="sm" fz="sm" fw={600}>
              Admin
            </Text>
          </Group>
        </UnstyledButton>
      </Menu.Target>

      <Menu.Dropdown>
        <Menu.Label>Account</Menu.Label>
        <Menu.Item leftSection={<IconSettings size={16} />}>Settings</Menu.Item>
        <Menu.Divider />
        <Menu.Item
          color="red"
          leftSection={<IconLogout size={16} />}
          onClick={handleLogout}
        >
          Sign out
        </Menu.Item>
      </Menu.Dropdown>
    </Menu>
  );
}
