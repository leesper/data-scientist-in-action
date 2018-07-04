package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	header = `<?xml version="1.0" encoding="UTF-8"?>
<?xml-stylesheet type="text/xsl" href="configuration.xsl"?>
<!--
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License. See accompanying LICENSE file.
-->

<!-- Put site-specific property overrides in this file. -->

`
)

// Conf represents Hadoop configuration
type Conf struct {
	XMLName    xml.Name `xml:"configuration"`
	Properties []Property
}

// Property represents individual property in configuration
type Property struct {
	XMLName xml.Name `xml:"property"`
	Name    string   `xml:"name"`
	Value   string   `xml:"value"`
}

func loadEnvs(envPath string) map[string]map[string]string {
	f, err := os.Open(envPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	data, _ := ioutil.ReadAll(f)
	var envs map[string]map[string]string
	err = json.Unmarshal(data, &envs)
	if err != nil {
		log.Fatalln(err)
	}
	return envs
}

func configure(confPath, module string, envs map[string]map[string]string) {
	fmt.Printf("configuring %s\n", module)

	f, err := os.Open(confPath)
	if err != nil {
		log.Fatalln(err)
	}
	data, _ := ioutil.ReadAll(f)
	f.Close()

	var conf Conf
	err = xml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalln(err)
	}

	for name, value := range envs[strings.ToUpper(module)] {
		property := Property{
			Name:  name,
			Value: value,
		}
		conf.Properties = append(conf.Properties, property)
	}

	data, err = xml.MarshalIndent(&conf, "", "    ")
	if err != nil {
		log.Fatalln(err)
	}

	f, err = os.Create(confPath)
	defer f.Close()
	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(header + string(data))
	if err != nil {
		log.Fatalln(err)
	}

	writer.Flush()
}

func main() {
	envs := loadEnvs("hadoop.env.json")

	configure("/etc/hadoop/core-site.xml", "core", envs)
	configure("/etc/hadoop/hdfs-site.xml", "hdfs", envs)
	configure("/etc/hadoop/yarn-site.xml", "yarn", envs)
	configure("/etc/hadoop/mapred-site.xml", "mapred", envs)
}
