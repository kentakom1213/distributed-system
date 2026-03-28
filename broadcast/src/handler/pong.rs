use crate::message::NodeId;

pub fn handle(node_id: NodeId, from: NodeId) {
    tracing::info!("[{}] Received Pong from {}", node_id, from);
}
