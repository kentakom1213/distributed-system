use std::net::SocketAddr;

pub struct Network {
    listen_addr: SocketAddr,
    peers: Vec<SocketAddr>,
}
