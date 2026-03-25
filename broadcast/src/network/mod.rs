use std::net::SocketAddr;

use crate::NetworkConf;

pub mod connection;
pub mod listener;

#[derive(Debug)]
pub struct Network {
    conf: NetworkConf,
}

impl Network {
    pub fn new(conf: NetworkConf) -> Self {
        Self { conf }
    }

    pub fn listen_addr(&self) -> SocketAddr {
        self.conf.listen_addr
    }

    pub fn peers(&self) -> &[SocketAddr] {
        &self.conf.peers
    }
}
