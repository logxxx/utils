package ffmpeg

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type FFProbe string

type VideoFile struct {
	JSON        FFProbeJSON
	AudioStream *FFProbeStream
	VideoStream *FFProbeStream

	Path         string
	Title        string
	Comment      string
	Container    string
	Duration     float64
	StartTime    float64
	Bitrate      int64
	Size         int64
	CreationTime time.Time

	VideoCodec   string
	VideoBitrate int64
	Width        int
	Height       int
	FrameRate    float64
	Rotation     int64
	FrameCount   int64

	AudioCodec string
}

// FFProbeJSON is the JSON output of ffprobe.
type FFProbeJSON struct {
	Format struct {
		BitRate        string `json:"bit_rate"`
		Duration       string `json:"duration"`
		Filename       string `json:"filename"`
		FormatLongName string `json:"format_long_name"`
		FormatName     string `json:"format_name"`
		NbPrograms     int    `json:"nb_programs"`
		NbStreams      int    `json:"nb_streams"`
		ProbeScore     int    `json:"probe_score"`
		Size           string `json:"size"`
		StartTime      string `json:"start_time"`
		Tags           struct {
			CompatibleBrands string   `json:"compatible_brands"`
			CreationTime     JSONTime `json:"creation_time"`
			Encoder          string   `json:"encoder"`
			MajorBrand       string   `json:"major_brand"`
			MinorVersion     string   `json:"minor_version"`
			Title            string   `json:"title"`
			Comment          string   `json:"comment"`
		} `json:"tags"`
	} `json:"format"`
	Streams []FFProbeStream `json:"streams"`
	Error   struct {
		Code   int    `json:"code"`
		String string `json:"string"`
	} `json:"error"`
}

// FFProbeStream is a JSON representation of an ffmpeg stream.
type FFProbeStream struct {
	AvgFrameRate       string `json:"avg_frame_rate"`
	BitRate            string `json:"bit_rate"`
	BitsPerRawSample   string `json:"bits_per_raw_sample,omitempty"`
	ChromaLocation     string `json:"chroma_location,omitempty"`
	CodecLongName      string `json:"codec_long_name"`
	CodecName          string `json:"codec_name"`
	CodecTag           string `json:"codec_tag"`
	CodecTagString     string `json:"codec_tag_string"`
	CodecTimeBase      string `json:"codec_time_base"`
	CodecType          string `json:"codec_type"`
	CodedHeight        int    `json:"coded_height,omitempty"`
	CodedWidth         int    `json:"coded_width,omitempty"`
	DisplayAspectRatio string `json:"display_aspect_ratio,omitempty"`
	Disposition        struct {
		AttachedPic     int `json:"attached_pic"`
		CleanEffects    int `json:"clean_effects"`
		Comment         int `json:"comment"`
		Default         int `json:"default"`
		Dub             int `json:"dub"`
		Forced          int `json:"forced"`
		HearingImpaired int `json:"hearing_impaired"`
		Karaoke         int `json:"karaoke"`
		Lyrics          int `json:"lyrics"`
		Original        int `json:"original"`
		TimedThumbnails int `json:"timed_thumbnails"`
		VisualImpaired  int `json:"visual_impaired"`
	} `json:"disposition"`
	Duration          string `json:"duration"`
	DurationTs        int    `json:"duration_ts"`
	HasBFrames        int    `json:"has_b_frames,omitempty"`
	Height            int    `json:"height,omitempty"`
	Index             int    `json:"index"`
	IsAvc             string `json:"is_avc,omitempty"`
	Level             int    `json:"level,omitempty"`
	NalLengthSize     string `json:"nal_length_size,omitempty"`
	NbFrames          string `json:"nb_frames"`
	NbReadFrames      string `json:"nb_read_frames"`
	PixFmt            string `json:"pix_fmt,omitempty"`
	Profile           string `json:"profile"`
	RFrameRate        string `json:"r_frame_rate"`
	Refs              int    `json:"refs,omitempty"`
	SampleAspectRatio string `json:"sample_aspect_ratio,omitempty"`
	StartPts          int    `json:"start_pts"`
	StartTime         string `json:"start_time"`
	Tags              struct {
		CreationTime JSONTime `json:"creation_time"`
		HandlerName  string   `json:"handler_name"`
		Language     string   `json:"language"`
		Rotate       string   `json:"rotate"`
	} `json:"tags"`
	TimeBase      string `json:"time_base"`
	Width         int    `json:"width,omitempty"`
	BitsPerSample int    `json:"bits_per_sample,omitempty"`
	ChannelLayout string `json:"channel_layout,omitempty"`
	Channels      int    `json:"channels,omitempty"`
	MaxBitRate    string `json:"max_bit_rate,omitempty"`
	SampleFmt     string `json:"sample_fmt,omitempty"`
	SampleRate    string `json:"sample_rate,omitempty"`
}

// NewVideoFile runs ffprobe on the given path and returns a VideoFile.
func (f *FFProbe) NewVideoFile(videoPath string) (*VideoFile, error) {
	args := []string{"-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", "-show_error", videoPath}
	//log.Printf("args:%v", args)
	cmd := exec.Command(string(*f), args...)
	out, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf("FFProbe encountered an error with <%s>.\nError JSON:\n%s\nError: %s", videoPath, string(out), err.Error())
	}

	probeJSON := &FFProbeJSON{}
	if err := json.Unmarshal(out, probeJSON); err != nil {
		return nil, fmt.Errorf("error unmarshalling video data for <%s>: %s", videoPath, err.Error())
	}

	return parse(videoPath, probeJSON)
}

func parse(filePath string, probeJSON *FFProbeJSON) (*VideoFile, error) {
	if probeJSON == nil {
		return nil, fmt.Errorf("failed to get ffprobe json for <%s>", filePath)
	}

	result := &VideoFile{}
	result.JSON = *probeJSON

	if result.JSON.Error.Code != 0 {
		return nil, fmt.Errorf("ffprobe error code %d: %s", result.JSON.Error.Code, result.JSON.Error.String)
	}

	result.Path = filePath
	result.Title = probeJSON.Format.Tags.Title

	result.Comment = probeJSON.Format.Tags.Comment
	result.Bitrate, _ = strconv.ParseInt(probeJSON.Format.BitRate, 10, 64)

	result.Container = probeJSON.Format.FormatName
	duration, _ := strconv.ParseFloat(probeJSON.Format.Duration, 64)
	result.Duration = math.Round(duration*100) / 100
	fileStat, err := os.Stat(filePath)
	if err != nil {
		statErr := fmt.Errorf("error statting file <%s>: %w", filePath, err)
		log.Printf("%v", statErr)
		return nil, statErr
	}
	result.Size = fileStat.Size()
	result.StartTime, _ = strconv.ParseFloat(probeJSON.Format.StartTime, 64)
	result.CreationTime = probeJSON.Format.Tags.CreationTime.Time

	audioStream := result.getAudioStream()
	if audioStream != nil {
		result.AudioCodec = audioStream.CodecName
		result.AudioStream = audioStream
	}

	videoStream := result.getVideoStream()
	if videoStream != nil {
		result.VideoStream = videoStream
		result.VideoCodec = videoStream.CodecName
		result.FrameCount, _ = strconv.ParseInt(videoStream.NbFrames, 10, 64)
		if videoStream.NbReadFrames != "" { // if ffprobe counted the frames use that instead
			fc, _ := strconv.ParseInt(videoStream.NbReadFrames, 10, 64)
			if fc > 0 {
				result.FrameCount, _ = strconv.ParseInt(videoStream.NbReadFrames, 10, 64)
			} else {
				log.Printf("[ffprobe] <%s> invalid Read Frames count", videoStream.NbReadFrames)
			}
		}
		result.VideoBitrate, _ = strconv.ParseInt(videoStream.BitRate, 10, 64)
		var framerate float64
		if strings.Contains(videoStream.AvgFrameRate, "/") {
			frameRateSplit := strings.Split(videoStream.AvgFrameRate, "/")
			numerator, _ := strconv.ParseFloat(frameRateSplit[0], 64)
			denominator, _ := strconv.ParseFloat(frameRateSplit[1], 64)
			framerate = numerator / denominator
		} else {
			framerate, _ = strconv.ParseFloat(videoStream.AvgFrameRate, 64)
		}
		if math.IsNaN(framerate) {
			framerate = 0
		}
		result.FrameRate = math.Round(framerate*100) / 100
		if rotate, err := strconv.ParseInt(videoStream.Tags.Rotate, 10, 64); err == nil && rotate != 180 {
			result.Width = videoStream.Height
			result.Height = videoStream.Width
		} else {
			result.Width = videoStream.Width
			result.Height = videoStream.Height
		}
	}

	return result, nil
}

func (v *VideoFile) getAudioStream() *FFProbeStream {
	index := v.getStreamIndex("audio", v.JSON)
	if index != -1 {
		return &v.JSON.Streams[index]
	}
	return nil
}

func (v *VideoFile) getVideoStream() *FFProbeStream {
	index := v.getStreamIndex("video", v.JSON)
	if index != -1 {
		return &v.JSON.Streams[index]
	}
	return nil
}

func (v *VideoFile) getStreamIndex(fileType string, probeJSON FFProbeJSON) int {
	ret := -1
	for i, stream := range probeJSON.Streams {
		// skip cover art/thumbnails
		if stream.CodecType == fileType && stream.Disposition.AttachedPic == 0 {
			// prefer default stream
			if stream.Disposition.Default == 1 {
				return i
			}

			// backwards compatible behaviour - fallback to first matching stream
			if ret == -1 {
				ret = i
			}
		}
	}

	return ret
}
