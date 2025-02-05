package utils

import (
	"fmt"

	"github.com/tronprotocol/tron-docker/config"
)

func ShowSnapshotDataSourceList() {
	fmt.Printf("\nMain network - Lite Fullnode Data Source: \n")
	for _, items := range config.SnapshotDataSource[config.STLiteLevelSG] {
		fmt.Printf("  Region: %s\n", items.Region)
		fmt.Printf("    DBType: %s\n", items.DBType)
		fmt.Printf("    Domain: %s\n", items.Domain)
		fmt.Printf("    Description: %s\n\n", items.Description)
	}

	fmt.Printf("Main network - Fullnode Data Source: \n")
	for _, items := range config.SnapshotDataSource[config.STFullLevelSG] {
		fmt.Printf("  Region: %s\n", items.Region)
		fmt.Printf("    DBType: %s\n", items.DBType)
		fmt.Printf("    Domain: %s\n", items.Domain)
		fmt.Printf("    Description: %s\n\n", items.Description)
	}
	for _, items := range config.SnapshotDataSource[config.STFullLevelNA] {
		fmt.Printf("  Region: %s\n", items.Region)
		fmt.Printf("    DBType: %s\n", items.DBType)
		fmt.Printf("    Domain: %s\n", items.Domain)
		fmt.Printf("    Description: %s\n\n", items.Description)
	}
	for _, items := range config.SnapshotDataSource[config.STFullLevelNAWithAccountHistory] {
		fmt.Printf("  Region: %s\n", items.Region)
		fmt.Printf("    DBType: %s\n", items.DBType)
		fmt.Printf("    Domain: %s\n", items.Domain)
		fmt.Printf("    Description: %s\n\n", items.Description)
	}
	for _, items := range config.SnapshotDataSource[config.STFullRocksSG] {
		fmt.Printf("  Region: %s\n", items.Region)
		fmt.Printf("    DBType: %s\n", items.DBType)
		fmt.Printf("    Domain: %s\n", items.Domain)
		fmt.Printf("    Description: %s\n\n", items.Description)
	}

	fmt.Printf("Nile network - Fullnode/Lite Fullnode Data Source: \n")
	for _, items := range config.SnapshotDataSource[config.STNileLevel] {
		fmt.Printf("  Region: %s\n", items.Region)
		fmt.Printf("    DBType: %s\n", items.DBType)
		fmt.Printf("    Domain: %s\n", items.Domain)
		fmt.Printf("    Description: %s\n\n", items.Description)
	}
}

func CheckDomain(domain string) bool {
	has := false

	for _, items := range config.SnapshotDataSource {
		for _, item := range items {
			if item.Domain == domain {
				return true
			}
		}
	}

	return has
}

func IsHttps(domain string) bool {
	has := false

	for mtype, items := range config.SnapshotDataSource {
		for _, item := range items {
			if item.Domain == domain {
				return mtype == config.STNileLevel
			}
		}
	}

	return has
}
