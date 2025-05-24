use super::{
    config::{CommandConfig, DevspaceCommandConfig, DevspaceCommandError},
    constants::{
        DEVSPACE_BINARY_NAME, DEVSPACE_COMMAND_DELETE, DEVSPACE_COMMAND_PRO, FLAG_IGNORE_NOT_FOUND,
    },
};

pub struct DeleteProInstanceCommand {
    pro_id: String,
}
impl DeleteProInstanceCommand {
    pub fn new(pro_id: String) -> Self {
        DeleteProInstanceCommand { pro_id }
    }
}
impl DevspaceCommandConfig<()> for DeleteProInstanceCommand {
    fn config(&self) -> CommandConfig {
        CommandConfig {
            binary_name: DEVSPACE_BINARY_NAME,
            args: vec![
                DEVSPACE_COMMAND_PRO,
                DEVSPACE_COMMAND_DELETE,
                &self.pro_id,
                FLAG_IGNORE_NOT_FOUND,
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
