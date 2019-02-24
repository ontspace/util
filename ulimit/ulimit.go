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

package ulimit

import (
	"fmt"
	"syscall"
)

var (
	supportsFDManagement = false

	// getlimit returns the soft and hard limits of file descriptors counts
	getLimit func() (uint64, uint64, error)
	// set limit sets the soft and hard limits of file descriptors counts
	setLimit func(uint64, uint64) error
)

// GetFdLimit return the current soft and hard limits of fd counts
func GetFdLimit() (uint64, uint64, error) {
	if !supportsFDManagement {
		return 0, 0, fmt.Errorf("FD management not support")
	}

	return getLimit()
}

// SetFdLimit raise the current max file descriptor count
// of the process
func SetFdLimit(fdlimit uint64) (newLimit uint64, err error) {
	if !supportsFDManagement {
		return 0, fmt.Errorf("FD management not support")
	}

	soft, hard, err := getLimit()
	if err != nil {
		return 0, err
	}

	if fdlimit <= soft {
		return soft, nil
	}

	// the soft limit is the value that the kernel enforces for the
	// corresponding resource
	// the hard limit acts as a ceiling for the soft limit
	// an unprivileged process may only set it's soft limit to a
	// alue in the range from 0 up to the hard limit
	if err = setLimit(fdlimit, fdlimit); err != nil {
		if err != syscall.EPERM {
			return 0, fmt.Errorf("error setting: ulimit: %s", err)
		}

		// the process does not have permission so we should only
		// set the soft value
		if fdlimit > hard {
			return 0, fmt.Errorf("cannot set rlimit, %d is larger than the hard limit", fdlimit)
		}

		if err = setLimit(fdlimit, hard); err != nil {
			return 0, fmt.Errorf("error setting ulimit without hard limit: %s", err)
		}
	}

	return fdlimit, nil
}
