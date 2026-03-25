use crate::{
    BroadcastError,
    broadcast::Broadcast,
    config::Config,
    network::{Network, listener::Listener},
};

#[derive(Debug)]
pub struct Node {
    pub config: Config,
    pub network: Network,
    pub broadcast: Broadcast,
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

    pub async fn run(&mut self) -> Result<(), BroadcastError> {
        let node_id = self.config.node_id;

        let listener = Listener::bind(self.network.listen_addr()).await?;

        tokio::spawn(async move {
            listener.run(node_id).await;
        });

        tracing::info!(
            "node={node_id}: Listening on {}",
            self.network.listen_addr()
        );

        tokio::signal::ctrl_c().await?;

        Ok(())
    }
}
