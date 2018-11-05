// +build linux

/*
 * Copyright Go-IIoT (https://github.com/goiiot)
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

package libserial

import (
	"os"

	"golang.org/x/sys/unix"
)

func sysReadBaudRate(fd uintptr) uint32 {
	tty, err := unix.IoctlGetTermios(int(fd), unix.TCGETS)
	if err != nil {
		return 0
	}

	return tty.Cflag & unix.CBAUD
}

func (s *SerialPort) sysOpen(f *os.File, timeout uint8) error {
	_, err := unix.IoctlGetTermios(int(f.Fd()), unix.TCGETS)
	if err != nil {
		return err
	}

	tty := &unix.Termios{
		Cflag: unix.CREAD | unix.CLOCAL | uint32(s.controlOptions),
		Iflag: uint32(s.inputOptions),
	}

	if timeout == 0 {
		// set blocking read with at least 1 byte have read if no timeout defined
		tty.Cc[unix.VMIN] = 1
	}
	// set read timeout
	tty.Cc[unix.VTIME] = timeout

	if err = unix.IoctlSetTermios(int(f.Fd()), unix.TCSETS, tty); err != nil {
		return err
	}

	// set blocking
	// if err = unix.SetNonblock(int(f.Fd()), false); err != nil {
	// 	return err
	// }
	return nil
}