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
}
