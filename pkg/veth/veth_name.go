package veth

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	hostVethPairPrefix = "veth"
)

func VethNameForWorkload(namespace, podname string) string {
	// A SHA1 is always 20 bytes long, and so is sufficient for generating the
	// veth name and mac addr.
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s.%s", namespace, podname)))
	return fmt.Sprintf("%s%s", hostVethPairPrefix, hex.EncodeToString(h.Sum(nil))[:11])
}

type CniArgs struct {
	Namespace   string
	PodName     string
	ContainerID string
}

func ParseArgs(args string) *CniArgs {
	// FIXME: 测试用
	if args == "" {
		return &CniArgs{
			Namespace:   "na",
			PodName:     "pod1234",
			ContainerID: "con1234",
		}
	}

	m := make(map[string]string)

	attrs := strings.Split(args, ";")


	for _, attr := range attrs {
		kv := strings.Split(attr, "=")
		if len(kv) != 2 {
			continue
		}

		m[kv[0]] = kv[1]
	}

	return &CniArgs{
		Namespace:   m["K8S_POD_NAMESPACE"],
		PodName:     m["K8S_POD_NAME"],
		ContainerID: m["K8S_POD_INFRA_CONTAINER_ID"],
	}
}
