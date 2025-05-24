import { useContext } from "react"
import { TProviderManager } from "../../types"
import { DevSpaceContext, TDevspaceContext } from "./DevSpaceProvider"
import { useProviderManager } from "./useProviderManager"

export function useProviders(): [TDevspaceContext["providers"], TProviderManager] {
  const providers = useContext(DevSpaceContext).providers
  const manager = useProviderManager()

  return [providers, manager]
}
