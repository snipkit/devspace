import { useContext } from "react"
import { TProInstanceManager } from "../../types"
import { DevSpaceContext, TDevspaceContext } from "./DevSpaceProvider"
import { useProInstanceManager } from "./useProInstanceManager"

export function useProInstances(): [TDevspaceContext["proInstances"], TProInstanceManager] {
  const proInstances = useContext(DevSpaceContext).proInstances
  const manager = useProInstanceManager()

  return [proInstances, manager]
}
