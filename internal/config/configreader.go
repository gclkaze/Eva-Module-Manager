// Copyright (c) 2025 Michail Dorgiakis - gclkaze@gmail.com - https://github.com/gclkaze
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package config

import (
	"fmt"
	"os"

	"github.com/magiconair/properties"
)

type ConfigReader struct {
	properties *properties.Properties
	onError    bool
	error      error
}

var TheConfigReader *ConfigReader = nil

func Init() {
	if TheConfigReader != nil {
		return
	}
	TheConfigReader = &ConfigReader{}
	TheConfigReader.error = TheConfigReader.readProperties()
	if TheConfigReader.error != nil {
		TheConfigReader.onError = true
	}
}

func InitWithPropertiesPath(prop string) {
	if TheConfigReader != nil {
		return
	}
	TheConfigReader = &ConfigReader{}
	TheConfigReader.error = TheConfigReader.readPropertiesFromPath(prop)
	if TheConfigReader.error != nil {
		TheConfigReader.onError = true
	}
}

func InitWithPropertiesMap(m *map[string]string) {
	if TheConfigReader != nil {
		return
	}
	TheConfigReader = &ConfigReader{}
	TheConfigReader.error = TheConfigReader.readPropertiesFromMap(m)
	if TheConfigReader.error != nil {
		TheConfigReader.onError = true
	}
}

func (c *ConfigReader) readPropertiesFromPath(prop string) error {
	c.properties = properties.MustLoadFile(prop, properties.UTF8)
	if c.properties == nil {
		return fmt.Errorf("couldn't read application properties file %s", prop)
	}
	return nil
}

func (c *ConfigReader) readPropertiesFromMap(m *map[string]string) error {
	c.properties = properties.LoadMap(*m)
	if c.properties == nil {
		return fmt.Errorf("couldn't read application properties map")
	}
	return nil
}

func (c *ConfigReader) readProperties() error {
	currentPath, _ := os.Getwd()
	f := currentPath + "\\internal\\config\\application.properties"
	c.properties = properties.MustLoadFile(f, properties.UTF8)
	if c.properties == nil {
		return fmt.Errorf("couldn't read application properties file %s", f)
	}
	return nil
}

func (c ConfigReader) IsOnError() bool {
	return c.onError
}

func (c *ConfigReader) GetProperties() *properties.Properties {
	return c.properties
}
