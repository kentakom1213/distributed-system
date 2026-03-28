use std::sync::Arc;

use tokio::sync::Mutex;

use crate::{
    BroadcastError,
    broadcast::Broadcast,
    config::Config,
    handler::MessageHandler,
    network::{Network, listener::Listener},
};

#[derive(Debug)]
pub struct Node {
    pub config: Config,
    pub network: Network,
    pub broadcast: Arc<Mutex<Broadcast>>,
}

impl Node {
    pub fn new(config: Config) -> Self {
        let network = Network::new(config.network.clone());
        let broadcast = Arc::new(Mutex::new(Broadcast::new(config.node_id)));

        Self {
            config,
            network,
            broadcast,
        }
    }

    pub async fn run(&mut self) -> Result<(), BroadcastError> {
        let node_id = self.config.node_id;

        let listener = Listener::bind(self.network.listen_addr()).await?;
        let handler = MessageHandler::new(
            node_id,
            self.network.peers().to_vec(),
            Arc::clone(&self.broadcast),
        );

        tokio::spawn(async move {
            listener.run(node_id, handler).await;
        });

        tracing::info!("[{node_id}] Listening on {}", self.network.listen_addr());

        tokio::signal::ctrl_c().await?;

        Ok(())
    }
}
