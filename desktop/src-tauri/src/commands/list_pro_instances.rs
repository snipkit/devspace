use serde::{Deserialize, Serialize};

use super::{
    config::{CommandConfig, DevspaceCommandConfig, DevspaceCommandError},
    constants::{
        DEVSPACE_BINARY_NAME, DEVSPACE_COMMAND_LIST, DEVSPACE_COMMAND_PRO, FLAG_OUTPUT_JSON,
    },
};

#[derive(Serialize, Deserialize, Debug, Eq, PartialEq)]
#[serde(rename_all(serialize = "camelCase", deserialize = "camelCase"))]
pub struct ProInstance {
    id: Option<String>,
    url: Option<String>,
    creation_timestamp: Option<chrono::DateTime<chrono::Utc>>,
}
impl ProInstance {
    pub fn id(&self) -> Option<&String> {
        self.id.as_ref()
    }
}

pub struct ListProInstancesCommand {}
impl ListProInstancesCommand {
    pub fn new() -> Self {
        ListProInstancesCommand {}
    }

    fn deserialize(&self, str: &str) -> Result<Vec<ProInstance>, DevspaceCommandError> {
        serde_json::from_str(str).map_err(DevspaceCommandError::Parse)
    }
}
impl DevspaceCommandConfig<Vec<ProInstance>> for ListProInstancesCommand {
    fn config(&self) -> CommandConfig {
        CommandConfig {
            binary_name: DEVSPACE_BINARY_NAME,
            args: vec![
                DEVSPACE_COMMAND_PRO,
                DEVSPACE_COMMAND_LIST,
                FLAG_OUTPUT_JSON,
            ],
        }
    }

    fn exec(self) -> Result<Vec<ProInstance>, DevspaceCommandError> {
        let output = self
            .new_command()?
            .output()
            .map_err(|_| DevspaceCommandError::Output)?;

        self.deserialize(&output.stdout)
    }
}
