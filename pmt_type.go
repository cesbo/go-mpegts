package mpegts

// StreamType is an elementary stream type
type StreamType int

const (
	StreamData StreamType = iota
	StreamVideoH261
	StreamVideoH262
	StreamVideoH263
	StreamVideoH264
	StreamVideoH265
	StreamAudioMP2
	StreamAudioMP3
	StreamAudioAAC
	StreamAudioLATM
	StreamAudioAC3
	StreamAudioEAC3
	StreamDataSCTE35
	StreamDataSubtitles
	StreamDataTeletext
	StreamDataAIT
)

var streamTypeDescription = []string{
	"Private Data",
	"Video H.261/11172 (MPEG-1)",
	"Video H.262/13818-2 (MPEG-2)",
	"Visual H.263/14496-2 (MPEG-4 Visual)",
	"Video H.264/14496-10 (MPEG-4/AVC)",
	"Video H.265/23008-2 (HEVC)",
	"Audio MP2/11172-3 (MPEG-1 Layer 2)",
	"Audio MP3/13818-3 (MPEG-2 Layer 3)",
	"Audio AAC/13818-7 (MPEG-2 with ADTS transport syntax)",
	"Audio MPEG-4/14496-3 (MPEG-4 with LATM transport syntax)",
	"Audio AC-3 (Dolby Digital)",
	"Audio E-AC-3 (Dolby Digital Plus)",
	"SCTE-35",
	"Subtitles",
	"Teletext",
	"AIT (Application Information Table)",
}

func (t StreamType) String() string {
	return streamTypeDescription[t]
}
