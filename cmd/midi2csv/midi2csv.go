package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"gitlab.com/gomidi/midi/v2/smf"
)

func main() {
	fp, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	x := smf.ReadTracksFrom(fp)
	nTracks := len(x.SMF().Tracks)

	events := make(map[int][]string)
	c := make(chan string)
	go func() {
		for msg := range c {
			fields := strings.Split(msg, ",")
			trackNo, err := strconv.Atoi(fields[0])
			if err != nil {
				continue
			}
			events[trackNo] = append(events[trackNo], msg)
		}
	}()

	x.Do(
		func(te smf.TrackEvent) {
			switch te.Message.Type().String() {

			// File Meta-Events
			case "MetaTrackName":
				var trackName string
				te.Message.GetMetaTrackName(&trackName)
				c <- fmt.Sprintf("%d,%d,Title_t,\"%s\"", te.TrackNo, te.AbsTicks, trackName)
				// copyright
				// instrument
				// marker
				// cue

			case "MetaLyric":
				var text string

				te.Message.GetMetaLyric(&text)
				c <- fmt.Sprintf("%d,%d,Lyric_t,%s", te.TrackNo, te.AbsTicks, text)
				// text
				//sequencenumber

			case "MetaPort":
				var port uint8

				te.Message.GetMetaPort(&port)
				c <- fmt.Sprintf("%d,%d,MIDI_port,%d", te.TrackNo, te.AbsTicks, port)
				// channel prefix

			case "MetaTimeSig":
				var numerator uint8
				var denominator uint8
				var clocksPerClick uint8
				var demiSemiQuaverPerQuarter uint8

				te.Message.GetMetaTimeSig(&numerator, &denominator, &clocksPerClick, &demiSemiQuaverPerQuarter)
				c <- fmt.Sprintf("%d,%d,Time_signature,%d,%d,%d,%d", te.TrackNo, te.AbsTicks, numerator, denominator, clocksPerClick, demiSemiQuaverPerQuarter)

			case "MetaKeySig":
				var key uint8
				var num uint8
				var isMajor bool
				var isFlat bool
				var mode string

				te.Message.GetMetaKeySig(&key, &num, &isMajor, &isFlat)
				if isMajor {
					mode = "major"
				} else {
					mode = "minor"
				}
				c <- fmt.Sprintf("%d,%d,Key_signature,%d,\"%s\"", te.TrackNo, te.AbsTicks, key, mode)

			case "MetaTempo":
				var bpm float64

				te.Message.GetMetaTempo(&bpm)
				tempo := float64(60_000_000) / bpm
				c <- fmt.Sprintf("%d,%d,Tempo,%d", te.TrackNo, te.AbsTicks, int(tempo))

			case "MetaSMPTEOffset":
				var hour uint8
				var minute uint8
				var second uint8
				var frame uint8
				var fractframe uint8

				te.Message.GetMetaSMPTEOffsetMsg(&hour, &minute, &second, &frame, &fractframe)
				c <- fmt.Sprintf("%d,0,SMPTE_offset,%d,%d,%d,%d,%d", te.TrackNo, hour, minute, second, frame, fractframe)

				// sequencer specific
				// unknown meta event

			case "MetaEndOfTrack":
				c <- fmt.Sprintf("%d,%d,End_track", te.TrackNo, te.AbsTicks)

			// Channel Events
			case "NoteOn":
				var channel uint8
				var key uint8
				var velocity uint8

				te.Message.GetNoteOn(&channel, &key, &velocity)
				c <- fmt.Sprintf("%d,%d,Note_on_c,%d,%d,%d", te.TrackNo, te.AbsTicks, channel, key, velocity)

			case "NoteOff":
				var channel uint8
				var key uint8
				var velocity uint8

				te.Message.GetNoteOff(&channel, &key, &velocity)
				c <- fmt.Sprintf("%d,%d,Note_off_c,%d,%d,%d", te.TrackNo, te.AbsTicks, channel, key, velocity)

			case "PitchBend":
				var channel uint8
				var relative int16
				var absolute uint16

				te.Message.GetPitchBend(&channel, &relative, &absolute)
				c <- fmt.Sprintf("%d,%d,Control_c,%d,%d", te.TrackNo, te.AbsTicks, channel, relative)

			case "ControlChange":
				var channel uint8
				var controller uint8
				var value uint8

				te.Message.GetControlChange(&channel, &controller, &value)
				c <- fmt.Sprintf("%d,%d,Control_c,%d,%d,%d", te.TrackNo, te.AbsTicks, channel, controller, value)

			case "ProgramChange":
				var channel uint8
				var program uint8

				te.Message.GetProgramChange(&channel, &program)
				c <- fmt.Sprintf("%d,%d,Program_c,%d,%d", te.TrackNo, te.AbsTicks, channel, program)
			//case "ChannelAftertouch"
			//case "PolyAftertouch"

			// System Exclusive Events
			// System_exclusive (length, data, ...)
			// System_exclusive_packet (length, data, ...)

			default:
				fmt.Printf("[%v] @%vms %s", te.TrackNo, te.AbsMicroSeconds/1000, te.Message.String())
				panic("unknown type: " + te.Message.Type().String())
			}
		},
	)

	trackOrder := make([]int, 0)
	for trackNo := range events {
		trackOrder = append(trackOrder, trackNo)
	}
	sort.Slice(trackOrder, func(i, j int) bool { return trackOrder[i] < trackOrder[j] })

	fmt.Printf("0,0,Header,%d,%d,%d\n", x.SMF().Format(), nTracks, x.SMF().TimeFormat)
	for _, trackNo := range trackOrder {
		fmt.Printf("%d,0,Start_track\n", trackNo)
		for _, event := range events[trackNo] {
			fmt.Println(event)
		}
	}
	fmt.Printf("0,0,End_of_file\n")
}
