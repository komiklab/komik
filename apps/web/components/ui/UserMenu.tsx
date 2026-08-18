import { Avatar, Group, Menu, UnstyledButton } from "@mantine/core";
import { IconLogout, IconSettings } from "@tabler/icons-react";
import { usePostAuthLogout, useGetAuthMe } from "../../api/komik";
import { useRouter } from "next/navigation";
import { useUiStore } from "../../stores/ui";

export default function UserMenu() {
  const router = useRouter();
  const setCurrentSection = useUiStore((state) => state.setCurrentSection);
  const redirectSettings = () => {
    setCurrentSection("agent");
    router.push("/settings");
  };
  const { data: authdata } = useGetAuthMe();
  const username =
    authdata?.status === 200 ? authdata.data.username : undefined;
  console.log("username is " + username);
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
    logout();
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
              {username?.toUpperCase()[0]}
            </Avatar>
          </Group>
        </UnstyledButton>
      </Menu.Target>

      <Menu.Dropdown>
        <Menu.Label>{username}</Menu.Label>
        <Menu.Item leftSection={<IconSettings size={16} />} onClick={redirectSettings}>Settings</Menu.Item>
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
