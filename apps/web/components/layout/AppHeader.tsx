import { Group, Text } from "@mantine/core";
import UserMenu from "../ui/UserMenu";

export default function AppHeader(){
    return(
        <Group justify="space-between">
            <Text>KomikLab</Text>
            <UserMenu/>
        </Group>
    )
}