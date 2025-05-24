import { Link, Text } from "@chakra-ui/react"
import { client } from "../client"
import { Khulnasoft } from "../icons"

export function KhulnasoftOSSBadge() {
    return (
        <Link
            display="flex"
            alignItems="center"
            justifyContent="start"
            onClick={() => client.open("https://khulnasoft.com/")}>
            <Text fontSize="sm" variant="muted" marginRight="2">
                Open sourced by
            </Text>
            <Khulnasoft width="10" height="6" />
        </Link>
    )
}
