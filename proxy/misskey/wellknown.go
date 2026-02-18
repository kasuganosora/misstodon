package misskey

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/models"
)

func NodeInfo(server string, ni models.NodeInfo) (models.NodeInfo, error) {
	var result models.NodeInfo
	_, err := client.R().
		SetResult(&result).
		Get(utils.JoinURL(server, "/nodeinfo/2.0"))
	if err != nil {
		return ni, err
	}
	ni.Usage = result.Usage
	ni.OpenRegistrations = result.OpenRegistrations
	ni.Metadata = result.Metadata
	return ni, err
}

func WebFinger(server, resource string, writer http.ResponseWriter) error {
	resp, err := client.R().
		SetDoNotParseResponse(true).
		SetQueryParam("resource", resource).
		Get(utils.JoinURL(server, "/.well-known/webfinger"))
	if err != nil {
		return err
	}
	defer resp.RawBody().Close()
	writer.Header().Set("Content-Type", resp.Header().Get("Content-Type"))
	writer.WriteHeader(resp.StatusCode())
	_, err = io.Copy(writer, resp.RawBody())
	return err
}

func HostMeta(server string, writer http.ResponseWriter) error {
	resp, err := client.R().
		SetDoNotParseResponse(true).
		Get(utils.JoinURL(server, "/.well-known/host-meta"))
	if err != nil {
		return err
	}
	defer resp.RawBody().Close()

	// Read body and replace server URL with proxy host
	body, err := io.ReadAll(resp.RawBody())
	if err != nil {
		return err
	}

	// Replace the Misskey server URL with proxy server URL
	result := strings.ReplaceAll(string(body), server, "")

	writer.Header().Set("Content-Type", resp.Header().Get("Content-Type"))
	writer.WriteHeader(resp.StatusCode())
	_, err = writer.Write([]byte(result))
	return err
}

// HostMetaXML returns a properly formatted host-meta XML with the proxy host
func HostMetaXML(proxyHost string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<XRD xmlns="http://docs.oasis-open.org/ns/xri/xrd-1.0">
  <Link rel="lrdd" template="https://%s/.well-known/webfinger?resource={uri}"/>
</XRD>
`, proxyHost)
}
