import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { ReactQueryDevtools } from "@tanstack/react-query-devtools"
import dayjs from "dayjs"
import relativeTime from "dayjs/plugin/relativeTime"
import { StrictMode } from "react"
import ReactDOM from "react-dom/client"
import { RouterProvider } from "react-router"
import "xterm/css/xterm.css"
import { DevSpaceProvider, SettingsProvider } from "./contexts"
import { router } from "./routes"
import { ThemeProvider } from "./Theme"

dayjs.extend(relativeTime)

const queryClient = new QueryClient({
  logger: {
    log(...args) {
      console.log(args)
    },
    warn(...args) {
      console.warn(args)
    },
    error(...args) {
      const maybeError = args[0]
      if (maybeError instanceof Error) {
        console.error(maybeError.name, maybeError.message, maybeError.cause, maybeError)

        return
      }

      console.error(args)
    },
  },
})

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(<Root />)

function Root() {
  return (
    <StrictMode>
      <SettingsProvider>
        <ThemeProvider>
          <QueryClientProvider client={queryClient}>
            <DevSpaceProvider>
              <RouterProvider router={router} />
            </DevSpaceProvider>
            {/* Will be disabled in production automatically */}
            <ReactQueryDevtools
              position="bottom-right"
              toggleButtonProps={{ style: { margin: "0.5em 0.5em 2rem" } }}
            />
          </QueryClientProvider>
        </ThemeProvider>
      </SettingsProvider>
    </StrictMode>
  )
}
