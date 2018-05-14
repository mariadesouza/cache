package redisproxy

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

//Redis : internal struct to manage redis connection
type redisServerConn struct {
	redisConnection net.Conn
	reader          *bufio.Reader
	writer          *bufio.Writer
}

var (
	arrayPrefixString = []byte{'*'}
	bulkPrefixString  = []byte{'$'}
	lineEndingBytes   = []byte{'\r', '\n'}
)

const (
	simpleStringPrefix = '+'
	integerPrefix      = ':'
	arrayPrefix        = '*'
	errorPrefix        = '-'
	bulkPrefix         = '$'
)

func newRedisConnection(redisServer string, port string) (*redisServerConn, error) {
	// connecting to Redis
	dialer := net.Dialer{KeepAlive: time.Minute * 5}
	var redis redisServerConn
	var err error
	redis.redisConnection, err = dialer.Dial("tcp", redisServer+":"+port)
	if err != nil {
		return nil, err
	}
	redis.writer = bufio.NewWriter(redis.redisConnection)
	redis.reader = bufio.NewReader(redis.redisConnection)
	return &redis, nil
}

func (r *redisServerConn) CloseConnection() error {
	err := r.redisConnection.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisServerConn) writeString(s string) error {
	_, err := r.writer.WriteString(s)
	if err != nil {
		return err
	}
	_, err = r.writer.Write(lineEndingBytes)
	return err
}

//Send : Send cmd and args to redis connection. cmds can be GET, SET
func (r *redisServerConn) Send(args ...string) error {

	_, err := r.writer.Write(arrayPrefixString)
	if err != nil {
		return err
	}

	if err = r.writeString(strconv.Itoa(len(args))); err != nil {
		return err
	}

	for _, arg := range args {

		_, err = r.writer.Write(bulkPrefixString)
		if err != nil {
			return err
		}

		if err := r.writeString(strconv.Itoa(len(arg))); err != nil {
			return err
		}

		if err := r.writeString(arg); err != nil {
			return err
		}
	}
	return r.writer.Flush()
}

//Receive: receive the response from redis
//RESP uses prefix so we first need to read the first byte and then decide how to handle the following bytes
func (r *redisServerConn) Receive() ([]byte, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return nil, fmt.Errorf("inadequate response line")
	}

	switch line[0] {
	case arrayPrefix:
		return r.readArray(line)
	case simpleStringPrefix, integerPrefix, errorPrefix:
		return line, nil
	case bulkPrefix:
		result, err := r.readBulkString(line)
		if err != nil {
			fmt.Println(err)
		} else {			fmt.Println(result)
		}
		return result, err
	default:
		return nil, errors.New("RESP: invalid syntax")
	}
}

func (r *redisServerConn) getLen(line []byte) (int, error) {
	if len(line) == 0 {
		return -1, fmt.Errorf("malformed length")
	}
	if line[0] == '-' && len(line) == 2 && line[1] == '1' { // null bulk string
		return -1, fmt.Errorf("returned null")
	}

	var n int
	for _, b := range line {
		n *= 10
		n += int(b - '0')
	}
	return n, nil
}

func (r *redisServerConn) readBulkString(line []byte) ([]byte, error) {

	count, err := r.getLen(line[1:])
	if err != nil {
		return nil, err
	}

	if count == -1 {
		return line, nil
	}
	buf := make([]byte, count)
	_, err = io.ReadFull(r.reader, buf) //read into buf
	if err != nil {
		return nil, err
	}

	line, err = r.readLine()
	if err != nil {
		return nil, err
	} else if len(line) != 0 {
		return nil, fmt.Errorf("bad bulk string format")
	}

	return buf, nil

}

// readArray : not tested
func (r *redisServerConn) readArray(line []byte) ([]byte, error) {
	end := bytes.IndexByte(line, '\r')
	count, _ := strconv.Atoi(string(line[1:end]))
	for i := 0; i < count; i++ {
		buf, err := r.Receive() //recursively read lines
		if err != nil {
			return nil, err
		}
		line = append(line, buf...)
	}
	return line, nil
}

// readLine reads one line till we find \r\n. It will return the line without the line endings
func (r *redisServerConn) readLine() ([]byte, error) {
	line, err := r.reader.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	if len(line) > 1 && line[len(line)-2] == '\r' {
		return line[:len(line)-2], nil
	}
	return nil, errors.New("resp: invalid syntax")

}
