yamlconv
---
Print & Search Yaml struct in go
- [import in your project](#import-in-your-project)
- [example of printing yaml struct](#example-of-printing-yaml-struct)
	- [yaml example](#yaml-example)
	- [sample code](#sample-code)
	- [run](#run)
- [example of searching keys in yaml struct](#example-of-searching-keys-in-yaml-struct)
	- [yaml example](#yaml-example-1)
	- [sample code](#sample-code-1)
	- [run](#run-1)


# import in your project
```go
import (
    "github.com/HaesungSeo/yamlconv"
)
```

# example of printing yaml struct
## yaml example
```
---
sriov:
  - network: resource01 # network name
    interface: net1
    ip: 10.10.0.101     # network ip
gpu:
 drivers: video,compute,utility

service:
  type:
    NodePort: 30080

#cloud-config
password: centos
chpasswd: { expire: False }
ssh_pwauth: True
```

## sample code
```
package main

import (
    "fmt"
    "strings"

    "github.com/HaesungSeo/yamlconv"
    "gopkg.in/yaml.v2"
)

func main() {
    buf := []string{
        "sriov:",
        "  - network: resource01 # network name",
        "    interface: net1",
        "    ip: 10.10.0.101     # network ip",
        "gpu:",
        " drivers: video,compute,utility",
        "service:",
        "  type:",
        "    NodePort: 30080",
        "#cloud-config",
        "password: centos",
        "chpasswd: { expire: False }",
        "ssh_pwauth: True",
    }

    data := yaml.MapSlice{}
    yaml.Unmarshal([]byte(strings.Join(buf, "\n")), &data)
    yamlconv.Print(data, "  ")
    fmt.Printf("\n")
}
```

## run
```
$ ./print

M[sriov]
  A[0/1]
    M[network] Str[resource01]
    M[interface] Str[net1]
    M[ip] Str[10.10.0.101]
M[gpu]
  M[drivers] Str[video,compute,utility]
M[service]
  M[type]
    M[NodePort] Int[30080]
M[password] Str[centos]
M[chpasswd]
  M[expire] Bool[false]
M[ssh_pwauth] Bool[true]
$
```

# example of searching keys in yaml struct
## yaml example
```
---
sriov:
  - network: resource01 # network name
    interface: net1
    ip: 10.10.0.101     # network ip
gpu:
 drivers: video,compute,utility

service:
  type:
    NodePort: 30080

#cloud-config
password: centos
chpasswd: { expire: False }
ssh_pwauth: True
```

## sample code
search the keys "sriov[0].ip" in above yaml sample<br>
expected result: "10.10.0.101"
```
package main

import (
	"fmt"
	"strings"

	"github.com/HaesungSeo/yamlconv"
	"gopkg.in/yaml.v2"
)

func main() {
	buf := []string{
		"sriov:",
		"  - network: resource01 # network name",
		"    interface: net1",
		"    ip: 10.10.0.101     # network ip",
		"gpu:",
		" drivers: video,compute,utility",
		"service:",
		"  type:",
		"    NodePort: 30080",
		"#cloud-config",
		"password: centos",
		"chpasswd: { expire: False }",
		"ssh_pwauth: True",
	}

	data := yaml.MapSlice{}
	yaml.Unmarshal([]byte(strings.Join(buf, "\n")), &data)
	ret, _ := yamlconv.Search(data, []string{"sriov", "[0]", "ip"})
	yamlconv.Print(ret, "  ")
	fmt.Printf("\n")
}
```

## run
```
$ ./search
 Str[10.10.0.101]
$
```
