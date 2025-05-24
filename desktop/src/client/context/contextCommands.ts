import { Result, ResultError, Return, getErrorFromChildProcess } from "../../lib"
import { TContextOptionName, TContextOptions } from "../../types"
import { Command, isOk, serializeRawOptions } from "../command"
import {
  DEVSPACE_COMMAND_CONTEXT,
  DEVSPACE_COMMAND_OPTIONS,
  DEVSPACE_COMMAND_SET_OPTIONS,
  DEVSPACE_FLAG_DEBUG,
  DEVSPACE_FLAG_JSON_LOG_OUTPUT,
  DEVSPACE_FLAG_JSON_OUTPUT,
} from "../constants"

export class ContextCommands {
  static DEBUG = false

  private static newCommand(args: string[]): Command {
    return new Command([...args, ...(ContextCommands.DEBUG ? [DEVSPACE_FLAG_DEBUG] : [])])
  }

  static async SetOptions(
    rawOptions: Partial<Record<TContextOptionName, string>>
  ): Promise<ResultError> {
    const optionsFlag = serializeRawOptions(rawOptions)
    const result = await ContextCommands.newCommand([
      DEVSPACE_COMMAND_CONTEXT,
      DEVSPACE_COMMAND_SET_OPTIONS,
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

  static async ListOptions(): Promise<Result<TContextOptions>> {
    const result = await ContextCommands.newCommand([
      DEVSPACE_COMMAND_CONTEXT,
      DEVSPACE_COMMAND_OPTIONS,
      DEVSPACE_FLAG_JSON_OUTPUT,
    ]).run()
    if (result.err) {
      return result
    }

    if (!isOk(result.val)) {
      return getErrorFromChildProcess(result.val)
    }

    const options = JSON.parse(result.val.stdout) as TContextOptions

    return Return.Value(options)
  }
}
