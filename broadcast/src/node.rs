use crate::{broadcast::Broadcast, config::Config, network::Network};

#[derive(Debug)]
pub struct Node {
    config: Config,
    network: Network,
    broadcast: Broadcast,
}
