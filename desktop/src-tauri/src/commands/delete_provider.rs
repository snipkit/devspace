use super::{
    config::{CommandConfig, DevspaceCommandConfig, DevspaceCommandError},
    constants::{DEVSPACE_BINARY_NAME, DEVSPACE_COMMAND_DELETE, DEVSPACE_COMMAND_PROVIDER},
};

pub struct DeleteProviderCommand {
    provider_id: String,
}
impl DeleteProviderCommand {
    pub fn new(provider_id: String) -> Self {
        DeleteProviderCommand { provider_id }
    }
}
impl DevspaceCommandConfig<()> for DeleteProviderCommand {
    fn config(&self) -> CommandConfig {
        CommandConfig {
            binary_name: DEVSPACE_BINARY_NAME,
            args: vec![
                DEVSPACE_COMMAND_PROVIDER,
                DEVSPACE_COMMAND_DELETE,
                &self.provider_id,
            ],
        }
    }

    fn exec(self) -> Result<(), DevspaceCommandError> {
        let cmd = self.new_command()?;

        cmd.status()
            .map_err(DevspaceCommandError::Failed)?
            .success()
            .then_some(())
            .ok_or_else(|| DevspaceCommandError::Exit)
    }
}
