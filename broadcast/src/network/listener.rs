use std::net::SocketAddr;

use tokio::net::TcpListener;

use crate::{handler::MessageHandler, message::NodeId, network::connection::handle_connection};

pub struct Listener {
    listener: TcpListener,
}

impl Listener {
    pub async fn bind(addr: SocketAddr) -> std::io::Result<Self> {
        let listener = TcpListener::bind(addr).await?;

        Ok(Self { listener })
    }

    pub async fn run(self, node_id: NodeId, handler: MessageHandler) {
        loop {
            match self.listener.accept().await {
                Ok((stream, addr)) => {
                    tracing::info!("[{node_id}] Connected from {addr}");
                    let handler = handler.clone();

                    tokio::spawn(async move {
                        if let Err(e) = handle_connection(handler, stream, addr).await {
                            tracing::error!("[{node_id}] Connection error: {e}");
                        }
                    });
                }
                Err(e) => {
                    tracing::error!("[{node_id}] Accept error: {e}");
                }
            }
        }
    }
}
