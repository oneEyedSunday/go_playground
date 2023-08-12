###

Code from [article](https://allenkim67.github.io/programming/2016/05/04/how-to-make-your-own-bittorrent-client.html#grouping-messages) and [repo](https://github.com/allenkim67/allen-torrent/blob/master)


### Improvments

TODO:

    - Add a graphic user interface
    - Optimize for better download speeds and more efficient cpu usage. For example some clients calculate which pieces are the rarest and download those first.
    - There’s also something called distributed hash tables which makes it possible to share torrents without the use of centralized trackers.
    - You could write code to reconnect dropped connections
    - You could look for more peers periodically.
    - You could support pausing and resuming downloads.
    - You could support uploading since currently our client only downloads.
    - Sometimes peers are unable to connect to each other because they are behind a NAT which gives a proxy ip. You could look into NAT traversal strategies.
    It’s possible to bring bittorrent to the web using WebRTC, which creates a direct peer to peer connection.
