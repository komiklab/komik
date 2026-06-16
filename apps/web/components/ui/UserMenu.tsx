import { Menu, Avatar, Text } from "@mantine/core";

export default function UserMenu() {
  return (
    <Menu shadow="md" width={200}>
      <Menu.Target>
         <Avatar radius="xl" />
      </Menu.Target>

      <Menu.Dropdown>
        <Menu.Label>Application</Menu.Label>
        <Menu.Item>Settings</Menu.Item>
        <Menu.Item>Signout</Menu.Item>
      </Menu.Dropdown>
    </Menu>
  );
}
