import { useCallback, useSyncExternalStore } from "react"
import { TWorkspace } from "../../../types"
import { devSpaceStore } from "../devSpaceStore"

export function useWorkspaces(): readonly TWorkspace[] {
  const workspaces = useSyncExternalStore(
    useCallback((listener) => devSpaceStore.subscribe(listener), []),
    () => devSpaceStore.getAll()
  )

  return workspaces
}
