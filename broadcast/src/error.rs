#[derive(Debug, thiserror::Error)]
pub enum BroadcastError {
    #[error("IO error: {0}")]
    Io(#[from] std::io::Error),

    #[error("serialization error: {0}")]
    Serialization(#[from] serde_json::Error),

    #[error("invalid message: {0}")]
    InvalidMessage(String),

    #[error("network error: {0}")]
    Network(String),

    #[error("configuration error: {0}")]
    Config(String),

    #[error("config parse error: {0}")]
    ConfigParse(#[from] toml::de::Error),
}
