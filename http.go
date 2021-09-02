package fy

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func RequestString(req *http.Request) string {
	var build strings.Builder
	build.WriteString(">>>>>>>>> Request >>>>>>>>>\n")
	if req == nil {
		return ""
	}
	build.WriteString(fmt.Sprintf("%s %s %s\n", req.Method, req.URL.Path, req.Proto))
	build.WriteString("Host: ")
	build.WriteString(req.URL.Host)
	build.WriteByte('\n')
	for k, vs := range req.Header {
		build.WriteString(k)
		build.WriteString(": ")
		for _, v := range vs {
			build.WriteString(v)
		}
		build.WriteByte('\n')
	}
	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return err.Error()
		}
		build.Write(body)
	}
	build.WriteString("\n<<<<<<<<< Request <<<<<<<<<\n\n")
	return build.String()
}

func ResponseStr(rsp *http.Response, body []byte) string {
	var build strings.Builder
	build.WriteString(">>>>>>>>> Response >>>>>>>>>\n")
	if rsp == nil {
		return ""
	}
	build.WriteString(rsp.Proto)
	build.WriteByte(' ')
	build.WriteString(rsp.Status)
	build.WriteByte('\n')
	for k, vs := range rsp.Header {
		build.WriteString(k)
		build.WriteString(": ")
		for _, v := range vs {
			build.WriteString(v)
		}
		build.WriteByte('\n')
	}

	build.Write(body)
	build.WriteString("\n<<<<<<<<< Response <<<<<<<<<\n\n")
	return build.String()
}

func ResponseString(rsp *http.Response) string {
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return err.Error()
	}
	return ResponseStr(rsp, body)
}
