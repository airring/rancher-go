package controllers

import (
	"crypto/tls"
	"encoding/json"
	"reflect"

	"github.com/astaxie/beego"
	"github.com/tidwall/gjson"
	"gopkg.in/resty.v1"
)

type RancherController struct {
	beego.Controller
}

type defaulturl struct {
	serverURL string
	token     string
}

type workretrun struct {
	Name           string
	Hostname       string
	Wtype          string
	Port           string
	Image          string
	Tcpport        []map[string]gjson.Result
	LivenessProbe  map[string]gjson.Result
	ReadinessProbe map[string]gjson.Result
	Environment    []map[string]interface{}
	Podgroup       []map[string]interface{}
}

func (client defaulturl) Getrancher(url string) (string, error) {
	client.token = beego.AppConfig.String("rancher_token")
	client.serverURL = beego.AppConfig.String("rancher_url")
	client1 := resty.New()
	client1.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client1.SetAuthToken(client.token)
	response, err := client1.R().Get(client.serverURL + "/v3" + url)
	if err != nil {
		return "", err
	}
	body := string(response.Body()[:])

	return body, nil
}

func (client defaulturl) Updaterancher(url string, path interface{}) (string, error) {
	client.token = beego.AppConfig.String("rancher_token")
	client.serverURL = beego.AppConfig.String("rancher_url")
	client1 := resty.New()
	client1.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client1.SetAuthToken(client.token)

	response, err := client1.R().SetBody(path).Put(client.serverURL + "/v3" + url)
	if err != nil {
		return "", err
	}
	body := string(response.Body()[:])
	beego.Info(reflect.TypeOf(body))
	return body, nil
}

func (u RancherController) Gettoken() {
	a := defaulturl{}
	var entities []map[string]interface{}
	// 获取集群信息
	response, err := a.Getrancher("/clusters")
	if err != nil {
		beego.Error("访问失败:", err)
	}
	result := gjson.Get(response, "data")
	result.ForEach(func(key, value gjson.Result) bool {
		name := value.Get("name").String()
		id := value.Get("id").String()
		// 获取集群下namespace
		response, err := a.Getrancher("/clusters/" + id + "/projects")
		if err != nil {
			beego.Error("访问失败:", err)
		}
		result1 := gjson.Get(response, "data")

		var project []map[string]string
		result1.ForEach(func(key, value gjson.Result) bool {
			name := value.Get("name").String()
			id := value.Get("id").String()
			project = append(project, map[string]string{
				"id":   id,
				"name": name,
			})
			return true
		})

		capacity := value.Get("capacity").String()
		requested := value.Get("requested").String()
		nodeCount := value.Get("nodeCount").String()
		entities = append(entities, map[string]interface{}{
			"id":        id,
			"name":      name,
			"requested": requested,
			"capacity":  capacity,
			"nodeCount": nodeCount,
			"project":   project,
		})
		return true
	})
	var JsonOutput = entities
	beego.Info(JsonOutput)
	u.Data["json"] = JsonOutput
	u.ServeJSON()
}

func (u RancherController) Getproject() {
	a := defaulturl{}
	// var entities []map[string]string
	groupname := make(map[string][]map[string]interface{})
	projectid := u.GetString("projectid")
	beego.Info(projectid)
	response, err := a.Getrancher("/projects/" + projectid + "/workloads")
	if err != nil {
		beego.Error("访问失败:", err)
	}
	result := gjson.Get(response, "data")
	result.ForEach(func(key, value gjson.Result) bool {
		name := value.Get("name").String()
		id := value.Get("id").String()
		namespaceId := value.Get("namespaceId").String()
		// containers := value.Get("containers").String()
		groupname[namespaceId] = append(groupname[namespaceId], map[string]interface{}{
			"id":          id,
			"name":        name,
			"namespaceId": namespaceId,
			// "containers":  containers,
		})
		return true
	})

	var JsonOutput = groupname
	beego.Info(JsonOutput)
	u.Data["json"] = JsonOutput
	u.ServeJSON()
}

func (u RancherController) Getworker() {
	a := defaulturl{}
	// groupname := make(map[string][]map[string]string)
	var JsonOutput workretrun
	var podgroup []map[string]interface{}
	var tcpport1 []gjson.Result
	var envresult map[string]gjson.Result

	projectid := u.GetString("projectid")
	workerid := u.GetString("workerid")

	// 获取workload信息
	response, err := a.Getrancher("/projects/" + projectid + "/workloads/" + workerid)
	if err != nil {
		beego.Error("访问失败:", err)
	}

	JsonOutput.Name = gjson.Get(response, "name").String()
	JsonOutput.Wtype = gjson.Get(response, "type").String()
	JsonOutput.Hostname = gjson.Get(response, "publicEndpoints.#.hostname").String()
	JsonOutput.Port = gjson.Get(response, "publicEndpoints.#.port").String()
	JsonOutput.Image = gjson.Get(response, "containers.#.image").String()
	environment := gjson.Get(response, "containers.#.environment").Array()
	tcpport := gjson.Get(response, "containers.#.ports").Array()
	livenessProbe := gjson.Get(response, "containers.#.livenessProbe").Array()
	readinessProbe := gjson.Get(response, "containers.#.readinessProbe").Array()
	for _, value := range environment {
		envresult = value.Map()
	}
	for item, value := range envresult {
		JsonOutput.Environment = append(JsonOutput.Environment, map[string]interface{}{
			"key":   item,
			"value": value.String(),
		})
	}
	for _, value := range tcpport {
		tcpport1 = value.Array()
	}

	for _, value := range tcpport1 {
		JsonOutput.Tcpport = append(JsonOutput.Tcpport, value.Map())
	}

	for _, value := range livenessProbe {
		JsonOutput.LivenessProbe = value.Map()
	}
	for _, value := range readinessProbe {
		JsonOutput.ReadinessProbe = value.Map()
	}
	// beego.Info(environment)
	// 筛选pod信息
	response, err = a.Getrancher("/projects/" + projectid + "/pod/")
	if err != nil {
		beego.Error("访问失败:", err)
	}
	result := gjson.Get(response, "data")
	result.ForEach(func(key, value gjson.Result) bool {
		if workerid == value.Get("workloadId").String() {
			// 获取环境变量以及image信息
			podgroup = append(podgroup, map[string]interface{}{
				"image":  value.Get("containers.#.image").String(),
				"name":   value.Get("name").String(),
				"state":  value.Get("state").String(),
				"nodeip": value.Get("status.nodeIp").String(),
			})
		}
		return true
	})
	// 转化map
	// beego.Info(JsonOutput)
	// m, ok := gjson.Parse(name).Value().(map[string]interface{})
	// if !ok {
	// 	beego.Error("转化map失败")
	// }
	m := make(map[string]interface{})
	JsonOutput.Podgroup = podgroup
	beego.Info(JsonOutput)
	j, error := json.Marshal(JsonOutput)
	if error != nil {
		beego.Error(error)
	}
	json.Unmarshal(j, &m)
	beego.Info(m)
	u.Data["json"] = m
	u.ServeJSON()
}

func (u RancherController) Changeworker() {
	a := defaulturl{}
	// var newcontainers []map[string]interface{}
	projectid := u.GetString("projectid")
	workerid := u.GetString("workerid")
	image := u.GetString("image")
	name := u.GetString("name") // 多容器使用需要根据不通的name来修改
	response, err := a.Getrancher("/projects/" + projectid + "/workloads/" + workerid)
	if err != nil {
		beego.Error("访问失败:", err)
	}
	m, ok := gjson.Parse(response).Value().(map[string]interface{})
	if !ok {
		beego.Error("转化map失败")
	}
	containers := m["containers"]
	b := containers.([]interface{})
	for _, n := range b {
		if name != "" {
			a := n.(map[string]interface{})
			if a["name"] == name {
				a["image"] = image
			}
		} else {
			beego.Info("进入此方法")
			a := n.(map[string]interface{})
			a["image"] = image
			beego.Info(a)
		}
	}
	m["containers"] = b
	_, err = a.Updaterancher("/projects/"+projectid+"/workloads/"+workerid, m)
	if err != nil {
		beego.Error("更新失败:", err)
	}
	u.Data["json"] = 1
	u.ServeJSON()
}
