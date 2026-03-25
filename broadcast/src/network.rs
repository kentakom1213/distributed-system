use std::net::SocketAddr;

#[derive(Debug)]
pub struct Network {
    listen_addr: SocketAddr,
    peers: Vec<SocketAddr>,
}
