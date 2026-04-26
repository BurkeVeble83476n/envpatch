// Package merge implements .env file merging for the envpatch utility.
//
// It supports combining a base environment file with an overlay environment
// file, which is useful when managing configuration across multiple
// environments (e.g. development, staging, production).
//
// # Strategies
//
// Two merge strategies are available:
//
//   - StrategyOverlay: values in the overlay take precedence over base values
//     when the same key exists in both files. New keys from the overlay are
//     always added.
//
//   - StrategyKeepBase: values in the base are preserved when a key conflict
//     occurs. New keys from the overlay are still added.
//
// # Usage
//
//	result, err := merge.Merge(baseMap, overlayMap, merge.StrategyOverlay)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result.Merged)
package merge
