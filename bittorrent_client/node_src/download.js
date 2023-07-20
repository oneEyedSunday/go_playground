'use strict';

const net = require('net');
const tracker = require('./tracker');
const message = require('./message');

function download(peer, torrent) {
    const socket = new net.Socket();

    socket.on('error', console.error);
    socket.connect(peer.port, peer.ip, () => {
        socket.write(message.buildHandshake(torrent));
    });
    onWholeMsg(socker, data => msgHandler(data, socket));
}

function msgHandler(msg, socket) {
    if (isHandshake(msg)) socket.write(message.buildInterested());
}

function isHandshake(msg) {
    return msg.length === msg.readUInt8(0) + 49 &&
         msg.toString('utf8', 1) === 'BitTorrent protocol';
}

function onWholeMsg(socket, callback) {
    // this whole func is suss
    let savedBuf = Buffer.alloc(0);
    let handshake = true;
    socket.on('data', recvBuf => {
        // msgLen calculates the length of a whole message
        // How do we come across tehse magic numbers?
        // This is because the handshake message doesn’t tell you its length as part of the message. 
        // The only way you can tell you’re receiving a handshake message is that it’s always the first message you’ll receive. That’s why I start with handshake set to true, and then the first time we receive a whole message I set it to false.
        const msgLen = () => handshake ? savedBuf.readUInt8(0) + 49 : savedBuf.readInt32BE(0) + 4;

        while (savedBuf.length >= 4 && savedBuf.length >= msgLen()) {
            callback(savedBuf.slice(0, msgLen()));
            savedBuf = savedBuf.slice(msgLen());
            handshake = false;
        }
    });
}

module.exports = (torrent) => {
    tracker.getPeers(torrent, peers => {
        peers.forEach(peer => download(peer, torrent));
    });
};
