use tokio::{io::BufReader, net::TcpStream};

use crate::{
    BroadcastError,
    message::{NodeId, receive_message},
};

pub async fn handle_connection(node_id: NodeId, stream: TcpStream) -> Result<(), BroadcastError> {
    let mut reader = BufReader::new(stream);

    loop {
        match receive_message(&mut reader).await {
            Ok(msg) => tracing::info!("[{node_id}] received {msg:?}"),
            Err(BroadcastError::Closed) => {
                tracing::info!("[{node_id}] Connection Closed");
                break;
            }
            Err(e) => tracing::error!("[{node_id}] Error: {e}"),
        }
    }

    Ok(())
}
