import { useCallback, useSyncExternalStore } from "react"
import { TActionObj } from "../action"
import { devSpaceStore } from "../devSpaceStore"

export function useAllWorkspaceActions() {
  const actions = useSyncExternalStore(
    useCallback((listener) => devSpaceStore.subscribe(listener), []),
    () => devSpaceStore.getAllActions()
  )

  return { active: actions.active, history: actions.history.slice().sort(sortByCreationDesc) }
}

function sortByCreationDesc(a: TActionObj, b: TActionObj) {
  return b.createdAt - a.createdAt
}
