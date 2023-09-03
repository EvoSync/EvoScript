package EvoScript

import "bytes"

// ExecuteToString will attempt to execute the current to string representation
func ExecuteToString(source string, elements map[string]any) (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	if err := ExecuteString(source, buf, elements); err != nil {
		return "", err
	}

	return buf.String(), nil
}