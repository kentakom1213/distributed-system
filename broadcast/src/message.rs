use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct NodeId(usize);

#[derive(Debug)]
pub struct MessageId {
    pub node: NodeId,
    pub seq: u64,
}

#[derive(Debug)]
pub struct Message {
    pub id: MessageId,
    pub sender: NodeId,
    pub payload: String,
}
