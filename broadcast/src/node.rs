use crate::{broadcast::Broadcast, config::Config, network::Network};

pub struct Node {
    config: Config,
    network: Network,
    broadcast: Broadcast,
}
