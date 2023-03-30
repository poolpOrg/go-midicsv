package main

import (
	"fmt"
	"os"

	"gitlab.com/gomidi/midi/v2/smf"
)

func main() {
	fp, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	x := smf.ReadTracksFrom(fp)
	fmt.Printf("0,0,Header,%d,%d,%d\n", x.SMF().Format(), len(x.SMF().Tracks), 0)
	//0, 0, Header, format, nTracks, division

	x.Do(
		func(te smf.TrackEvent) {
			switch te.Message.Type().String() {
			case "MetaTrackName":
				var trackName string
				te.Message.GetMetaTrackName(&trackName)
				fmt.Printf("%d,%d,Title_t,%s\n", te.TrackNo, te.AbsTicks, trackName)

			case "MetaSMPTEOffset":
				var hour uint8
				var minute uint8
				var second uint8
				var frame uint8
				var fractframe uint8
				te.Message.GetMetaSMPTEOffsetMsg(&hour, &minute, &second, &frame, &fractframe)
				fmt.Printf("%d,0,SMPTE_offset,%d,%d,%d,%d,%d\n", te.TrackNo, hour, minute, second, frame, fractframe)

			case "MetaTimeSig":
				var numerator uint8
				var denominator uint8
				var clocksPerClick uint8
				var demiSemiQuaverPerQuarter uint8

				te.Message.GetMetaTimeSig(&numerator, &denominator, &clocksPerClick, &demiSemiQuaverPerQuarter)
				fmt.Printf("%d,%d,Time_signature,%d,%d,%d,%d\n", te.TrackNo, te.AbsTicks, numerator, denominator, clocksPerClick, demiSemiQuaverPerQuarter)

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
				fmt.Printf("%d,%d,Key_signature,%d,\"%s\"\n", te.TrackNo, te.AbsTicks, key, mode)

			case "MetaTempo":
				var bpm float64

				te.Message.GetMetaTempo(&bpm)
				tempo := float64(60_000_000) / bpm
				fmt.Printf("%d,%d,Tempo,%d\n", te.TrackNo, te.AbsTicks, int(tempo))

			case "MetaPort":
				var port uint8

				te.Message.GetMetaPort(&port)
				fmt.Printf("%d,%d,MIDI_port,%d\n", te.TrackNo, te.AbsTicks, port)

			case "MetaLyric":
				var text string

				te.Message.GetMetaLyric(&text)
				fmt.Printf("%d,%d,Lyric_t,%s\n", te.TrackNo, te.AbsTicks, text)

			case "MetaEndOfTrack":
				fmt.Printf("%d,%d,End_track\n", te.TrackNo, te.AbsTicks)

			case "ProgramChange":
				var channel uint8
				var program uint8

				te.Message.GetProgramChange(&channel, &program)
				fmt.Printf("%d,%d,Program_c,%d,%d\n", te.TrackNo, te.AbsTicks, channel, program)

			case "ControlChange":
				var channel uint8
				var controller uint8
				var value uint8

				te.Message.GetControlChange(&channel, &controller, &value)
				fmt.Printf("%d,%d,Control_c,%d,%d,%d\n", te.TrackNo, te.AbsTicks, channel, controller, value)

			case "NoteOn":
				var channel uint8
				var key uint8
				var velocity uint8

				te.Message.GetNoteOn(&channel, &key, &velocity)
				fmt.Printf("%d,%d,Note_on_c,%d,%d,%d\n", te.TrackNo, te.AbsTicks, channel, key, velocity)

			case "NoteOff":
				var channel uint8
				var key uint8
				var velocity uint8

				te.Message.GetNoteOff(&channel, &key, &velocity)
				fmt.Printf("%d,%d,Note_off_c,%d,%d,%d\n", te.TrackNo, te.AbsTicks, channel, key, velocity)

			default:
				fmt.Printf("[%v] @%vms %s\n", te.TrackNo, te.AbsMicroSeconds/1000, te.Message.String())
				panic("unknown type: " + te.Message.Type().String())
			}
		},
	)
	fmt.Printf("0,0,End_of_file\n")
}
