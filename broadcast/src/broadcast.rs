use std::collections::HashSet;

use crate::message::{MessageId, NodeId};

pub struct Broadcast {
    node_id: NodeId,
    seq: u64,
    seen_messages: HashSet<MessageId>,
}
