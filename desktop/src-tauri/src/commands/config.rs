use std::collections::HashMap;

use tauri::api::process::Command;
use thiserror::Error;

use crate::commands::constants::DEVSPACE_BINARY_NAME;

use super::constants::DEVSPACE_UI_ENV_VAR;

pub struct CommandConfig<'a> {
    pub(crate) binary_name: &'static str,
    pub(crate) args: Vec<&'a str>,
}

impl<'a> CommandConfig<'_> {
    pub fn binary_name(&self) -> &'static str {
        self.binary_name
    }

    pub fn args(&self) -> &Vec<&str> {
        &self.args
    }
}

#[derive(Error, Debug)]
pub enum DevspaceCommandError {
    #[error("unable to parse command response")]
    Parse(#[from] serde_json::Error),
    #[error("unable to find sidecar binary")]
    Sidecar,
    #[error("unable to collect output from command")]
    Output,
    #[error("command failed")]
    Failed(#[from] tauri::api::Error),
    #[error("command exited with non-zero code")]
    Exit,
}
impl serde::Serialize for DevspaceCommandError {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        serializer.serialize_str(self.to_string().as_ref())
    }
}
pub trait DevspaceCommandConfig<T> {
    fn config(&self) -> CommandConfig {
        CommandConfig {
            binary_name: DEVSPACE_BINARY_NAME,
            args: vec![],
        }
    }
    fn exec(self) -> Result<T, DevspaceCommandError>;

    fn new_command(&self) -> Result<Command, DevspaceCommandError> {
        let config = self.config();
        let env_vars: HashMap<String, String> =
            HashMap::from([(DEVSPACE_UI_ENV_VAR.into(), "true".into())]);

        let cmd = Command::new_sidecar(config.binary_name())
            .map_err(|_| DevspaceCommandError::Sidecar)?
            .envs(env_vars)
            .args(config.args());

        Ok(cmd)
    }
}
