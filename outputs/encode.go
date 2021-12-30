package outputs

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/rakshasa/bencode-go"
	"github.com/rakshasa/rbedit/objects"
	"github.com/rakshasa/rbedit/types"
)

func NewEncodeGenericBencode() types.EncodeFunc {
	return func(metadata types.IOMetadata, object interface{}) (types.IOMetadata, []byte, error) {
		var buf bytes.Buffer

		if err := bencode.Marshal(&buf, object); err != nil {
			return types.IOMetadata{}, nil, fmt.Errorf("failed to encode data: %v", err)
		}

		return metadata, buf.Bytes(), nil
	}
}

func NewEncodeTorrentBencode() types.EncodeFunc {
	return func(metadata types.IOMetadata, object interface{}) (types.IOMetadata, []byte, error) {
		torrentInfo, err := objects.NewTorrentInfo(object)
		if err != nil {
			return types.IOMetadata{}, nil, fmt.Errorf("failed to encode data to bencoded torrent, not a valid torrent: %v", err)
		}
		metadata.OutputTorrentInfo = &torrentInfo

		var buf bytes.Buffer

		if err := bencode.Marshal(&buf, object); err != nil {
			return types.IOMetadata{}, nil, fmt.Errorf("failed to encode data to bencoded torrent: %v", err)
		}

		return metadata, buf.Bytes(), nil
	}
}

func NewEncodePrint() types.EncodeFunc {
	return func(metadata types.IOMetadata, object interface{}) (types.IOMetadata, []byte, error) {
		return metadata, []byte(objects.SprintObject(object)), nil
	}
}

func NewEncodePrintList() types.EncodeFunc {
	return func(metadata types.IOMetadata, object interface{}) (types.IOMetadata, []byte, error) {
		stringList, err := SprintListOfStrings(object)
		if err != nil {
			return types.IOMetadata{}, nil, fmt.Errorf("cannot print object: %v", err)
		}

		var str string
		for _, uri := range stringList {
			str += fmt.Sprintf("%s\n", uri)
		}

		return metadata, []byte(strings.TrimSuffix(str, "\n")), nil
	}
}

func NewEncodePrintAsListOfLists() types.EncodeFunc {
	return func(metadata types.IOMetadata, object interface{}) (types.IOMetadata, []byte, error) {
		parentList, ok := objects.AsList(object)
		if !ok {
			return types.IOMetadata{}, nil, fmt.Errorf("cannot print object: not a list")
		}

		var str string
		for idx, childListObject := range parentList {
			stringList, err := SprintListOfStrings(childListObject)
			if err != nil {
				return types.IOMetadata{}, nil, fmt.Errorf("cannot print object: %v", err)
			}

			for _, s := range stringList {
				str += fmt.Sprintf("%d: %s\n", idx, s)
			}
		}

		return metadata, []byte(strings.TrimSuffix(str, "\n")), nil
	}
}

func NewEncodeAsHexString() types.EncodeFunc {
	return func(metadata types.IOMetadata, object interface{}) (types.IOMetadata, []byte, error) {
		str, ok := objects.AsString(object)
		if !ok {
			return types.IOMetadata{}, nil, fmt.Errorf("not a string")
		}

		return metadata, []byte(hex.EncodeToString([]byte(str))), nil
	}
}
