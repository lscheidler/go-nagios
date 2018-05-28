/*
 * Copyright 2018 Lars Eric Scheidler
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package nagios
package nagios

import (
  "fmt"
  "math"
  "os"
  "strings"
)

// Nagios holds the nagios message representation
type Nagios struct {
  ok       []string
  warning  []string
  critical []string
  unknown  []string
  perfdata []string
  ShowPerfdata bool
}

// Init initialize basis nagios struct
func Init() Nagios {
  return Nagios{
    ShowPerfdata: true,
  }
}

// Exit exits with the corresponding exit code
func (nagios *Nagios) Exit() {
  exitcode, msg := nagios.getMsg()
  fmt.Println(msg)
  os.Exit(exitcode)
}

// Critical adds new critical message
func (nagios *Nagios) Critical(msg string) {
  nagios.critical = append(nagios.critical, msg)
}

// Warning adds new warning message
func (nagios *Nagios) Warning(msg string) {
  nagios.warning = append(nagios.warning, msg)
}

// Ok adds new ok message
func (nagios *Nagios) Ok(msg string) {
  nagios.ok = append(nagios.ok, msg)
}

// Unknown adds new unknown message
func (nagios *Nagios) Unknown(msg string) {
  nagios.unknown = append(nagios.unknown, msg)
}

// Perfdata adds perdata key and value
func (nagios *Nagios) Perfdata(key string, value string) {
  nagios.perfdata = append(nagios.perfdata, key + "=" + value)
}

// CheckThreshold checks a value againts thresholds
func (nagios *Nagios) CheckThreshold(name string, value float64, warning float64, critical float64) {
  nagios.checkThreshold(name, value, warning, critical, "")
}

// CheckPercentageThreshold checks the percentage calculated from value and max against thresholds
func (nagios *Nagios) CheckPercentageThreshold(name string, value float64, max float64, warning float64, critical float64) {
  percentage := value/max*100.0
  msg := fmt.Sprintf("%s=%.2f%%", name, percentage)
  nagios.checkThreshold(name, percentage, warning, critical, msg)
}

// checkThreshold
func (nagios *Nagios) checkThreshold(name string, value float64, warning float64, critical float64, msg string) {
  if msg == "" {
    msg = fmt.Sprintf("%s=%f", name, value)
  }

  if ! math.IsNaN(critical) && value > critical {
    nagios.critical = append(nagios.critical, msg)
  } else if value > warning {
    nagios.warning = append(nagios.warning, msg)
  } else {
    nagios.ok = append(nagios.ok, msg)
  }
  nagios.perfdata = append(nagios.perfdata, fmt.Sprintf("%s=%f", name, value) )
}

// getMsg
func (nagios *Nagios) getMsg() (int, string) {
  prefix := "OK -"
  exitcode := 0

  msg := ""
  if len(nagios.critical) > 0 {
    prefix = "CRITICAL -"
    exitcode = 1
    msg += " critical(" + strings.Join(nagios.critical, ", ") + ")"
  }
  if len(nagios.warning) > 0 {
    if exitcode == 0 {
      prefix = "WARNING -"
      exitcode = 2
    }
    msg += " warning(" + strings.Join(nagios.warning, ", ") + ")"
  }
  if len(nagios.unknown) > 0 {
    if exitcode == 0 {
      prefix = "UNKNOWN -"
      exitcode = 3
    }
    msg += " unknown(" + strings.Join(nagios.unknown, ", ") + ")"
  }
  if len(nagios.ok) > 0 {
    msg += " ok(" + strings.Join(nagios.ok, ", ") + ")"
  }
  if msg == "" {
    msg += " Everything is ok"
  }
  return exitcode, prefix + msg + nagios.getPerfdataMsg()
}

// getPerfdataMsg
func (nagios *Nagios) getPerfdataMsg() string {
  msg := ""
  if nagios.ShowPerfdata && len(nagios.perfdata) > 0 {
    msg += " | " + strings.Join(nagios.perfdata, " ")
  }
  return msg
}
