package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func ValidateFlags(
	rawurl string,
	method string,
	data string,
	customHeaders []string,
) error {
	//rawurlのフォーマットをチェック
	if err := validateRawURL(rawurl); err != nil {
		return err
	}

	//methodの整合性をチェック
	if err := validateMethod(method); err != nil {
		return err
	}

	//dataのフォーマットをチェック
	if err := validateData(data); err != nil {
		return err
	}

	//customHeadersのフォーマットをチェック
	if err := validateHeader(customHeaders); err != nil {
		return err
	}

	return nil
}

func validateRawURL(rawurl string) error {
	url, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return err
	}
	if url.Scheme != "http" && url.Scheme != "https" {
		return fmt.Errorf("url schema '%s' is not supported", url.Scheme)
	}

	return nil
}

func validateMethod(method string) error {
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		return nil
	}
	return fmt.Errorf("HTTP method '%s' is not supported", method)
}

func validateData(data string) error {
	s := strings.TrimSpace(data)
	if s == "" {
		return nil
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return fmt.Errorf("json parse error: %s", err.Error())
	}
	return nil
}

func validateHeader(customHeaders []string) error {
	for _, v := range customHeaders {
		s := strings.TrimSpace(v)
		kv := strings.Split(s, ":")
		if len(kv) != 2 {
			return fmt.Errorf("invalid format header: %s", s)
		}
		if len(kv[0]) == 0 || len(kv[1]) == 0 {
			return fmt.Errorf("invalid format header: %s", s)
		}
	}
	return nil
}
