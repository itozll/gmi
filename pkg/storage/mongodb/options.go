package mongodb

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type ClientOptions interface {
	ClientOptions() *options.ClientOptions
	GetDatabase() string
}

type Options struct {
	URI                    string        `json:"uri" yaml:"uri" toml:"uri"`
	Database               string        `json:"database" yaml:"database" toml:"database"`
	ConnectTimeout         time.Duration `json:"connect_timeout,omitempty" yaml:"connect_timeout" toml:"connect_timeout"`
	HeartbeatInterval      time.Duration `json:"heartbeat_interval,omitempty" yaml:"heartbeat_interval" toml:"heartbeat_interval"`
	LocalThreshold         time.Duration `json:"local_threshold,omitempty" yaml:"local_threshold" toml:"local_threshold"`
	MaxConnIdleTime        time.Duration `json:"max_conn_idletime,omitempty" yaml:"max_conn_idletime" toml:"max_conn_idletime"`
	MaxPoolSize            uint64        `json:"max_pool_size,omitempty" yaml:"max_pool_size" toml:"max_pool_size"`
	MinPoolSize            uint64        `json:"min_pool_size,omitempty" yaml:"min_pool_size" toml:"min_pool_size"`
	ServerSelectionTimeout time.Duration `json:"server_selection_timeout,omitempty" yaml:"server_selection_timeout" toml:"server_selection_timeout"`
	SocketTimeout          time.Duration `json:"socket_timeout,omitempty" yaml:"socket_timeout" toml:"socket_timeout"`
	ReadPref               string        `json:"read_pref,omitempty" yaml:"read_pref" toml:"read_pref"`
	ReadConcern            string        `json:"read_concern,omitempty" yaml:"read_concern"`
}

var (
	DefaultReadPreferred = readpref.SecondaryPreferred()

	readPrefMap = map[string]*readpref.ReadPref{
		"primaryPreferred":   readpref.PrimaryPreferred(),
		"secondary":          readpref.Secondary(),
		"secondaryPreferred": readpref.SecondaryPreferred(),
		"nearest":            readpref.Nearest(),
	}

	readConcernMap = map[string]*readconcern.ReadConcern{
		"local":        readconcern.Local(),
		"majority":     readconcern.Majority(),
		"linearizable": readconcern.Linearizable(),
		"available":    readconcern.Available(),
		"snapshot":     readconcern.Snapshot(),
	}
)

func (o *Options) ClientOptions() *options.ClientOptions {
	opt := options.Client().ApplyURI(o.URI).
		SetLocalThreshold(o.LocalThreshold).
		SetMaxConnIdleTime(o.MaxConnIdleTime).
		SetMaxPoolSize(o.MaxPoolSize).
		SetConnectTimeout(o.ConnectTimeout).
		SetServerSelectionTimeout(o.ServerSelectionTimeout).
		SetSocketTimeout(o.SocketTimeout)

	if o.HeartbeatInterval > 0 {
		opt.SetHeartbeatInterval(o.HeartbeatInterval)
	}

	if o.MinPoolSize > 0 {
		opt.SetMinPoolSize(o.MinPoolSize)
	}

	// https://docs.mongodb.com/manual/core/read-preference/
	if pref, ok := readPrefMap[strings.ToLower(strings.TrimSpace(o.ReadPref))]; ok {
		opt.SetReadPreference(pref)
	} else {
		opt.SetReadPreference(DefaultReadPreferred)
	}

	if concern, ok := readConcernMap[strings.ToLower(strings.TrimSpace(o.ReadConcern))]; ok {
		opt.SetReadConcern(concern)
	}

	return opt
}

func (o *Options) GetDatabase() string {
	return o.Database
}
