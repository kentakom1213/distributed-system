use std::net::SocketAddr;

use serde::Deserialize;

use crate::message::NodeId;

#[derive(Debug, Deserialize)]
pub struct Config {
    pub node_id: NodeId,
    pub network: NetworkConf,
}

#[derive(Debug, Clone, Deserialize)]
pub struct NetworkConf {
    pub listen_addr: SocketAddr,
    pub peers: Vec<SocketAddr>,
}
