import { useCallback, useId, useMemo, useRef, useSyncExternalStore } from "react"
import { client, TStreamEventListenerFn } from "../../../client"
import { exists } from "../../../lib"
import {
  TDeepNonNullable,
  TStreamID,
  TWorkspace,
  TWorkspaceID,
  TWorkspaceStartConfig,
} from "../../../types"
import { TActionID, TActionObj, useConnectAction, useReplayAction } from "../action"
import { devSpaceStore } from "../devSpaceStore"

export type TWorkspaceResult = Readonly<{
  data: TWorkspace | undefined
  isLoading: boolean
  current:
    | (TActionObj & Readonly<{ connect: (listener: TStreamEventListenerFn) => void }>)
    | undefined
  history: Readonly<{
    // all: readonly TActionObj[]
    replay: (actionID: TActionID, listener: TStreamEventListenerFn) => void
  }>
  start: (config: TWorkspaceStartConfig, onStream?: TStreamEventListenerFn) => TActionID | undefined
  create: (
    config: Omit<TWorkspaceStartConfig, "sourceConfig"> &
      Pick<TDeepNonNullable<TWorkspaceStartConfig>, "sourceConfig">,
    onStream?: TStreamEventListenerFn
  ) => TActionID
  stop: (onStream?: TStreamEventListenerFn) => TActionID | undefined
  remove: (force: boolean, onStream?: TStreamEventListenerFn) => TActionID | undefined
  rebuild: (onStream?: TStreamEventListenerFn) => TActionID | undefined
  reset: (onStream?: TStreamEventListenerFn) => TActionID | undefined
  checkStatus: (onStream?: TStreamEventListenerFn) => TActionID | undefined
}>

export function useWorkspaceActions(
  workspaceID: TWorkspaceID | undefined
): TActionObj[] | undefined {
  const dataCache = useRef<TActionObj[]>()
  const data = useSyncExternalStore(
    useCallback((listener) => devSpaceStore.subscribe(listener), []),
    () => {
      if (workspaceID === undefined) {
        return undefined
      }

      const workspaceActions = devSpaceStore.getWorkspaceActions(workspaceID)
      if (!dataCache.current || dataCache.current.length !== workspaceActions.length) {
        dataCache.current = workspaceActions

        return dataCache.current
      }

      // compare actions
      const diff = dataCache.current.filter(
        (action) => !workspaceActions.find((workspaceAction) => action.id === workspaceAction.id)
      )
      if (diff.length > 0) {
        dataCache.current = workspaceActions

        return dataCache.current
      }

      return dataCache.current
    }
  )

  return data
}

export function useWorkspace(workspaceID: TWorkspaceID | undefined): TWorkspaceResult {
  const viewID = useId()
  const data = useSyncExternalStore(
    useCallback((listener) => devSpaceStore.subscribe(listener), []),
    () => (workspaceID !== undefined ? devSpaceStore.get(workspaceID) : undefined)
  )
  const create = useCallback<TWorkspaceResult["create"]>(
    (config, onStream) => {
      return devSpaceStore.startAction({
        actionName: "start",
        workspaceID: config.id,
        actionFn: async (ctx) => {
          const result = await client.workspaces.start(config, onStream, {
            id: config.id,
            actionID: ctx.id,
            streamID: viewID,
          })
          if (result.err) {
            return result
          }
          devSpaceStore.setStatus(config.id, result.val)

          return result
        },
      })
    },
    [viewID]
  )

  const start = useCallback<TWorkspaceResult["start"]>(
    (config, onStream) => {
      if (workspaceID === undefined) {
        return
      }

      return startWorkspaceAction({ workspaceID, config, onStream, streamID: viewID })
    },
    [viewID, workspaceID]
  )

  const checkStatus = useCallback<TWorkspaceResult["checkStatus"]>(
    (onStream) => {
      if (workspaceID === undefined) {
        return
      }

      return devSpaceStore.startAction({
        actionName: "checkStatus",
        workspaceID,
        actionFn: async (ctx) => {
          const result = await client.workspaces.checkStatus(onStream, {
            id: workspaceID,
            actionID: ctx.id,
            streamID: viewID,
          })
          if (result.err) {
            return result
          }
          devSpaceStore.setStatus(workspaceID, result.val)

          return result
        },
      })
    },
    [viewID, workspaceID]
  )

  const stop = useCallback<TWorkspaceResult["stop"]>(
    (onStream) => {
      if (workspaceID === undefined) {
        return
      }

      return stopWorkspaceAction({ workspaceID, onStream, streamID: viewID })
    },
    [viewID, workspaceID]
  )

  const rebuild = useCallback<TWorkspaceResult["rebuild"]>(
    (onStream) => {
      if (workspaceID === undefined) {
        return
      }

      return devSpaceStore.startAction({
        actionName: "rebuild",
        workspaceID,
        actionFn: async (ctx) => {
          const result = await client.workspaces.rebuild(onStream, {
            id: workspaceID,
            actionID: ctx.id,
            streamID: viewID,
          })
          if (result.err) {
            return result
          }
          devSpaceStore.setStatus(workspaceID, result.val)

          return result
        },
      })
    },
    [viewID, workspaceID]
  )

  const reset = useCallback<TWorkspaceResult["reset"]>(
    (onStream) => {
      if (workspaceID === undefined) {
        return
      }

      return devSpaceStore.startAction({
        actionName: "reset",
        workspaceID,
        actionFn: async (ctx) => {
          const result = await client.workspaces.reset(onStream, {
            id: workspaceID,
            actionID: ctx.id,
            streamID: viewID,
          })
          if (result.err) {
            return result
          }
          devSpaceStore.setStatus(workspaceID, result.val)

          return result
        },
      })
    },
    [viewID, workspaceID]
  )

  const remove = useCallback<TWorkspaceResult["remove"]>(
    (force, onStream) => {
      if (workspaceID === undefined) {
        return
      }

      return removeWorkspaceAction({ force, workspaceID, onStream, streamID: viewID })
    },
    [viewID, workspaceID]
  )

  const currentAction = useSyncExternalStore(
    useCallback((listener) => devSpaceStore.subscribe(listener), []),
    () => (workspaceID !== undefined ? devSpaceStore.getCurrentAction(workspaceID) : undefined)
  )
  const isLoading = useMemo(() => exists(currentAction), [currentAction])

  const connect = useConnectAction(currentAction, viewID)
  const current = useMemo<TWorkspaceResult["current"]>(() => {
    if (currentAction === undefined) {
      return undefined
    }

    return {
      ...currentAction,
      connect,
    }
  }, [currentAction, connect])

  const replay = useReplayAction()
  const history = useMemo<TWorkspaceResult["history"]>(() => {
    return { replay }
  }, [replay])

  return useMemo(
    () => ({
      data,
      isLoading,
      current,
      history,
      create,
      start,
      stop,
      rebuild,
      reset,
      remove,
      checkStatus,
    }),
    [data, isLoading, current, history, create, start, stop, rebuild, reset, remove, checkStatus]
  )
}

type TStartWorkspaceActionArgs = Readonly<{
  config: TWorkspaceStartConfig
  onStream?: TStreamEventListenerFn
  workspaceID: TWorkspaceID
  streamID: TStreamID
}>
export function startWorkspaceAction({
  workspaceID,
  streamID,
  config,
  onStream,
}: TStartWorkspaceActionArgs): TActionObj["id"] {
  return devSpaceStore.startAction({
    actionName: "start",
    workspaceID,
    actionFn: async (ctx) => {
      const result = await client.workspaces.start(config, onStream, {
        id: workspaceID,
        actionID: ctx.id,
        streamID,
      })
      if (result.err) {
        return result
      }
      devSpaceStore.setStatus(workspaceID, result.val)

      return result
    },
  })
}

type TStopWorkspaceActionArgs = Readonly<{
  onStream?: TStreamEventListenerFn
  workspaceID: TWorkspaceID
  streamID: TStreamID
}>
export function stopWorkspaceAction({
  workspaceID,
  onStream,
  streamID,
}: TStopWorkspaceActionArgs): TActionObj["id"] {
  return devSpaceStore.startAction({
    actionName: "stop",
    workspaceID,
    actionFn: async (ctx) => {
      const result = await client.workspaces.stop(onStream, {
        id: workspaceID,
        actionID: ctx.id,
        streamID,
      })
      if (result.err) {
        return result
      }
      devSpaceStore.setStatus(workspaceID, result.val)

      return result
    },
  })
}

type TRemoveWorkspaceActionArgs = Readonly<{
  onStream?: TStreamEventListenerFn
  workspaceID: TWorkspaceID
  streamID: TStreamID
  force: boolean
}>
export function removeWorkspaceAction({
  workspaceID,
  onStream,
  streamID,
  force,
}: TRemoveWorkspaceActionArgs): TActionObj["id"] {
  return devSpaceStore.startAction({
    actionName: "remove",
    workspaceID,
    actionFn: async (ctx) => {
      const result = await client.workspaces.remove(force, onStream, {
        id: workspaceID,
        actionID: ctx.id,
        streamID,
      })
      if (result.err) {
        return result
      }
      devSpaceStore.removeWorkspace(workspaceID)

      return result
    },
  })
}
