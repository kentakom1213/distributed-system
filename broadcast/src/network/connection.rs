use std::net::SocketAddr;

use tokio::{io::BufReader, net::TcpStream};

use crate::{BroadcastError, handler::MessageHandler, message::receive_message};

pub async fn handle_connection(
    handler: MessageHandler,
    stream: TcpStream,
    peer_addr: SocketAddr,
) -> Result<(), BroadcastError> {
    let mut reader = BufReader::new(stream);
    let node_id = handler.node_id();

    loop {
        match receive_message(&mut reader).await {
            Ok(msg) => handler.handle(msg, peer_addr).await?,
            Err(BroadcastError::Closed) => {
                tracing::info!("[{node_id}] Connection Closed");
                break;
            }
            Err(e) => tracing::error!("[{node_id}] Error: {e}"),
        }
    }

    Ok(())
}
