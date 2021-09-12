package merge

import (
	"time"
)

var DefaultOptions = &Options{
	Processor:     10,
	QueueSize:     1_000,
	Timeout:       10 * time.Second,
	WaitTimes:     3,
	TriggerNumber: 500,
	WaitNumber:    100,
}

type Options struct {
	// 处理进程数
	Processor int32 `json:"processor,omitempty" bson:"processor" toml:"processor" yaml:"processor"`

	// 队列大小
	QueueSize int32 `json:"queue_size,omitempty" bson:"queue_size" toml:"queue_size" yaml:"queue_size"`

	// 入队超时
	Timeout time.Duration `json:"timeout,omitempty" bson:"timeout" toml:"timeout" yaml:"timeout"`

	// 周期间隔时间
	Interval time.Duration `json:"interval,omitempty" bson:"interval" toml:"interval" yaml:"interval"`

	// 重试周期数
	TickCount int `json:"tick_count,omitempty" bson:"tick_count" toml:"tick_count" yaml:"tick_count"`

	// 触发操作的数量
	TriggerNumber int `json:"trigger_number,omitempty" bson:"trigger_number" toml:"trigger_number" yaml:"trigger_number"`

	// 周期等待次数
	WaitTimes int `json:"wait_times,omitempty" bson:"wait_times" toml:"wait_times" yaml:"wait_times"`

	// 周期中，用户不足时等待下个周期
	WaitNumber int `json:"wait_number,omitempty" bson:"wait_number" toml:"wait_number" yaml:"wait_number"`
}

func (opt *Options) init() *Options {
	if opt == nil {
		return DefaultOptions
	}

	if opt.Processor <= 0 {
		opt.Processor = DefaultOptions.Processor
	}

	if opt.QueueSize <= 0 {
		opt.QueueSize = opt.Processor * 100
	}

	if opt.Timeout <= 1*time.Millisecond {
		opt.Timeout = DefaultOptions.Timeout
	}

	if opt.WaitTimes <= 0 {
		opt.WaitTimes = DefaultOptions.WaitTimes
	}

	if opt.TriggerNumber <= 0 {
		opt.TriggerNumber = DefaultOptions.TriggerNumber
	}

	if opt.WaitNumber <= 0 {
		opt.WaitNumber = 60
	}

	return opt
}
