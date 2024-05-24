package chat

import "time"

// Time allowed to write a message to the peer.
const WRITE_WAIT = 10 * time.Second

// Time allowed to read the next pong message from the peer.
const PONG_WAIT = 60 * time.Second

// Send pings to peer with this period. Must be less than pongWait.
const PING_PERIOD = (PONG_WAIT * 9) / 10

// Maximum message size allowed from peer.
const MAX_MESSAGE_SIZE = 512
