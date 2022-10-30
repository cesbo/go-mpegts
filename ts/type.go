package ts

type PacketType uint32

const (
	PACKET_UNKNOWN PacketType = iota + 0
	PACKET_PAT
	PACKET_CAT
	PACKET_PMT
	PACKET_NIT
	PACKET_BAT
	PACKET_SDT
	PACKET_EIT
	PACKET_TDT
	PACKET_TOT
	PACKET_ECM
	PACKET_EMM
	PACKET_VIDEO_H261
	PACKET_VIDEO_H262
	PACKET_VIDEO_H263
	PACKET_VIDEO_H264
	PACKET_VIDEO_H265
	PACKET_AUDIO_MP2
	PACKET_AUDIO_MP3
	PACKET_AUDIO_AAC
	PACKET_AUDIO_LATM
	PACKET_AUDIO_AC_3
	PACKET_AUDIO_EAC_3
	PACKET_SUB
	PACKET_TTX
	PACKET_AIT
	PACKET_DATA
	PACKET_NULL
)

var typeDescription = []string{
	"Unknown",
	"PAT (Program Association Table)",
	"CAT (Conditional access Table)",
	"PMT (Program Map Table)",
	"NIT (Network Information Table)",
	"BAT (Bouquet Association Table)",
	"SDT (Service Description Table)",
	"EIT (Event Information Table)",
	"TDT (Time and Date Table)",
	"TOT (Time Offset Table)",
	"ECM (Entitlement Control Message)",
	"EMM (Entitlement Management Messages)",
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
	"Subtitling",
	"Teletext",
	"AIT (Application Information Table)",
	"Private Data",
	"NULL-TS",
}

func (t PacketType) String() string {
	return typeDescription[t]
}
