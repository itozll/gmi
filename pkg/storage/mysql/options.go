package mysql

type Options struct {
	Dsn                    string `json:"dsn,omitempty" yaml:"dsn" toml:"dsn"`
	Type                   string `json:"type,omitempty" yaml:"type" toml:"type"`
	PluralTable            bool   `json:"plural_table,omitempty" yaml:"plural_table"`
	SkipDefaultTransaction bool   `json:"skip_default_transaction,omitempty" yaml:"skip_default_transaction" toml:"skip_default_transaction"`
	Debug                  bool   `json:"debug,omitempty" yaml:"debug" toml:"debug"`
}

func DefaultOption() *Options { return &Options{} }

func (o *Options) WithPluralTable(pluralTable bool) *Options { o.PluralTable = pluralTable; return o }
func (o *Options) WithType(typ string) *Options              { o.Type = typ; return o }
func (o *Options) WithDebug(debug bool) *Options             { o.Debug = debug; return o }
func (o *Options) WithDsn(Dsn string) *Options               { o.Dsn = Dsn; return o }
func (o *Options) WithSkipDefaultTransaction(SkipDefaultTransaction bool) *Options {
	o.SkipDefaultTransaction = SkipDefaultTransaction
	return o
}
