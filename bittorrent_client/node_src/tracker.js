'use strict';

const dgram = require('dgram');
const Buffer = require('buffer').Buffer;
const urlParse = require('url').parse;
const crypto = require('crypto');
const parser = require('./parser');
const util = require('./util');

function getPeers(torrent, callback) {
    const socket = dgram.createSocket('udp4');
    const url = torrent.announce.toString('utf8');

    udpSend(socket, buildConnReq(), url);

    socket.on('message', response => {
        if (respType(response) === 'connect') {
        // 2. receive and parse connect response
        const connResp = parseConnResp(response);
        // 3. send announce request
        const announceReq = buildAnnounceReq(connResp.connectionId, torrent);
        udpSend(socket, announceReq, url);
        } else if (respType(response) === 'announce') {
        // 4. parse announce response
        const announceResp = parseAnnounceResp(response);
        // 5. pass peers to callback
        callback(announceResp.peers);
        }
    })
}

/**
 * 
 * @param {dgram.Socket} socket 
 * @param {import('buffer').Buffer} message 
 * @param {string}} rawUrl 
 * @param {() => void} callback 
 */
function udpSend(socket, message, rawUrl, callback = () => {}) {
    const url = urlParse(rawUrl);
    socket.send(message, 0, message.length, url.port, url.host, callback);
}

function respType(resp) {
    const action = resp.readUInt32BE(0);
    if (action === 0) return 'connect';
    if (action === 1) return 'announce';
}

function buildConnReq() {
    const buf = Buffer.alloc(16);

    buf.writeUInt32BE(0x417, 0);
    buf.writeUInt32BE(0x27101980, 4);
    buf.writeUInt32BE(0, 8);
    crypto.randomBytes(4).copy(buf, 12);

    return buf;
}

/**
 * 
 * @param {Buffer} resp 
 */
function parseConnResp(resp) {
    return {
        action: resp.readUInt32BE(0),
        transactionId: resp.readUInt32BE(4),
        connectionId: resp.slice(8)
    }
}

/**
 * 
 * @param {Buffer} connId 
 * @param {*} torrent 
 * @param {number=} port 
 */
function buildAnnounceReq(connId, torrent, port = 6881) {
    const buf = Buffer.allocUnsafe(98);

    // connection id
    connId.copy(buf, 0);
    // action
    buf.writeUInt32BE(1, 8);
    // transaction id
    crypto.randomBytes(4).copy(buf, 12);
    // info hash
    parser.infoHash(torrent).copy(buf, 16);
    // peerId
    util.genId().copy(buf, 36);
    // downloaded
    Buffer.alloc(8).copy(buf, 56);
    // left
    parser.size(torrent).copy(buf, 64);
    // uploaded
    Buffer.alloc(8).copy(buf, 72);
    // event
    buf.writeUInt32BE(0, 80);
    // ip address
    buf.writeUInt32BE(0, 84);
    // key
    crypto.randomBytes(4).copy(buf, 88);
    // num want
    buf.writeInt32BE(-1, 92);
    // port
    buf.writeUInt16BE(port, 96);

    return buf;
}

/**
 * 
 * @param {Buffer} resp
 */
function parseAnnounceResp(resp) {
    function group(iterable, groupSize) {
        const groups = [];
        for (let i = 0; i < iterable.length; i += groupSize) {
            groups.push(iterable.slice(i, i + groupSize));
        }

        return groups;
    }
    return {
        action: resp.readUInt32BE(0),
        transactionId: resp.readUInt32BE(4),
        leechers: resp.readUInt32BE(8),
        seeders: resp.readUInt32BE(12),
        peers: group(resp.slice(20), 6).map(address => ({
            ip: address.slice(0, 4).join('.'),
            port: address.readUInt16BE(4)
        }))
    }
}


module.exports = {
    getPeers
};
