package uidgo

import (
	"fmt"
	"sync"
	"time"
)

// ref: https://en.wikipedia.org/wiki/Snowflake_ID

var (
	// set the beginning time
	epoch = time.Date(time.Now().Year(), time.January, 01, 00, 00, 00, 00, time.UTC).UnixMilli()
)

const (
	// timestamp occupancy bits
	timestampBits = 41
	// dataCenterId occupancy bits
	dataCenterIdBits = 5
	// workerId occupancy bits
	workerIdBits = 5
	// sequence occupancy bits
	seqBits = 12

	// timestamp max value, just like 2^41-1 = 2199023255551
	timestampMaxValue = (1 << timestampBits) - 1
	// dataCenterId max value, just like 2^5-1 = 31
	dataCenterIdMaxValue = (1 << dataCenterIdBits) - 1
	// workId max value, just like 2^5-1 = 31
	workerIdMaxValue = (1 << workerIdBits) - 1
	// sequence max value, just like 2^12-1 = 4095
	seqMaxValue = (1 << seqBits) - 1

	// number of workId offsets (seqBits)
	workIdShift = 12
	// number of dataCenterId offsets (seqBits + workerIdBits)
	dataCenterIdShift = 17
	// number of timestamp offsets (seqBits + workerIdBits + dataCenterIdBits)
	timestampShift = 22

	defaultInitValue = 0
)

type SnowflakeSeqGenerator struct {
	timestamp    int64
	dataCenterId int64
	workerId     int64
	sequence     int64
	mu           *sync.Mutex
}

// NewSnowflakeSeqGenerator initiates the snowflake generator
func NewSnowflakeSeqGenerator(dataCenterId, workId int64) (r *SnowflakeSeqGenerator, err error) {
	if dataCenterId < 0 || dataCenterId > dataCenterIdMaxValue {
		err = fmt.Errorf("dataCenterId should between 0 and %d", dataCenterIdMaxValue-1)
		return nil, err
	}

	if workId < 0 || workId > workerIdMaxValue {
		err = fmt.Errorf("workId should between 0 and %d", dataCenterIdMaxValue-1)
		return nil, err
	}

	return &SnowflakeSeqGenerator{
		mu:           new(sync.Mutex),
		timestamp:    defaultInitValue - 1,
		dataCenterId: dataCenterId,
		workerId:     workId,
		sequence:     defaultInitValue,
	}, nil
}

// GenerateId timestamp + dataCenterId + workId + sequence
func (S *SnowflakeSeqGenerator) GenerateId1() (string, error) {
	S.mu.Lock()
	defer S.mu.Unlock()

	now := time.Now().UnixMilli()

	if S.timestamp > now { // Clock callback
		return "", fmt.Errorf("Clock moved backwards. Refusing to generate ID, last timestamp is %d, now is %d", S.timestamp, now)
	} else if S.timestamp == now {
		// generate multiple IDs in the same millisecond, incrementing the sequence number to prevent conflicts
		S.sequence = (S.sequence + 1) & seqMaxValue
		if S.sequence == 0 {
			// sequence overflow, waiting for next millisecond
			for now <= S.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// initialized sequences are used directly at different millisecond timestamps
		S.sequence = defaultInitValue
	}
	tmp := now - epoch
	if tmp > timestampMaxValue {
		return "", fmt.Errorf("epoch should between 0 and %d", timestampMaxValue-1)
	}
	S.timestamp = now

	// combine the parts to generate the final ID and convert the 64-bit binary to decimal digits.
	r := (tmp)<<timestampShift |
		(S.dataCenterId << dataCenterIdShift) |
		(S.workerId << workIdShift) |
		(S.sequence)

	return fmt.Sprintf("%d", r), nil
}

func (S *SnowflakeSeqGenerator) GenerateId2() (uint64, error) {
	S.mu.Lock()
	defer S.mu.Unlock()

	now := time.Now().UnixMilli()

	if S.timestamp > now { // Clock callback
		return 0, fmt.Errorf("Clock moved backwards. Refusing to generate ID, last timestamp is %d, now is %d", S.timestamp, now)
	} else if S.timestamp == now {
		// generate multiple IDs in the same millisecond, incrementing the sequence number to prevent conflicts
		S.sequence = (S.sequence + 1) & seqMaxValue
		if S.sequence == 0 {
			// sequence overflow, waiting for next millisecond
			for now <= S.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// initialized sequences are used directly at different millisecond timestamps
		S.sequence = defaultInitValue
	}
	tmp := now - epoch
	if tmp > timestampMaxValue {
		return 0, fmt.Errorf("epoch should between 0 and %d", timestampMaxValue-1)
	}
	S.timestamp = now

	// combine the parts to generate the final ID and convert the 64-bit binary to decimal digits.
	r := (tmp)<<timestampShift |
		(S.dataCenterId << dataCenterIdShift) |
		(S.workerId << workIdShift) |
		(S.sequence)

	return uint64(r), nil
}

func (S *SnowflakeSeqGenerator) GenerateId3() (uint64, string, error) {
	S.mu.Lock()
	defer S.mu.Unlock()

	now := time.Now().UnixMilli()

	if S.timestamp > now { // Clock callback
		return 0, "", fmt.Errorf("Clock moved backwards. Refusing to generate ID, last timestamp is %d, now is %d", S.timestamp, now)
	} else if S.timestamp == now {
		// generate multiple IDs in the same millisecond, incrementing the sequence number to prevent conflicts
		S.sequence = (S.sequence + 1) & seqMaxValue
		if S.sequence == 0 {
			// sequence overflow, waiting for next millisecond
			for now <= S.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// initialized sequences are used directly at different millisecond timestamps
		S.sequence = defaultInitValue
	}
	tmp := now - epoch
	if tmp > timestampMaxValue {
		return 0, "", fmt.Errorf("epoch should between 0 and %d", timestampMaxValue-1)
	}
	S.timestamp = now

	// combine the parts to generate the final ID and convert the 64-bit binary to decimal digits.
	r := (tmp)<<timestampShift |
		(S.dataCenterId << dataCenterIdShift) |
		(S.workerId << workIdShift) |
		(S.sequence)

	return uint64(r), fmt.Sprintf("%d", r), nil
}
