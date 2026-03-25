use std::net::SocketAddr;

use tokio::net::TcpListener;

use crate::{BroadcastError, network::connection};

pub async fn run_listener(addr: SocketAddr) -> Result<(), BroadcastError> {
    let listener = TcpListener::bind(addr).await?;

    loop {
        let (stream, peer_addr) = listener.accept().await?;

        tokio::spawn(async move {
            connection::handle_connection(stream).await;
        });
    }
}
