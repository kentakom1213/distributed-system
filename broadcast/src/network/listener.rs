use std::net::SocketAddr;

use tokio::net::TcpListener;

use crate::{message::NodeId, network::connection::handle_connection};

pub struct Listener {
    listener: TcpListener,
}

impl Listener {
    pub async fn bind(addr: SocketAddr) -> std::io::Result<Self> {
        let listener = TcpListener::bind(addr).await?;

        Ok(Self { listener })
    }

    pub async fn run(self, node_id: NodeId) {
        loop {
            match self.listener.accept().await {
                Ok((stream, addr)) => {
                    tracing::info!("{}: connected from {}", node_id, addr);

                    tokio::spawn(async move {
                        if let Err(e) = handle_connection(node_id, stream).await {
                            tracing::error!("{}: connection error: {}", node_id, e);
                        }
                    });
                }
                Err(e) => {
                    tracing::error!("{node_id}: accept error: {e}");
                }
            }
        }
    }
}
