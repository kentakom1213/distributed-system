use std::collections::HashSet;

use crate::message::{MessageId, NodeId};

#[derive(Debug)]
pub struct Broadcast {
    node_id: NodeId,
    seq: u64,
    seen_messages: HashSet<MessageId>,
}

impl Broadcast {
    pub fn new(node_id: NodeId) -> Self {
        Self {
            node_id,
            seq: 0,
            seen_messages: HashSet::new(),
        }
    }

    pub fn node_id(&self) -> NodeId {
        self.node_id
    }

    pub fn next_message_id(&mut self) -> MessageId {
        self.seq += 1;
        MessageId {
            node: self.node_id,
            seq: self.seq,
        }
    }

    pub fn record_received(&mut self, id: MessageId) -> bool {
        self.seen_messages.insert(id)
    }
}
