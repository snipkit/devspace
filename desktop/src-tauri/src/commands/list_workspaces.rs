use super::{
    config::{CommandConfig, DevspaceCommandConfig, DevspaceCommandError},
    constants::{DEVSPACE_BINARY_NAME, DEVSPACE_COMMAND_LIST, FLAG_OUTPUT_JSON},
};
use crate::workspaces::WorkspacesState;

pub struct ListWorkspacesCommand {}
impl ListWorkspacesCommand {
    pub fn new() -> Self {
        ListWorkspacesCommand {}
    }

    fn deserialize(&self, str: &str) -> Result<WorkspacesState, DevspaceCommandError> {
        serde_json::from_str(str).map_err(DevspaceCommandError::Parse)
    }
}
impl DevspaceCommandConfig<WorkspacesState> for ListWorkspacesCommand {
    fn config(&self) -> CommandConfig {
        CommandConfig {
            binary_name: DEVSPACE_BINARY_NAME,
            args: vec![DEVSPACE_COMMAND_LIST, FLAG_OUTPUT_JSON],
        }
    }

    fn exec(self) -> Result<WorkspacesState, DevspaceCommandError> {
        let output = self
            .new_command()?
            .output()
            .map_err(|_| DevspaceCommandError::Output)?;

        self.deserialize(&output.stdout)
    }
}
