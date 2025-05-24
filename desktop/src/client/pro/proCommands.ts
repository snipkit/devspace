import { Result, ResultError, Return, getErrorFromChildProcess } from "@/lib"
import { TImportWorkspaceConfig, TListProInstancesConfig, TProID, TProInstance } from "@/types"
import { Command, isOk, serializeRawOptions, toFlagArg } from "../command"
import {
  DEVSPACE_COMMAND_DELETE,
  DEVSPACE_COMMAND_IMPORT_WORKSPACE,
  DEVSPACE_COMMAND_LIST,
  DEVSPACE_COMMAND_LOGIN,
  DEVSPACE_COMMAND_PRO,
  DEVSPACE_FLAG_ACCESS_KEY,
  DEVSPACE_FLAG_DEBUG,
  DEVSPACE_FLAG_JSON_LOG_OUTPUT,
  DEVSPACE_FLAG_JSON_OUTPUT,
  DEVSPACE_FLAG_LOGIN,
  DEVSPACE_FLAG_PROVIDER,
  DEVSPACE_FLAG_USE,
  DEVSPACE_FLAG_WORKSPACE_ID,
  DEVSPACE_FLAG_WORKSPACE_PROJECT,
  DEVSPACE_FLAG_WORKSPACE_UID,
} from "../constants"
import { TStreamEventListenerFn } from "../types"

export class ProCommands {
  static DEBUG = false

  private static newCommand(args: string[]): Command {
    return new Command([...args, ...(ProCommands.DEBUG ? [DEVSPACE_FLAG_DEBUG] : [])])
  }

  static async Login(
    host: string,
    providerName?: string,
    accessKey?: string,
    listener?: TStreamEventListenerFn
  ): Promise<ResultError> {
    const maybeProviderNameFlag = providerName
      ? [toFlagArg(DEVSPACE_FLAG_PROVIDER, providerName)]
      : []
    const maybeAccessKeyFlag = accessKey ? [toFlagArg(DEVSPACE_FLAG_ACCESS_KEY, accessKey)] : []
    const useFlag = toFlagArg(DEVSPACE_FLAG_USE, "false")

    const cmd = ProCommands.newCommand([
      DEVSPACE_COMMAND_PRO,
      DEVSPACE_COMMAND_LOGIN,
      host,
      useFlag,
      DEVSPACE_FLAG_JSON_LOG_OUTPUT,
      ...maybeProviderNameFlag,
      ...maybeAccessKeyFlag,
    ])
    if (listener) {
      return cmd.stream(listener)
    } else {
      const result = await cmd.run()
      if (result.err) {
        return result
      }

      if (!isOk(result.val)) {
        return getErrorFromChildProcess(result.val)
      }

      return Return.Ok()
    }
  }

  static async ListProInstances(
    config?: TListProInstancesConfig
  ): Promise<Result<readonly TProInstance[]>> {
    const maybeLoginFlag = config?.authenticate ? [DEVSPACE_FLAG_LOGIN] : []
    const result = await ProCommands.newCommand([
      DEVSPACE_COMMAND_PRO,
      DEVSPACE_COMMAND_LIST,
      DEVSPACE_FLAG_JSON_OUTPUT,
      ...maybeLoginFlag,
    ]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    const instances = JSON.parse(result.val.stdout) as readonly TProInstance[]

    return Return.Value(instances)
  }

  static async RemoveProInstance(id: TProID) {
    const result = await ProCommands.newCommand([
      DEVSPACE_COMMAND_PRO,
      DEVSPACE_COMMAND_DELETE,
      id,
      DEVSPACE_FLAG_JSON_LOG_OUTPUT,
    ]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    return Return.Ok()
  }

  static async ImportWorkspace(config: TImportWorkspaceConfig): Promise<ResultError> {
    const optionsFlag = config.options ? serializeRawOptions(config.options) : []
    const result = await new Command([
      DEVSPACE_COMMAND_PRO,
      DEVSPACE_COMMAND_IMPORT_WORKSPACE,
      config.devSpaceProHost,
      DEVSPACE_FLAG_WORKSPACE_ID,
      config.workspaceID,
      DEVSPACE_FLAG_WORKSPACE_UID,
      config.workspaceUID,
      DEVSPACE_FLAG_WORKSPACE_PROJECT,
      config.project,
      ...optionsFlag,
      DEVSPACE_FLAG_JSON_LOG_OUTPUT,
    ]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    return Return.Ok()
  }
}
