use std::{net::SocketAddr, sync::Arc};

use tokio::sync::Mutex;

use crate::{
    BroadcastError,
    broadcast::Broadcast,
    message::{Message, NodeId},
};

pub mod broadcast;
pub mod ping;
pub mod pong;

#[derive(Clone, Debug)]
pub struct MessageHandler {
    node_id: NodeId,
    peers: Vec<SocketAddr>,
    broadcast: Arc<Mutex<Broadcast>>,
}

impl MessageHandler {
    pub fn new(node_id: NodeId, peers: Vec<SocketAddr>, broadcast: Arc<Mutex<Broadcast>>) -> Self {
        Self {
            node_id,
            peers,
            broadcast,
        }
    }

    pub fn node_id(&self) -> NodeId {
        self.node_id
    }

    pub async fn handle(&self, msg: Message, peer_addr: SocketAddr) -> Result<(), BroadcastError> {
        match msg {
            Message::Ping { from } => ping::handle(self.node_id, from, peer_addr).await,
            Message::Pong { from } => {
                pong::handle(self.node_id, from);
                Ok(())
            }
            Message::Broadcast {
                id,
                sender,
                payload,
            } => {
                let mut state = self.broadcast.lock().await;
                broadcast::handle(&mut state, self.node_id, &self.peers, id, sender, &payload).await
            }
        }
    }
}
