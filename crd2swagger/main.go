package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/HaesungSeo/yamlconv"
	"gopkg.in/yaml.v2"
)

func crdObjects(data interface{}) (objs map[string]interface{}, err error) {
	// check the it confirms crd definitions
	fapiVersion := "apiVersion"
	fkind := "kind"
	fmetaName := "metadata.name"
	fspec := "spec"

	_, err = yamlconv.Search(data, strings.Split(fapiVersion, "."))
	if err != nil {
		return nil, fmt.Errorf("NOT FOUND: %s: %+v", fapiVersion, err)
	}
	_, err = yamlconv.Search(data, strings.Split(fkind, "."))
	if err != nil {
		return nil, fmt.Errorf("NOT FOUND: %s: %+vn", fkind, err)
	}
	mmeta, err := yamlconv.Search(data, strings.Split(fmetaName, "."))
	if err != nil {
		return nil, fmt.Errorf("NOT FOUND: %s: %+vn", fmetaName, err)
	}
	crdName, ok := mmeta.(string)
	if !ok {
		return nil, fmt.Errorf("CAST ERROR: %s", fmetaName)
	}
	_, err = yamlconv.Search(data, strings.Split(fspec, "."))
	if err != nil {
		return nil, fmt.Errorf("NOT FOUND: %s: %+v", fspec, err)
	}
	objs = make(map[string]interface{}, 0)

	// Build CR Struct from CRD
	// name: .spec.names.singular
	//   apiVersion: .spec.group '/' .spec.versions[].name
	//   kind: .spec.names.kind
	//   spec: .spec.versions[].schema.openAPIV3Schema.properties.spec.properties
	fSpecGrp := "spec.group"
	fspecKind := "spec.names.kind"
	fspecName := "spec.names.singular"
	fspecVers := "spec.versions"
	fscheme := "schema.openAPIV3Schema.properties.spec.properties"

	mSpecGrp, err := yamlconv.Search(data, strings.Split(fSpecGrp, "."))
	if err != nil {
		return nil, fmt.Errorf("NOT FOUND: %s: %s: %+v", crdName, fSpecGrp, err)
	}
	mkind, err := yamlconv.Search(data, strings.Split(fspecKind, "."))
	if err != nil {
		return nil, fmt.Errorf("NOT FOUND: %s: %s: %+v", crdName, fspecKind, err)
	}
	mname, err := yamlconv.Search(data, strings.Split(fspecName, "."))
	if err != nil {
		return nil, fmt.Errorf("NOT FOUND: %s: %s: %+v", crdName, fspecName, err)
	}
	mvers, err := yamlconv.Search(data, strings.Split(fspecVers, "."))
	if err != nil {
		return nil, fmt.Errorf("NOT FOUND: %s: %s: %+v", crdName, fspecVers, err)
	}
	vers, ok := mvers.([]interface{})
	if !ok {
		return nil, fmt.Errorf("CAST ERROR: %s: %s: %+v", crdName, fspecVers, err)
	}

	for ii, versionObj := range vers {
		vObj, ok := versionObj.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("CAST ERROR: %s: spec.versions[%d]", crdName, ii)
		}
		mspec, err := yamlconv.Search(versionObj, strings.Split(fscheme, "."))
		if err != nil {
			// don't panic, its common case
			// return nil, fmt.Errorf("NOT FOUND: %s: [%d]%s: %+v", crdName, ii, fscheme, err)
			continue
		}
		switch m := mspec.(type) {
		case map[interface{}]interface{}:
			name, _ := mname.(string)
			apiVersion := fmt.Sprintf("%s.%s", mSpecGrp, vObj["name"])

			// build apiVersion enum
			apiObj := make(map[string]interface{})
			apiObj["type"] = "string"
			apiObj["enum"] = []string{apiVersion}

			// build kind enum
			kindObj := make(map[string]interface{})
			kindObj["type"] = "string"
			kindObj["enum"] = []string{fmt.Sprintf("%s", mkind)}

			// build spec enum
			specObj := make(map[string]interface{})
			specObj["type"] = "object"
			specObj["properties"] = m

			// build obj properties
			properties := make(map[string]interface{}, 0)
			properties["apiVersion"] = apiObj
			properties["kind"] = kindObj
			properties["spec"] = specObj

			// build obj, finally
			outObj := make(map[string]interface{}, 0)
			outObj["type"] = "object"
			outObj["properties"] = properties

			objName := name + "." + apiVersion
			objs[objName] = outObj
		}
	}

	return objs, nil
}

type schema struct {
	Ref string `yaml:"$ref"`
}

type contentBody struct {
	Schema schema `yaml:"schema"`
}

type resp struct {
	Desc    string                 `yaml:"description"`
	Content map[string]contentBody `yaml:"content"`
}

type responses struct {
	Responses map[int]resp `yaml:"responses"`
}

type method struct {
	Get responses `yaml:"get"`
}

func swaggerPrint(swagger map[string]interface{}) {
	buf, err := yaml.Marshal(swagger)
	if err != nil {
		panic(fmt.Sprintf("yaml marshal: %s", err.Error()))
	}
	os.Stdout.Write(buf)
}

func crdPrint(objs map[string]interface{}) {
	// build contact
	contact := make(map[string]string, 0)
	contact["name"] = "a"
	contact["url"] = "http://localhost"
	contact["email"] = "a@a.com"

	// build info
	info := make(map[string]interface{}, 0)
	info["title"] = "CRD to REST-API"
	info["description"] = "CRD to REST-API conversion"
	info["version"] = "V3.0.0-oas4"
	info["contact"] = contact

	// print openapi
	swagger := make(map[string]interface{}, 0)
	swagger["openapi"] = "3.0.0"
	swaggerPrint(swagger)

	// print info
	swagger = make(map[string]interface{}, 0)
	swagger["info"] = info
	swaggerPrint(swagger)

	// build paths
	paths := make(map[string]method, 0)
	for k, _ := range objs {
		jsonBody := contentBody{
			Schema: schema{
				Ref: fmt.Sprintf("#/components/schemas/%s", k),
			},
		}
		rsp := resp{
			Desc:    "",
			Content: make(map[string]contentBody, 0),
		}
		rsp.Content["application/json"] = jsonBody
		get := method{
			Get: responses{
				Responses: make(map[int]resp),
			},
		}
		get.Get.Responses[200] = rsp
		paths["/"+k] = get
	}

	// print paths
	swagger = make(map[string]interface{}, 0)
	swagger["paths"] = paths
	swaggerPrint(swagger)

	// build components
	components := make(map[string]interface{}, 0)
	components["schemas"] = objs
	swagger = make(map[string]interface{}, 0)
	swagger["components"] = components
	swaggerPrint(swagger)
}

func main() {
	yamlpath := flag.String("f", "/dev/stdin", "yaml file")
	flag.Parse()

	// read yaml file into buffer
	var filebuf []byte
	filename, err := filepath.Abs(*yamlpath)
	if err != nil {
		panic(err.Error())
	}
	filebuf, err = os.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}

	// parse it
	var data interface{}
	err = yaml.Unmarshal(filebuf, &data)
	if err != nil {
		panic(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	objs := make(map[string]interface{}, 0)

	// check it contains array of yamls
	switch arr := data.(type) {
	case []interface{}:
		for ii, item := range arr {
			obj, err := crdObjects(item)
			if err != nil {
				panic(fmt.Sprintf("ERROR: yaml[%d] %s", ii, err.Error()))
			}
			for k, v := range obj {
				objs[k] = v
			}
		}
	case interface{}:
		obj, err := crdObjects(data)
		if err != nil {
			panic(fmt.Sprintf("ERROR: %s", err.Error()))
		}
		for k, v := range obj {
			objs[k] = v
		}
	default:
		panic(fmt.Sprintf("ERROR: unknown yaml type %T", arr))
	}

	crdPrint(objs)
}
