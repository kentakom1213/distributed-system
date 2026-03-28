use serde::{Deserialize, Serialize};
use tokio::io::{AsyncBufReadExt, AsyncWriteExt};

use crate::BroadcastError;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash, Serialize, Deserialize)]
pub struct NodeId(pub usize);

impl std::fmt::Display for NodeId {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "node{}", self.0)
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash, Serialize, Deserialize)]
pub struct MessageId {
    pub node: NodeId,
    pub seq: u64,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(tag = "type", rename_all = "snake_case")]
pub enum Message {
    Ping {
        from: NodeId,
    },
    Pong {
        from: NodeId,
    },
    Broadcast {
        id: MessageId,
        sender: NodeId,
        payload: String,
    },
}

/// メッセージを改行区切りの JSON 形式で送信する
pub async fn send_message<W: AsyncWriteExt + Unpin>(
    mut w: W,
    msg: &Message,
) -> Result<(), BroadcastError> {
    let json = serde_json::to_string(msg)?;

    w.write_all(json.as_bytes()).await?;
    w.write_all(b"\n").await?;
    w.flush().await?;

    Ok(())
}

/// 改行区切りの JSON メッセージを受信する
pub async fn receive_message<R: AsyncBufReadExt + Unpin>(
    r: &mut R,
) -> Result<Message, BroadcastError> {
    let mut line = String::new();

    let n = r.read_line(&mut line).await?;
    if n == 0 {
        return Err(BroadcastError::Closed);
    }

    let msg: Message = serde_json::from_str(&line.trim())?;
    Ok(msg)
}

#[cfg(test)]
mod test_message {
    use crate::message::{Message, MessageId, NodeId};

    #[test]
    fn test_serialize() {
        let ping = Message::Ping { from: NodeId(1) };
        assert_eq!(
            serde_json::to_string(&ping).unwrap(),
            r#"{"type":"ping","from":1}"#
        );

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
        let pong_json = r#"{"type":"pong","from":20}"#;
        let pong: Message = serde_json::from_str(pong_json).unwrap();
        assert!(matches!(pong, Message::Pong { from: NodeId(20) }));

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
