use tokio::net::TcpStream;

use crate::{BroadcastError, message::NodeId};

pub async fn handle_connection(
    node_id: NodeId,
    mut stream: TcpStream,
) -> Result<(), BroadcastError> {
    tracing::info!("{node_id}: Receive connection");

    Ok(())
}
