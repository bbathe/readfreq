package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/tarm/serial"
)

var (
	errBand = fmt.Errorf("out of band")
)

// readCIVMessageFromPort reads bytes from port and returns CIV message
func readCIVMessageFromPort(p *serial.Port) ([]byte, error) {
	var buf bytes.Buffer
	b := []byte{0}

	for {
		n, err := p.Read(b)
		if err != nil {
			log.Printf("%+v", err)
			return []byte{}, err
		}

		if n > 0 {
			// accumulate message bytes
			buf.Write(b)

			// message terminator?
			if b[0] == 0xFD {
				// return CIV message
				return buf.Bytes(), nil
			}
		} else {
			// no data available
			return []byte{}, nil
		}
	}
}

// return amateur band corresponding to frequency
func bandFromFrequency(freq int) (int, error) {
	switch {
	case freq >= 1810000 && freq <= 2000000:
		return 160, nil
	case freq >= 3500000 && freq <= 3800000:
		return 80, nil
	case freq >= 5250000 && freq <= 5450000:
		return 60, nil
	case freq >= 7000000 && freq <= 7200000:
		return 40, nil
	case freq >= 10100000 && freq <= 10150000:
		return 30, nil
	case freq >= 14000000 && freq <= 14350000:
		return 20, nil
	case freq >= 18068000 && freq <= 18168000:
		return 17, nil
	case freq >= 21000000 && freq <= 21450000:
		return 15, nil
	case freq >= 24890000 && freq <= 24990000:
		return 12, nil
	case freq >= 28000000 && freq <= 29700000:
		return 10, nil
	case freq >= 50000000 && freq <= 52000000:
		return 6, nil
	}

	return 0, errBand
}

// return band codes as used by the KPA500
func bcdFromBand(band int) ([4]int, error) {
	switch band {
	case 160:
		return [4]int{0, 0, 0, 1}, nil
	case 80:
		return [4]int{0, 0, 1, 0}, nil
	case 60:
		return [4]int{0, 0, 0, 0}, nil
	case 40:
		return [4]int{0, 0, 1, 1}, nil
	case 30:
		return [4]int{0, 1, 0, 0}, nil
	case 20:
		return [4]int{0, 1, 0, 1}, nil
	case 17:
		return [4]int{0, 1, 1, 0}, nil
	case 15:
		return [4]int{0, 1, 1, 1}, nil
	case 12:
		return [4]int{1, 0, 0, 0}, nil
	case 10:
		return [4]int{1, 0, 0, 1}, nil
	case 6:
		return [4]int{1, 0, 1, 0}, nil
	}

	return [4]int{1, 1, 1, 1}, errBand
}

func main() {
	c := &serial.Config{
		Name: "COM9",
		Baud: 9600,
	}

	p, err := serial.OpenPort(c)
	if err != nil {
		log.Printf("%+v", err)
		return
	}
	defer p.Close()

	for {
		// read ci-v message
		r, err := readCIVMessageFromPort(p)
		if err != nil {
			log.Fatalf("%+v", err)
		}

		// is it transfer operating frequency data?
		if len(r) == 11 && r[2] == 0x00 && r[3] == 0x94 {
			fmt.Printf("%X ", r)

			// radio sends as least significant byte first, flip order of bytes
			fd := fmt.Sprintf("%02X%02X%02X%02X%02X", r[9], r[8], r[7], r[6], r[5])
			fmt.Print(fd, " ")

			// convert to number
			freq, err := strconv.Atoi(fd)
			if err != nil {
				log.Printf("%+v", err)
				return
			}
			fmt.Print(freq, " ")

			// get corresponding band
			band, err := bandFromFrequency(freq)
			if err == nil {
				fmt.Print(band, " ")

				// get bcd bits that would be sent to kpa500
				bcd, err := bcdFromBand(band)
				if err == nil {
					fmt.Print(bcd)
				}
			}

			// next
			fmt.Println("")
		}
	}
}
