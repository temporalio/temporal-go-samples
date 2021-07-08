package compressionconverter

import (
	"bytes"
	"fmt"
	"io"

	commonpb "go.temporal.io/api/common/v1"

	"go.temporal.io/sdk/converter"

	"github.com/pierrec/lz4"
)

const (
	// MetadataCompressionAlgorithmKey is "compression-algorithm"
	MetadataCompressionAlgorithmKey = "compression-algorithm"

	// MetadataCompressionAlgorithm is "lz4"
	MetadataCompressionAlgorithm = "lz4"

	// MetadataEncodingCompressed is "binary/compressed"
	MetadataEncodingCompressed = "binary/compressed"
)

// CompressionConverter implements a DataConverter using LZ4.
type CompressionConverter struct {
	dataConverter converter.DataConverter
}

// NewCompressionConverter creates a new instance of CompressionConvertter wrapping a DataConverter
func NewCompressionConverter(dataConverter converter.DataConverter) *CompressionConverter {
	return &CompressionConverter{
		dataConverter: dataConverter,
	}
}

// ToPayloads converts a list of values.
func (dc *CompressionConverter) ToPayloads(values ...interface{}) (*commonpb.Payloads, error) {
	result := &commonpb.Payloads{}

	for i, value := range values {
		payload, err := dc.dataConverter.ToPayload(value)
		if err != nil {
			return nil, fmt.Errorf("values[%d]: %w", i, err)
		}

		payload, err = dc.compressPayload(payload)
		if err != nil {
			return nil, fmt.Errorf("values[%d]: %w", i, err)
		}

		result.Payloads = append(result.Payloads, payload)
	}

	return result, nil
}

func (dc *CompressionConverter) compressPayload(payload *commonpb.Payload) (*commonpb.Payload, error) {
	serializedPayload, err := payload.Marshal()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", converter.ErrUnableToEncode, err)
	}

	var compressedPayload bytes.Buffer

	w := lz4.NewWriter(&compressedPayload)
	_, err = w.Write(serializedPayload)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", converter.ErrUnableToEncode, err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", converter.ErrUnableToEncode, err)
	}

	return &commonpb.Payload{
		Metadata: map[string][]byte{
			converter.MetadataEncoding:      []byte(MetadataEncodingCompressed),
			MetadataCompressionAlgorithmKey: []byte(MetadataCompressionAlgorithm),
		},
		Data: compressedPayload.Bytes(),
	}, nil
}

// ToPayload converts single value to payload.
func (dc *CompressionConverter) ToPayload(value interface{}) (*commonpb.Payload, error) {
	return dc.dataConverter.ToPayload(value)
}

func isCompressedPayload(payload *commonpb.Payload) bool {
	metadata := payload.GetMetadata()
	if metadata == nil {
		return false
	}

	if encoding, ok := metadata[converter.MetadataEncoding]; ok {
		return string(encoding) == MetadataEncodingCompressed
	}

	return false
}

// FromPayloads converts to a list of values of different types.
func (dc *CompressionConverter) FromPayloads(payloads *commonpb.Payloads, valuePtrs ...interface{}) error {
	for i, payload := range payloads.GetPayloads() {
		var err error

		if i >= len(valuePtrs) {
			break
		}

		if isCompressedPayload(payload) {
			payload, err = dc.decompressPayload(payload)
			if err != nil {
				return fmt.Errorf("args[%d]: %w", i, err)
			}
		}

		err = dc.dataConverter.FromPayload(payload, valuePtrs[i])
		if err != nil {
			return fmt.Errorf("args[%d]: %w", i, err)
		}
	}

	return nil
}

func (dc *CompressionConverter) decompressPayload(payload *commonpb.Payload) (*commonpb.Payload, error) {
	metadata := payload.GetMetadata()
	if metadata == nil {
		return nil, converter.ErrMetadataIsNotSet
	}

	algorithm, ok := metadata[MetadataCompressionAlgorithmKey]
	if !ok {
		return nil, fmt.Errorf("%w: %s", converter.ErrUnableToDecode, "no compression algorithm metadata")
	}

	if string(algorithm) != MetadataCompressionAlgorithm {
		return nil, fmt.Errorf("%w: %s '%s'", converter.ErrUnableToDecode, "unsupported compression algorithm", algorithm)
	}

	var uncompressedPayload bytes.Buffer

	r := lz4.NewReader(bytes.NewReader(payload.GetData()))
	_, err := io.Copy(&uncompressedPayload, r)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", converter.ErrUnableToDecode, err)
	}

	result := &commonpb.Payload{}
	err = result.Unmarshal(uncompressedPayload.Bytes())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", converter.ErrUnableToDecode, err)
	}

	return result, nil
}

// FromPayload converts single value from payload.
func (dc *CompressionConverter) FromPayload(payload *commonpb.Payload, valuePtr interface{}) error {
	return dc.dataConverter.FromPayload(payload, valuePtr)
}

// ToStrings converts payloads object into human readable strings.
func (dc *CompressionConverter) ToStrings(payloads *commonpb.Payloads) []string {
	var result []string
	for _, payload := range payloads.GetPayloads() {
		result = append(result, dc.ToString(payload))
	}

	return result
}

// ToString converts payload object into human readable string.
func (dc *CompressionConverter) ToString(payload *commonpb.Payload) string {
	if isCompressedPayload(payload) {
		var err error
		payload, err = dc.decompressPayload(payload)
		if err != nil {
			return err.Error()
		}
	}

	return dc.dataConverter.ToString(payload)
}
