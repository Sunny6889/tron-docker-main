package config

import (
	"fmt"
)

type SnapshotType int

const (
	STFullLevelNA SnapshotType = iota // 1
	STFullLevelSG
	STFullLevelNAWithAccountHistory
	STFullRocksSG
	STLiteLevelSG
	STNileLevel
)

type SnapshotDataSourceItem struct {
	DataType    string
	DBType      string
	Region      string
	Domain      string
	Description string
}

// https://nile-snapshots.s3-accelerate.amazonaws.com/backup20250204/FullNode_output-directory.tgz
// http://34.143.247.77/backup20250205/FullNode_output-directory.tgz
var SnapshotDataSource = map[SnapshotType]map[int]SnapshotDataSourceItem{
	STFullLevelNA: {
		1: {
			DataType:    "Fullnode Data Source",
			DBType:      "LevelDB",
			Region:      "America",
			Domain:      "34.86.86.229",
			Description: "Exclude internal transactions (About 2094G on 25-Jan-2025)",
		},
	},
	STFullLevelSG: {
		1: {
			DataType:    "Fullnode Data Source",
			DBType:      "LevelDB",
			Region:      "Singapore",
			Domain:      "34.143.247.77",
			Description: "Exclude internal transactions (About 2093G on 24-Jan-2025)",
		},
		2: {
			DataType:    "Fullnode Data Source",
			DBType:      "LevelDB",
			Region:      "Singapore",
			Domain:      "35.247.128.170",
			Description: "Include internal transactions (About 2278G on 24-Jan-2025)",
		},
	},
	STFullLevelNAWithAccountHistory: {
		1: {
			DataType:    "Fullnode Data Source",
			DBType:      "LevelDB",
			Region:      "America",
			Domain:      "34.48.6.163",
			Description: "Exclude internal transactions, include account history TRX balance (About 2627G on 24-Jan-2025)",
		},
	},
	STFullRocksSG: {
		1: {
			DataType:    "Fullnode Data Source",
			DBType:      "RocksDB",
			Region:      "America",
			Domain:      "35.197.17.205",
			Description: "Exclude internal transactions (About 2067G on 24-Jan-2025)",
		},
	},
	STLiteLevelSG: {
		1: {
			DataType:    "Lite Fullnode Data Source",
			DBType:      "LevelDB",
			Region:      "Singapore",
			Domain:      "34.143.247.77",
			Description: "(About 46G on 24-Jan-2025)",
		},
	},
	STNileLevel: {
		1: {
			DataType:    "Data Source for Nile",
			DBType:      "LevelDB",
			Region:      "Singapore",
			Domain:      "database.nileex.io",
			Description: "Fullnode/Lite Fullnode",
		},
	},
}

// Define custom error type for "not supported" errors
type NotSupportedError struct {
	Name  string
	Value string
}

func (e *NotSupportedError) Error() string {
	return fmt.Sprintf("%s '%s' is not supported", e.Name, e.Value)
}
