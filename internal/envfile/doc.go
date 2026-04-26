// Package envfile provides primitives for parsing and representing .env files.
//
// A .env file is a newline-delimited list of KEY=VALUE pairs, optionally
// containing blank lines and comments (lines starting with '#').
//
// Values may be quoted with single or double quotes, and may include an
// inline comment separated by " #".
//
// Example usage:
//
//	f, err := os.Open(".env")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer f.Close()
//
//	ef, err := envfile.Parse(f)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, entry := range ef.Entries {
//		if entry.Key == "" {
//			continue // blank line or comment
//		}
//		fmt.Printf("%s = %s\n", entry.Key, entry.Value)
//	}
//
// The parsed EnvFile preserves entry order and provides an index for O(1)
// key lookups, making it suitable for diff and merge operations.
package envfile
