// Package validate provides schema-based validation for .env files.
//
// It compares a target environment map against a schema (reference)
// environment map and reports missing required keys as errors and
// undeclared extra keys as warnings.
//
// Basic usage:
//
//	schema, _ := envfile.Parse(schemaReader)
//	target, _ := envfile.Parse(targetReader)
//
//	result, err := validate.Validate(schema, target)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, issue := range result.Issues {
//		fmt.Println(issue)
//	}
//	if result.HasErrors() {
//		os.Exit(1)
//	}
package validate
