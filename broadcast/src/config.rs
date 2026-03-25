use std::net::SocketAddr;

use serde::Deserialize;

use crate::message::NodeId;

#[derive(Debug, Deserialize)]
pub struct Config {
    node_id: NodeId,
    network: NetworkConf,
}

#[derive(Debug, Deserialize)]
pub struct NetworkConf {
    listen_addr: SocketAddr,
    peers: Vec<SocketAddr>,
}
