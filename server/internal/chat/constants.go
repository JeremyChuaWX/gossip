package chat

import "time"

const WRITE_WAIT = 10 * time.Second

const PONG_WAIT = 60 * time.Second

const PING_PERIOD = PONG_WAIT * 9 / 10

const MAX_MESSAGE_SIZE = 10000

const BUFFER_SIZE = 4096
