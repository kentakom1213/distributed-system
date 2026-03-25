use crate::{BroadcastError, broadcast::Broadcast, config::Config, network::Network};

#[derive(Debug)]
pub struct Node {
    config: Config,
    network: Network,
    broadcast: Broadcast,
}

impl Node {
    pub fn new(config: Config) -> Self {
        let network = Network::new(config.network.clone());
        let broadcast = Broadcast::new(config.node_id);

        Self {
            config,
            network,
            broadcast,
        }
    }

    pub fn run(&mut self) -> Result<(), BroadcastError> {
        loop {}
    }
}
