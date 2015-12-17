package video

import "io"
import "github.com/32bitkid/mpeg/util"
import "errors"

var log = util.NewLog("mpeg:video:frame_provider")

var ErrUnsupportedVideoStream_ISO_IEC_11172_2 = errors.New("unsupported video stream ISO/IEC 11172-2")

type FrameProvider interface {
	Next() (interface{}, error)
}

func NewFrameProvider(source io.Reader) FrameProvider {
	return &frameProvider{
		br: util.NewBitReader(source),
	}
}

type frameProvider struct {
	br util.BitReader32
}

func (fp *frameProvider) Next() (interface{}, error) {
	br := fp.br

	// Align to next start code
	err := next_start_code(br)
	if err != nil {
		panic(err)
	}

	// Read sequence_header
	sqh, err := sequence_header(br)
	if err != nil {
		panic(err)
	}

	// peek for sequence_extension
	val, err := br.Peek32(32)
	if err != nil {
		panic(err)
	}

	if val == ExtensionStartCode {

		se, err := sequence_extension(br)

		log.Printf("%#v\n", se)

		for {
			extension_and_user_data(0, br)

			for {
				nextbits, err := br.Peek32(32)
				if err != nil {
					panic(err)
				}

				if StartCode(nextbits) == GroupStartCode {
					group_of_pictures_header(br)
					extension_and_user_data(1, br)
				}
				picture_header(br)
			}

			val, err := br.Peek32(32)
			log.Printf("%x\n", val)
			if err != nil {
				panic(err)
			}

			if val == SequenceEndStartCode {
				break
			}
		}

		err = br.Trash(32)

		return sqh, err
	} else {
		// Stream is MPEG-1 Video
		return nil, ErrUnsupportedVideoStream_ISO_IEC_11172_2
	}

}