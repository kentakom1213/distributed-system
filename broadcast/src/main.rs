use std::{fs, path::PathBuf};

use broadcast::{BroadcastError, Config, Node, TimeNodeFormatter};
use clap::Parser;

/// Initialize broadcast node
#[derive(Debug, clap::Parser)]
struct Args {
    /// must be toml
    conf_path: PathBuf,
}

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt()
        .event_format(TimeNodeFormatter)
        .init();

    if let Err(e) = run().await {
        eprintln!("fatal error: {e}");
        std::process::exit(1);
    }
}

async fn run() -> Result<(), BroadcastError> {
    let args = Args::parse();

    // ノードの設定を読み込む
    let conf_file = fs::read_to_string(args.conf_path)?;
    let conf: Config = toml::from_str(&conf_file)?;

    tracing::info!("{}: Initializing", conf.node_id);

    // ノードを起動
    let mut node = Node::new(conf);

    node.run().await?;

    Ok(())
}
