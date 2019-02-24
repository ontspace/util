/*
 * Copyright (C) 2019 The OntSpace Authors
 *
 * The program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *
 * The MIT License (MIT)
 *
 * Copyright (c) 2014-2016 Juan Batiz-Benet
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

// +build !windows

package ulimit_test

import (
	"strings"
	"syscall"
	"testing"

	"github.com/ontspace/util/ulimit"
)

func TestGetFdLimit(t *testing.T) {
	t.Log("Testing get fd count")
	soft, hard, err := ulimit.GetFdLimit()
	if err != nil {
		t.Fatalf("Get fd count: %s", err)
	}
	if soft > hard {
		t.Fatalf("get fd count: soft %d > hard %d", soft, hard)
	}
}

func TestSetFdLimit(t *testing.T) {
	t.Log("Testing file descriptor count")
	limit := uint64(20480)
	newLimit, err := ulimit.SetFdLimit(limit)
	if err != nil {
		t.Fatalf("Cannot manage file descriptors: %s", err)
	}
	if newLimit < limit {
		t.Fatalf("Maximum file descriptors default value changed: %d", newLimit)
	}
}

func TestManageInvalidNFds(t *testing.T) {
	t.Logf("Testing file descriptor invalidity")

	var err error
	rlimit := syscall.Rlimit{}
	if err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		t.Fatal("Cannot get the file descriptor count")
	}

	value := rlimit.Max + rlimit.Cur
	if newLimit, err := ulimit.SetFdLimit(value); err == nil {
		t.Errorf("SetFdLimit should return an error: %d, %s", newLimit, err)
	} else if err != nil {
		flag := strings.Contains(err.Error(),
			"is larger than the hard limit")
		if !flag {
			t.Errorf("SetFdLimit returned unexpected error: %s", err)
		}
	}
}

func TestSetFdLimitWithEnvSet(t *testing.T) {
	t.Logf("Testing file descriptor manager")

	var err error
	rlimit := syscall.Rlimit{}
	if err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		t.Fatal("Cannot get the file descriptor count")
	}

	value := rlimit.Max - rlimit.Cur + 1
	if _, err = ulimit.SetFdLimit(value); err != nil {
		t.Errorf("Cannot manage file descriptor count")
	}
}
