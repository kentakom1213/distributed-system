use tokio::{
    io::{BufReader, BufWriter},
    net::TcpStream,
};

use crate::{
    BroadcastError,
    message::{Message, NodeId, receive_message, send_message},
};

pub async fn handle_connection(node_id: NodeId, stream: TcpStream) -> Result<(), BroadcastError> {
    let (reader, writer) = tokio::io::split(stream);
    let mut reader = BufReader::new(reader);
    let mut writer = BufWriter::new(writer);

    loop {
        match receive_message(&mut reader).await {
            Ok(msg) => match msg {
                Message::Ping { from } => {
                    tracing::info!("[{node_id}] Received Ping from {from}");

                    // Pong を返す
                    let pong = Message::Pong { from: node_id };
                    send_message(&mut writer, &pong).await?;

                    tracing::info!("[{node_id}] Send Pong to {from}");
                }
                Message::Pong { from } => {
                    tracing::info!("[{node_id}] Received Pong from {from}");
                }
                Message::Broadcast {
                    sender, payload, ..
                } => {
                    tracing::info!("[{node_id}] Received Broadcast from {sender}: {payload}");
                }
            },
            Err(BroadcastError::Closed) => {
                tracing::info!("[{node_id}] Connection Closed");
                break;
            }
            Err(e) => tracing::error!("[{node_id}] Error: {e}"),
        }
    }

    Ok(())
}
