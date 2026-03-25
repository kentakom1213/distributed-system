use std::net::SocketAddr;

use crate::message::NodeId;

pub struct Config {
    node_id: NodeId,
    listen_addr: SocketAddr,
    peers: Vec<SocketAddr>,
}
