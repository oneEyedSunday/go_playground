'use strict';

const tracker = require('./tracker');
const parser = require('./parser');

const torrent = parser.open('puppy.torrent');

tracker.getPeers(torrent, peers => {
    console.log('list of peers: ', peers);
});
