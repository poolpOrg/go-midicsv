package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

func main() {
	fp, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	csvReader := csv.NewReader(fp)
	csvReader.FieldsPerRecord = -1
	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	midiSMF := smf.New()
	deltaTicks := make(map[int]uint32)

	for _, record := range records {
		trackNo, err := strconv.Atoi(record[0])
		if err != nil {
			panic(err)
		}

		time64, err := strconv.ParseUint(record[1], 10, 32)
		if err != nil {
			panic(err)
		}
		abstime := uint32(time64)
		_ = trackNo

		if _, exists := deltaTicks[trackNo]; !exists {
			deltaTicks[trackNo] = 0
		}
		time := abstime - deltaTicks[trackNo]
		deltaTicks[trackNo] = abstime

		recordType := record[2]

		switch recordType {
		case "Header":
			//format, err := strconv.Atoi(record[3])
			//if err != nil {
			//	panic(err)
			//}
			nTracks, err := strconv.Atoi(record[4])
			if err != nil {
				panic(err)
			}
			division, err := strconv.Atoi(record[5])
			if err != nil {
				panic(err)
			}
			midiSMF.Tracks = make([]smf.Track, nTracks+1)
			midiSMF.TimeFormat = smf.MetricTicks(division)

		case "Start_track":

			// pass

		case "Tempo":
			tempo, err := strconv.Atoi(record[3])
			if err != nil {
				panic(err)
			}
			bpm := float64(60_000_000) / float64(tempo)
			midiSMF.Tracks[trackNo].Add(time, smf.MetaTempo(bpm))
			// pass

		case "Time_signature":
			num, err := strconv.ParseUint(record[3], 10, 8)
			if err != nil {
				panic(err)
			}
			numerator := uint8(num)

			denom, err := strconv.ParseUint(record[4], 10, 8)
			if err != nil {
				panic(err)
			}
			denominator := uint8(denom)

			click, err := strconv.ParseUint(record[5], 10, 8)
			if err != nil {
				panic(err)
			}
			clocksPerClick := uint8(click)

			notesQ, err := strconv.ParseUint(record[6], 10, 8)
			if err != nil {
				panic(err)
			}
			demiSemiQuaverPerQuarter := uint8(notesQ)

			midiSMF.Tracks[trackNo].Add(time, smf.MetaTimeSig(numerator, denominator, clocksPerClick, demiSemiQuaverPerQuarter))

		case "Title_t":
			midiSMF.Tracks[trackNo].Add(time, smf.MetaTrackSequenceName(record[3]))

		case "End_track":
			midiSMF.Tracks[trackNo].Close(time)

		case "Program_c":
			c, err := strconv.ParseUint(record[3], 10, 8)
			if err != nil {
				panic(err)
			}
			channel := uint8(c)

			p, err := strconv.ParseUint(record[4], 10, 8)
			if err != nil {
				panic(err)
			}
			program := uint8(p)

			midiSMF.Tracks[trackNo].Add(time, midi.ProgramChange(channel, program))

		case "Key_signature":
			//midiSMF.Tracks[trackNo].Add(time, smf.MetaKey(key, bool, num, ifFlat))
			//pass

		case "Control_c":
			c, err := strconv.ParseUint(record[3], 10, 8)
			if err != nil {
				panic(err)
			}
			channel := uint8(c)

			ctrl, err := strconv.ParseUint(record[4], 10, 8)
			if err != nil {
				panic(err)
			}
			controller := uint8(ctrl)

			v, err := strconv.ParseUint(record[5], 10, 8)
			if err != nil {
				panic(err)
			}
			value := uint8(v)

			midiSMF.Tracks[trackNo].Add(time, midi.ControlChange(channel, controller, value))

		case "Note_on_c":
			c, err := strconv.ParseUint(record[3], 10, 8)
			if err != nil {
				panic(err)
			}
			channel := uint8(c)

			k, err := strconv.ParseUint(record[4], 10, 8)
			if err != nil {
				panic(err)
			}
			key := uint8(k)

			v, err := strconv.ParseUint(record[5], 10, 8)
			if err != nil {
				panic(err)
			}
			velocity := uint8(v)

			midiSMF.Tracks[trackNo].Add(time, midi.NoteOn(channel, key, velocity))

		case "Note_off_c":
			c, err := strconv.ParseUint(record[3], 10, 8)
			if err != nil {
				panic(err)
			}
			channel := uint8(c)

			k, err := strconv.ParseUint(record[4], 10, 8)
			if err != nil {
				panic(err)
			}
			key := uint8(k)

			v, err := strconv.ParseUint(record[5], 10, 8)
			if err != nil {
				panic(err)
			}
			_ = v

			midiSMF.Tracks[trackNo].Add(time, midi.NoteOff(channel, key))

		case "End_of_file":
			//pass
		default:
			fmt.Println(record)
			panic(err)
		}
	}

	err = midiSMF.WriteFile("/tmp/foobar.mid")
	if err != nil {
		panic(err)
	}

}
