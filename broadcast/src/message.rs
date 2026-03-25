pub struct NodeId(usize);

pub struct MessageId {
    pub node: NodeId,
    pub seq: u64,
}

pub struct Message {
    pub id: MessageId,
    pub sender: NodeId,
    pub payload: String,
}
