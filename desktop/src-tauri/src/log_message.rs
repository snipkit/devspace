use crate::{commands::DevspaceCommandError, AppState, UiMessage};
use log::{error, info, warn};

#[tauri::command]
pub fn log_message(message: String) {
    info!("logging message: {}", message);
}
