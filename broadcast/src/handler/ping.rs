use std::net::SocketAddr;

use crate::{
    BroadcastError,
    message::{Message, NodeId},
    network::sender::Sender,
};

pub async fn handle(
    node_id: NodeId,
    from: NodeId,
    peer_addr: SocketAddr,
) -> Result<(), BroadcastError> {
    tracing::info!("[{}] Received Ping from {}", node_id, from);

    let pong = Message::Pong { from: node_id };
    let mut sender = Sender::new(peer_addr).await?;
    sender.send(&pong).await?;

    tracing::info!("[{}] Send Pong to {}", node_id, from);

    Ok(())
}
