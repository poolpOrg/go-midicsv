package encoding

import (
	"fmt"
	"io"

	"gitlab.com/gomidi/midi/v2/smf"
)

type Encoder struct {
	rd *smf.TracksReader
}

func NewEncoder(rd io.Reader) *Encoder {
	return &Encoder{
		rd: smf.ReadTracksFrom(rd),
	}
}

func (e *Encoder) Encode() ([]string, error) {
	ret := make([]string, 0)
	ret = append(ret, fmt.Sprintf("0,0,Header,%d,%d,%d", e.rd.SMF().Format(), len(e.rd.SMF().Tracks), 0))
	e.rd.Do(
		func(te smf.TrackEvent) {
			switch te.Message.Type().String() {
			case "MetaTrackName":
				var trackName string
				te.Message.GetMetaTrackName(&trackName)
				ret = append(ret, fmt.Sprintf("%d,%d,Title_t,%s", te.TrackNo, te.AbsTicks, trackName))

			case "MetaSMPTEOffset":
				var hour uint8
				var minute uint8
				var second uint8
				var frame uint8
				var fractframe uint8
				te.Message.GetMetaSMPTEOffsetMsg(&hour, &minute, &second, &frame, &fractframe)
				ret = append(ret, fmt.Sprintf("%d,0,SMPTE_offset,%d,%d,%d,%d,%d", te.TrackNo, hour, minute, second, frame, fractframe))

			case "MetaTimeSig":
				var numerator uint8
				var denominator uint8
				var clocksPerClick uint8
				var demiSemiQuaverPerQuarter uint8

				te.Message.GetMetaTimeSig(&numerator, &denominator, &clocksPerClick, &demiSemiQuaverPerQuarter)
				ret = append(ret, fmt.Sprintf("%d,%d,Time_signature,%d,%d,%d,%d", te.TrackNo, te.AbsTicks, numerator, denominator, clocksPerClick, demiSemiQuaverPerQuarter))

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
				ret = append(ret, fmt.Sprintf("%d,%d,Key_signature,%d,\"%s\"", te.TrackNo, te.AbsTicks, key, mode))

			case "MetaTempo":
				var bpm float64

				te.Message.GetMetaTempo(&bpm)
				tempo := float64(60_000_000) / bpm
				ret = append(ret, fmt.Sprintf("%d,%d,Tempo,%d", te.TrackNo, te.AbsTicks, int(tempo)))

			case "MetaPort":
				var port uint8

				te.Message.GetMetaPort(&port)
				ret = append(ret, fmt.Sprintf("%d,%d,MIDI_port,%d", te.TrackNo, te.AbsTicks, port))

			case "MetaLyric":
				var text string

				te.Message.GetMetaLyric(&text)
				ret = append(ret, fmt.Sprintf("%d,%d,Lyric_t,%s", te.TrackNo, te.AbsTicks, text))

			case "MetaEndOfTrack":
				ret = append(ret, fmt.Sprintf("%d,%d,End_track", te.TrackNo, te.AbsTicks))

			case "ProgramChange":
				var channel uint8
				var program uint8

				te.Message.GetProgramChange(&channel, &program)
				ret = append(ret, fmt.Sprintf("%d,%d,Program_c,%d,%d", te.TrackNo, te.AbsTicks, channel, program))

			case "ControlChange":
				var channel uint8
				var controller uint8
				var value uint8

				te.Message.GetControlChange(&channel, &controller, &value)
				ret = append(ret, fmt.Sprintf("%d,%d,Control_c,%d,%d,%d", te.TrackNo, te.AbsTicks, channel, controller, value))

			case "NoteOn":
				var channel uint8
				var key uint8
				var velocity uint8

				te.Message.GetNoteOn(&channel, &key, &velocity)
				ret = append(ret, fmt.Sprintf("%d,%d,Note_on_c,%d,%d,%d", te.TrackNo, te.AbsTicks, channel, key, velocity))

			case "NoteOff":
				var channel uint8
				var key uint8
				var velocity uint8

				te.Message.GetNoteOff(&channel, &key, &velocity)
				ret = append(ret, fmt.Sprintf("%d,%d,Note_off_c,%d,%d,%d", te.TrackNo, te.AbsTicks, channel, key, velocity))

			default:
				//return nil, fmt.Errorf("unknown type: " + te.Message.Type().String())
			}
		},
	)

	ret = append(ret, "0,0,End_of_file")
	return ret, nil
}
