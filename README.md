yamlconv
---
Print & Search Yaml struct in go
- [import in your project](#import-in-your-project)
- [example of printing yaml struct](#example-of-printing-yaml-struct)
  - [yaml example](#yaml-example)
  - [sample code](#sample-code)
  - [run result](#run-result)
- [example of searching keys in yaml struct](#example-of-searching-keys-in-yaml-struct)
  - [yaml example](#yaml-example-1)
  - [sample code](#sample-code-1)
  - [run result](#run-result-1)
- [example of subtract some struct in yaml struct](#example-of-subtract-some-struct-in-yaml-struct)
  - [yaml example](#yaml-example-2)
  - [sample code](#sample-code-2)
  - [run result](#run-result-2)


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
lets go and play with [playground](https://go.dev/play/p/gPU6zdlDP7R)
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

## run result
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
lets go and play with [playground](https://go.dev/play/p/85ICpvMjTua)
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

## run result
```
$ ./search
 Str[10.10.0.101]
$
```

# example of subtract some struct in yaml struct
## yaml example
lets go and play with [playground](https://go.dev/play/p/85ICpvMjTua)
```
---
sriov:
  - network: resource01 # network name
    interface: net1
    ip: 10.10.0.101     # network ip
  - network: resource02 # network name
    interface: net2
    ip: 20.10.0.101     # network ip
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
remove the struct under the keys "sriov[1]", "sriov[0].network" and "password" in above yaml sample<br>
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
        "  - network: resource02 # network name",
        "    interface: net2",
        "    ip: 20.10.0.101     # network ip",
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
    ret, err := yamlconv.Subtract(data, []string{"sriov", "[0]"})
    if err != nil {
        panic(err.Error())
    }
    ret, err = yamlconv.Subtract(ret, []string{"sriov", "[0]", "network"})
    if err != nil {
        panic(err.Error())
    }
    ret, err = yamlconv.Subtract(ret, []string{"password"})
    if err != nil {
        panic(err.Error())
    }
    yamlconv.Print(ret, "  ")
    fmt.Printf("\n")
}
```

## run result
```
$ ./subtract

M[sriov]
  A[0/1]
    M[interface] Str{net2}
    M[ip] Str{20.10.0.101}
M[gpu]
  M[drivers] Str{video,compute,utility}
M[service]
  M[type]
    M[NodePort] Int{30080}
M[chpasswd]
  M[expire] Bool{false}
M[ssh_pwauth] Bool{true}
$
```
