'use strict';

const download = require('./download');
const parser = require('./parser');

const torrent = parser.open('puppy.torrent');

download(torrent, torrent.info.name);
