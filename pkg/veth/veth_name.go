package veth

import (
	"strings"
)

const (
	hostVethPairPrefix = "veth"
)

type CniArgs struct {
	Namespace   string
	PodName     string
	ContainerID string
}

// ParseArgs 解析传入参数
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
