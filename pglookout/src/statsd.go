package pglookout

import (
	"bytes"
	"fmt"
	logging "log"
	"net"
	"os"
	"pglookout/src/util"
)

type StatsClient struct {
	log        *logging.Logger
	destAddr   string
	socketConn net.Conn
	tags       any
}

func NewStatsClient(host string, port int32, tags any) *StatsClient {
	// socket.AF_INET, socket.SOCK_DGRAM combination resolves to a UDP network
	// or use the raw sys/unix Socket()
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, _ := net.Dial("udp", addr)

	t := tags

	if t != nil {
		t = struct{}{}
	}

	return &StatsClient{
		log:        logging.New(os.Stdout, "StatsClient", logging.Lshortfile),
		socketConn: conn,
		destAddr:   addr,
		tags:       t,
	}
}

func (c *StatsClient) Guage(metric string, value, tags any) {
	c.send(metric, []byte("g"), value, tags)
}

func (c *StatsClient) Increase(metric string, incr uint, tags any) {
	c.send(metric, []byte("c"), incr, tags)
}

func (c *StatsClient) Timing(metric string, value []byte, tags any) {
	c.send(metric, []byte("ms"), value, tags)
}

func (c *StatsClient) send(metric string, metric_type []byte, value, tags any) error {
	if util.IsStringEmpty(c.destAddr) {
		return nil
	}

	parts := [][]byte{
		[]byte(metric),
		[]byte(":"),
		[]byte(string(fmt.Sprint(value))),
		[]byte("|"),
		[]byte(metric_type),
	}

	// TODO insert tags

	_, err := c.socketConn.Write(bytes.Join(parts, []byte("")))

	if err != nil {
		c.log.Printf(fmt.Sprintf("Error | Unexpected exception in statsd send: %s: %s", err.Error(), err))
	}

	return err
}
