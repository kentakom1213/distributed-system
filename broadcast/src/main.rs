use std::{fs, path::PathBuf};

use broadcast::{BroadcastError, Config, Node};
use clap::Parser;

/// Initialize broadcast node
#[derive(Debug, clap::Parser)]
struct Args {
    /// must be toml
    conf_path: PathBuf,
}

fn main() {
    if let Err(e) = run() {
        eprintln!("fatal error: {e}");
        std::process::exit(1);
    }
}

fn run() -> Result<(), BroadcastError> {
    let args = Args::parse();

    // ノードの設定を読み込む
    let conf_file = fs::read_to_string(args.conf_path)?;
    let conf: Config = toml::from_str(&conf_file)?;

    // ノードを起動
    let mut node = Node::new(conf);
    node.run()?;

    Ok(())
}
