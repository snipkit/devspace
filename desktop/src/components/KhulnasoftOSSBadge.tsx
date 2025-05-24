import { Link, Text, useColorModeValue } from "@chakra-ui/react"
import { client } from "../client"
import { Khulnasoft } from "../icons"

export function KhulnasoftOSSBadge() {
  const textColor = useColorModeValue("gray.500", "gray.400")

  return (
    <Link
      display="flex"
      alignItems="center"
      justifyContent="start"
      onClick={() => client.open("https://khulnasoft.sh/")}>
      <Text fontSize="sm" color={textColor} marginRight="2">
        Open sourced by
      </Text>
      <Khulnasoft width="10" height="6" />
    </Link>
  )
}
