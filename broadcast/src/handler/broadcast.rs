use std::net::SocketAddr;

use crate::{
    BroadcastError,
    broadcast::Broadcast,
    message::{Message, MessageId, NodeId},
    network::sender::Sender,
};

pub async fn handle(
    state: &mut Broadcast,
    node_id: NodeId,
    peers: &[SocketAddr],
    id: MessageId,
    sender: NodeId,
    payload: &str,
) -> Result<(), BroadcastError> {
    let is_new = state.record_received(id);

    tracing::info!(
        "[{}] Received Broadcast from {}: {}",
        node_id,
        sender,
        payload
    );

    if !is_new {
        tracing::debug!("[{}] Ignore duplicated broadcast {:?}", node_id, id);
        return Ok(());
    }

    let msg = Message::Broadcast {
        id,
        sender,
        payload: payload.to_string(),
    };

    for &peer in peers {
        let mut outbound = Sender::new(peer).await?;
        outbound.send(&msg).await?;
    }

    Ok(())
}
