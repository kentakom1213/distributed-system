use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Copy, Serialize, Deserialize)]
pub struct NodeId(pub usize);

impl std::fmt::Display for NodeId {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "node{}", self.0)
    }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct MessageId {
    pub node: NodeId,
    pub seq: u64,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(tag = "type", rename_all = "snake_case")]
pub enum Message {
    Ping,
    Pong,
    Broadcast {
        id: MessageId,
        sender: NodeId,
        payload: String,
    },
}

#[cfg(test)]
mod test_message {
    use crate::message::{Message, MessageId, NodeId};

    #[test]
    fn test_serialize() {
        let ping = Message::Ping;
        assert_eq!(serde_json::to_string(&ping).unwrap(), r#"{"type":"ping"}"#);

        let broadcast = Message::Broadcast {
            id: MessageId {
                node: NodeId(2),
                seq: 1,
            },
            sender: NodeId(1),
            payload: "こんにちは".to_string(),
        };
        assert_eq!(
            serde_json::to_string(&broadcast).unwrap(),
            r#"{"type":"broadcast","id":{"node":2,"seq":1},"sender":1,"payload":"こんにちは"}"#
        );
    }

    #[test]
    fn test_deserialize() {
        let pong_json = r#"{"type":"pong"}"#;
        let pong: Message = serde_json::from_str(pong_json).unwrap();
        assert!(matches!(pong, Message::Pong));

        let broadcast_json =
            r#"{"type":"broadcast","id":{"node":2,"seq":1},"sender":1,"payload":"こんにちは"}"#;
        let broadcast: Message = serde_json::from_str(broadcast_json).unwrap();
        if let Message::Broadcast {
            id,
            sender,
            payload,
        } = broadcast
        {
            assert_eq!(id.node.0, 2);
            assert_eq!(id.seq, 1);
            assert_eq!(sender.0, 1);
            assert_eq!(payload, "こんにちは");
        } else {
            panic!("Expected Broadcast message");
        }
    }
}
