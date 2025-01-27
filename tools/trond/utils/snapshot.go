package utils

import (
	"fmt"

	"github.com/tronprotocol/tron-docker/config"
)

func ShowSnapshotDataSourceList() {
	fmt.Printf("\nLite Fullnode Data Source: \n")
	for _, items := range config.SnapshotDataSource[config.STLiteLevelSG] {
		fmt.Printf("  Region: %s\n", items.Region)
		fmt.Printf("    DBType: %s\n", items.DBType)
		fmt.Printf("    Host: %s\n", items.Host)
		fmt.Printf("    Description: %s\n\n", items.Description)
	}

	fmt.Printf("Fullnode Data Source: \n")
	for _, items := range config.SnapshotDataSource[config.STFullLevelSG] {
		fmt.Printf("  Region: %s\n", items.Region)
		fmt.Printf("    DBType: %s\n", items.DBType)
		fmt.Printf("    Host: %s\n", items.Host)
		fmt.Printf("    Description: %s\n\n", items.Description)
	}
	for _, items := range config.SnapshotDataSource[config.STFullLevelNA] {
		fmt.Printf("  Region: %s\n", items.Region)
		fmt.Printf("    DBType: %s\n", items.DBType)
		fmt.Printf("    Host: %s\n", items.Host)
		fmt.Printf("    Description: %s\n\n", items.Description)
	}
	for _, items := range config.SnapshotDataSource[config.STFullLevelNAWithAccountHistory] {
		fmt.Printf("  Region: %s\n", items.Region)
		fmt.Printf("    DBType: %s\n", items.DBType)
		fmt.Printf("    Host: %s\n", items.Host)
		fmt.Printf("    Description: %s\n\n", items.Description)
	}
	for _, items := range config.SnapshotDataSource[config.STFullRocksSG] {
		fmt.Printf("  Region: %s\n", items.Region)
		fmt.Printf("    DBType: %s\n", items.DBType)
		fmt.Printf("    Host: %s\n", items.Host)
		fmt.Printf("    Description: %s\n\n", items.Description)
	}
}

func CheckDomain(domain string) bool {
	has := false

	for _, items := range config.SnapshotDataSource {
		for _, item := range items {
			if item.Host == domain {
				return true
			}
		}
	}

	return has
}
