use std::net::SocketAddr;

use tokio::{io::AsyncWriteExt, net::TcpStream};

use crate::BroadcastError;

pub async fn send(addr: SocketAddr, buf: &[u8]) -> Result<(), BroadcastError> {
    let mut stream = TcpStream::connect(addr).await?;

    stream.write_all(buf).await?;

    Ok(())
}

pub async fn send_all(peers: &[SocketAddr], buf: &[u8]) {
    for peer in peers {
        if let Err(e) = send(*peer, buf).await {
            tracing::error!("send error {peer}: {e}");
        }
    }
}
