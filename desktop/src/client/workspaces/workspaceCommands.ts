import { exists, Result, Return } from "../../lib"
import {
  TWorkspace,
  TWorkspaceID,
  TWorkspaceStartConfig,
  TWorkspaceStatusResult,
  TWorkspaceWithoutStatus,
} from "../../types"
import { Command, isOk, serializeRawOptions, toFlagArg, toMultipleFlagArg } from "../command"
import {
  DEVSPACE_COMMAND_DELETE,
  DEVSPACE_COMMAND_GET_WORKSPACE_CONFIG,
  DEVSPACE_COMMAND_GET_WORKSPACE_NAME,
  DEVSPACE_COMMAND_HELPER,
  DEVSPACE_COMMAND_LIST,
  DEVSPACE_COMMAND_STATUS,
  DEVSPACE_COMMAND_STOP,
  DEVSPACE_COMMAND_UP,
  DEVSPACE_FLAG_DEBUG,
  DEVSPACE_FLAG_DEVCONTAINER_PATH,
  DEVSPACE_FLAG_FORCE,
  DEVSPACE_FLAG_ID,
  DEVSPACE_FLAG_IDE,
  DEVSPACE_FLAG_JSON_LOG_OUTPUT,
  DEVSPACE_FLAG_JSON_OUTPUT,
  DEVSPACE_FLAG_PREBUILD_REPOSITORY,
  DEVSPACE_FLAG_PROVIDER,
  DEVSPACE_FLAG_PROVIDER_OPTION,
  DEVSPACE_FLAG_RECREATE,
  DEVSPACE_FLAG_RESET,
  DEVSPACE_FLAG_SOURCE,
  DEVSPACE_FLAG_TIMEOUT,
} from "../constants"

type TRawWorkspaces = readonly (Omit<TWorkspace, "status" | "id"> &
  Readonly<{ id: string | null }>)[]

export class WorkspaceCommands {
  static DEBUG = false
  static ADDITIONAL_FLAGS = ""

  private static newCommand(args: string[]): Command {
    return new Command([...args, ...(WorkspaceCommands.DEBUG ? [DEVSPACE_FLAG_DEBUG] : [])])
  }

  static async ListWorkspaces(): Promise<Result<TWorkspaceWithoutStatus[]>> {
    const result = await new Command([DEVSPACE_COMMAND_LIST, DEVSPACE_FLAG_JSON_OUTPUT]).run()
    if (result.err) {
      return result
    }

    const rawWorkspaces = JSON.parse(result.val.stdout) as TRawWorkspaces

    return Return.Value(
      rawWorkspaces.filter((workspace): workspace is TWorkspaceWithoutStatus =>
        exists(workspace.id)
      )
    )
  }

  static async FetchWorkspaceStatus(
    id: string
  ): Promise<Result<Pick<TWorkspace, "id" | "status">>> {
    const result = await new Command([DEVSPACE_COMMAND_STATUS, id, DEVSPACE_FLAG_JSON_OUTPUT]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return Return.Failed(`Failed to get status for workspace ${id}: ${result.val.stderr}`)
    }

    const { state } = JSON.parse(result.val.stdout) as TWorkspaceStatusResult

    return Return.Value({ id, status: state })
  }

  static async GetWorkspaceID(source: string) {
    const result = await new Command([
      DEVSPACE_COMMAND_HELPER,
      DEVSPACE_COMMAND_GET_WORKSPACE_NAME,
      source,
    ]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return Return.Failed(`Failed to get ID for workspace source ${source}: ${result.val.stderr}`)
    }

    return Return.Value(result.val.stdout)
  }

  static GetStatusLogs(id: string) {
    return new Command([DEVSPACE_COMMAND_STATUS, id, DEVSPACE_FLAG_JSON_LOG_OUTPUT])
  }

  static StartWorkspace(id: TWorkspaceID, config: TWorkspaceStartConfig) {
    const maybeSource = config.sourceConfig?.source
    const maybeIDFlag = exists(maybeSource) ? [toFlagArg(DEVSPACE_FLAG_ID, id)] : []

    const maybeSourceType = config.sourceConfig?.type
    const maybeSourceFlag =
      exists(maybeSourceType) && exists(maybeSource)
        ? [toFlagArg(DEVSPACE_FLAG_SOURCE, `${maybeSourceType}:${maybeSource}`)]
        : []
    const identifier = exists(maybeSource) && exists(maybeIDFlag) ? maybeSource : id

    const maybeIdeName = config.ideConfig?.name
    const maybeIDEFlag = exists(maybeIdeName) ? [toFlagArg(DEVSPACE_FLAG_IDE, maybeIdeName)] : []

    const maybeProviderID = config.providerConfig?.providerID
    const maybeProviderFlag = exists(maybeProviderID)
      ? [toFlagArg(DEVSPACE_FLAG_PROVIDER, maybeProviderID)]
      : []
    const maybeProviderOptions = config.providerConfig?.options
    const maybeProviderOptionsFlag = exists(maybeProviderOptions)
      ? serializeRawOptions(maybeProviderOptions, DEVSPACE_FLAG_PROVIDER_OPTION)
      : []

    const maybePrebuildRepositories = config.prebuildRepositories?.length
      ? [toFlagArg(DEVSPACE_FLAG_PREBUILD_REPOSITORY, config.prebuildRepositories.join(","))]
      : []

    const maybeDevcontainerPath = config.devcontainerPath
      ? [toFlagArg(DEVSPACE_FLAG_DEVCONTAINER_PATH, config.devcontainerPath)]
      : []

    const additionalFlags =
      WorkspaceCommands.ADDITIONAL_FLAGS.length !== 0
        ? toMultipleFlagArg(WorkspaceCommands.ADDITIONAL_FLAGS)
        : []

    return WorkspaceCommands.newCommand([
      DEVSPACE_COMMAND_UP,
      identifier,
      ...maybeIDFlag,
      ...maybeSourceFlag,
      ...maybeIDEFlag,
      ...maybeProviderFlag,
      ...maybePrebuildRepositories,
      ...maybeDevcontainerPath,
      ...additionalFlags,
      ...maybeProviderOptionsFlag,
      DEVSPACE_FLAG_JSON_LOG_OUTPUT,
    ])
  }

  static StopWorkspace(id: TWorkspaceID) {
    return WorkspaceCommands.newCommand([DEVSPACE_COMMAND_STOP, id, DEVSPACE_FLAG_JSON_LOG_OUTPUT])
  }

  static RebuildWorkspace(id: TWorkspaceID) {
    return WorkspaceCommands.newCommand([
      DEVSPACE_COMMAND_UP,
      id,
      DEVSPACE_FLAG_JSON_LOG_OUTPUT,
      DEVSPACE_FLAG_RECREATE,
    ])
  }

  static ResetWorkspace(id: TWorkspaceID) {
    return WorkspaceCommands.newCommand([
      DEVSPACE_COMMAND_UP,
      id,
      DEVSPACE_FLAG_JSON_LOG_OUTPUT,
      DEVSPACE_FLAG_RESET,
    ])
  }

  static RemoveWorkspace(id: TWorkspaceID, force?: boolean) {
    const args = [DEVSPACE_COMMAND_DELETE, id, DEVSPACE_FLAG_JSON_LOG_OUTPUT]
    if (force) {
      args.push(DEVSPACE_FLAG_FORCE)
    }

    return WorkspaceCommands.newCommand(args)
  }

  static GetDevcontainerConfig(rawSource: string) {
    return new Command([
      DEVSPACE_COMMAND_HELPER,
      DEVSPACE_COMMAND_GET_WORKSPACE_CONFIG,
      rawSource,
      DEVSPACE_FLAG_TIMEOUT,
      "10s",
    ])
  }
}
