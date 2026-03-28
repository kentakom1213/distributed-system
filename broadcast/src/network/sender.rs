use std::net::SocketAddr;

use tokio::net::TcpStream;

use crate::{
    BroadcastError,
    message::{Message, send_message},
};

pub struct Sender {
    pub stream: TcpStream,
}

impl Sender {
    pub async fn new(addr: SocketAddr) -> Result<Self, BroadcastError> {
        let stream = TcpStream::connect(addr).await?;

        Ok(Self { stream })
    }

    pub async fn send(&mut self, msg: &Message) -> Result<(), BroadcastError> {
        send_message(&mut self.stream, msg).await
    }
}
