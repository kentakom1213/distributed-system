use serde::Deserialize;

#[derive(Debug, Clone, Copy, Deserialize)]
pub struct NodeId(pub usize);

impl std::fmt::Display for NodeId {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.0)
    }
}

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
