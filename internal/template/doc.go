// Package template implements variable substitution for .env file values.
//
// It supports two placeholder syntaxes commonly found in configuration files:
//
//   - Double-brace style:  {{VAR_NAME}}
//   - Shell style:         ${VAR_NAME}
//
// Placeholders are resolved against a caller-supplied context map (typically
// the current environment or another .env file). Unresolved placeholders are
// left unchanged in the output and reported as [ErrUnresolved] errors so that
// callers can decide whether to treat them as fatal or as warnings.
//
// # Usage
//
//	ctx := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
//	out, errs := template.Expand(env, ctx)
//	if len(errs) > 0 {
//		// handle or log unresolved placeholders
//	}
//
// The package is intentionally stateless; each call to Expand returns a new
// map and does not modify the input.
package template
