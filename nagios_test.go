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
	"math"
	"testing"
)

func TestDefault(t *testing.T) {
	nagios := Init()

	testMsg(t, &nagios, 0, "OK - Everything is ok")
}

func TestOk(t *testing.T) {
	nagios := Init()
	nagios.Ok("message")

	testMsg(t, &nagios, 0, "OK - ok(message)")
}

func TestCritical(t *testing.T) {
	nagios := Init()
	nagios.Critical("message")

	testMsg(t, &nagios, 2, "CRITICAL - critical(message)")
}

func TestWarning(t *testing.T) {
	nagios := Init()
	nagios.Warning("message")

	testMsg(t, &nagios, 1, "WARNING - warning(message)")
}

func TestUnknown(t *testing.T) {
	nagios := Init()
	nagios.Unknown("message")

	testMsg(t, &nagios, 3, "UNKNOWN - unknown(message)")
}

func TestPerfdata(t *testing.T) {
	nagios := Init()
	nagios.Ok("message")
	nagios.Perfdata("value", "555")

	testMsg(t, &nagios, 0, "OK - ok(message) | value=555")
}

func TestCheckThresholdCritical(t *testing.T) {
	nagios := Init()
	nagios.CheckThreshold("test", 5, 4, 4)

	testMsg(t, &nagios, 2, "CRITICAL - critical(test=5.000000) | test=5.000000")
}

func TestCheckThresholdWarning(t *testing.T) {
	nagios := Init()
	nagios.CheckThreshold("test", 5, 4, math.NaN())

	testMsg(t, &nagios, 1, "WARNING - warning(test=5.000000) | test=5.000000")
}

func TestCheckThresholdOk(t *testing.T) {
	nagios := Init()
	nagios.CheckThreshold("test", 5, 5, 5)

	testMsg(t, &nagios, 0, "OK - ok(test=5.000000) | test=5.000000")
}

func TestCheckThresholdMultiple(t *testing.T) {
	nagios := Init()
	nagios.CheckThreshold("test", 5, 4, 4)
	nagios.CheckThreshold("test2", 5, 4, math.NaN())

	testMsg(t, &nagios, 2, "CRITICAL - critical(test=5.000000) warning(test2=5.000000) | test=5.000000 test2=5.000000")

	nagios.CheckThreshold("test3", 5, 8, math.NaN())
	nagios.CheckThreshold("test4", 5, 8, 4)

	testMsg(t, &nagios, 2, "CRITICAL - critical(test=5.000000, test4=5.000000) warning(test2=5.000000) ok(test3=5.000000) | test=5.000000 test2=5.000000 test3=5.000000 test4=5.000000")
}

func TestCheckThresholdWithoutPerfdata(t *testing.T) {
	nagios := Init()
	nagios.ShowPerfdata = false
	nagios.CheckThreshold("test", 5, 4, 4)

	testMsg(t, &nagios, 2, "CRITICAL - critical(test=5.000000)")
}

func TestCheckPercentageThresholdCritical(t *testing.T) {
	nagios := Init()
	nagios.CheckPercentageThreshold("test", 5, 100, 4, 4)

	testMsg(t, &nagios, 2, "CRITICAL - critical(test=5.00%) | test=5.000000")
}

func TestCheckPercentageThresholdWarning(t *testing.T) {
	nagios := Init()
	nagios.CheckPercentageThreshold("test", 5, 100, 4, math.NaN())

	testMsg(t, &nagios, 1, "WARNING - warning(test=5.00%) | test=5.000000")
}

func TestCheckPercentageThresholdOk(t *testing.T) {
	nagios := Init()
	nagios.CheckPercentageThreshold("test", 5, 100, 5, 5)

	testMsg(t, &nagios, 0, "OK - ok(test=5.00%) | test=5.000000")
}

func TestCheckPercentageThresholdMultiple(t *testing.T) {
	nagios := Init()
	nagios.CheckPercentageThreshold("test", 5, 100, 4, 4)
	nagios.CheckPercentageThreshold("test2", 5, 100, 4, math.NaN())

	testMsg(t, &nagios, 2, "CRITICAL - critical(test=5.00%) warning(test2=5.00%) | test=5.000000 test2=5.000000")

	nagios.CheckPercentageThreshold("test3", 5, 100, 8, math.NaN())
	nagios.CheckPercentageThreshold("test4", 5, 100, 8, 4)

	testMsg(t, &nagios, 2, "CRITICAL - critical(test=5.00%, test4=5.00%) warning(test2=5.00%) ok(test3=5.00%) | test=5.000000 test2=5.000000 test3=5.000000 test4=5.000000")
}

func TestCheckPercentageThresholdWithoutPerfdata(t *testing.T) {
	nagios := Init()
	nagios.ShowPerfdata = false
	nagios.CheckPercentageThreshold("test", 5, 100, 4, 4)

	testMsg(t, &nagios, 2, "CRITICAL - critical(test=5.00%)")
}

func testMsg(t *testing.T, nagios *Nagios, expectedExitCode int, expectedMsg string) {
	exitcode, msg := nagios.getMsg()
	if exitcode != expectedExitCode {
		t.Errorf(
			"Exit code is not %d: %d", expectedExitCode, exitcode,
		)
	}

	if msg != expectedMsg {
		t.Errorf(
			"Output should output \"%s\" but returned \"%s\"", expectedMsg, msg,
		)
	}
}
