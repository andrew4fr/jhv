package server

var (
	globalConfig serviceConfig
)

type serviceConfig struct {
	XMLPath    string `long:"xml-path" env:"XML_PATH" required:"true" default:"https://www.treasury.gov/ofac/downloads/sdn.xml"`
	StorageDSN string `long:"storage-dsn" env:"STORAGE_DSN" required:"true"`
}
