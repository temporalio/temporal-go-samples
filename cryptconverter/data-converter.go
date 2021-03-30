package cryptconverter

import (
	"fmt"

	commonpb "go.temporal.io/api/common/v1"

	"go.temporal.io/sdk/converter"
)

const (
	// MetadataEncryptionKeyId is "encryption-key-id"
	MetadataEncryptionKeyId = "encryption-key-id"

	// MetadataEncodingEncrypted is "binary/encrypted"
	MetadataEncodingEncrypted = "binary/encrypted"
)

// CryptDataConverter implements DataConverter using AES Crypt.
type CryptDataConverter struct {
	dataConverter converter.DataConverter
	context       CryptContext
}

type Context interface {
	Value(interface{}) interface{}
}

func (dc *CryptDataConverter) WithValue(v interface{}) converter.DataConverter {
	ctx, ok := v.(Context)
	if !ok {
		return dc
	}

	if val := ctx.Value(PropagateKey); val != nil {
		return &CryptDataConverter{
			dataConverter: converter.WithValue(dc.dataConverter, ctx),
			context:       val.(CryptContext),
		}
	}

	return dc
}

func (dc *CryptDataConverter) getKey(keyId string) (key []byte) {
	// Key can be fetched from KMS or other secure storage.
	return []byte("test-key-test-key-test-key-test!")
}

// NewCryptDataConverter creates a new instance of CryptDataConverter wrapping a DataConverter
func NewCryptDataConverter(dataConverter converter.DataConverter) *CryptDataConverter {
	return &CryptDataConverter{
		dataConverter: dataConverter,
	}
}

// ToPayloads converts a list of values.
func (dc *CryptDataConverter) ToPayloads(values ...interface{}) (*commonpb.Payloads, error) {
	if dc.context.KeyId == "" {
		return dc.dataConverter.ToPayloads(values)
	}
	key := dc.getKey(dc.context.KeyId)

	result := &commonpb.Payloads{}

	for i, value := range values {
		payload, err := dc.dataConverter.ToPayload(value)
		if err != nil {
			return nil, fmt.Errorf("values[%d]: %w", i, err)
		}

		payload, err = dc.encryptPayload(payload, dc.context.KeyId, key)
		if err != nil {
			return nil, fmt.Errorf("values[%d]: %w", i, err)
		}

		result.Payloads = append(result.Payloads, payload)
	}

	return result, nil
}

func (dc *CryptDataConverter) encryptPayload(payload *commonpb.Payload, keyId string, key []byte) (*commonpb.Payload, error) {
	serializedPayload, err := payload.Marshal()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", converter.ErrUnableToEncode, err)
	}

	encryptedPayload, err := encrypt(serializedPayload, key)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", converter.ErrUnableToEncode, err)
	}

	return &commonpb.Payload{
		Metadata: map[string][]byte{
			converter.MetadataEncoding: []byte(MetadataEncodingEncrypted),
			MetadataEncryptionKeyId:    []byte(keyId),
		},
		Data: encryptedPayload,
	}, nil
}

// ToPayload converts single value to payload.
func (dc *CryptDataConverter) ToPayload(value interface{}) (*commonpb.Payload, error) {
	return dc.dataConverter.ToPayload(value)
}

func isEncryptedPayload(payload *commonpb.Payload) bool {
	metadata := payload.GetMetadata()
	if metadata == nil {
		return false
	}

	if encoding, ok := metadata[converter.MetadataEncoding]; ok {
		return string(encoding) == MetadataEncodingEncrypted
	}

	return false
}

// FromPayloads converts to a list of values of different types.
func (dc *CryptDataConverter) FromPayloads(payloads *commonpb.Payloads, valuePtrs ...interface{}) error {
	for i, payload := range payloads.GetPayloads() {
		var err error

		if i >= len(valuePtrs) {
			break
		}

		if isEncryptedPayload(payload) {
			payload, err = dc.decryptPayload(payload)
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

func (dc *CryptDataConverter) decryptPayload(payload *commonpb.Payload) (*commonpb.Payload, error) {
	metadata := payload.GetMetadata()
	if metadata == nil {
		return nil, converter.ErrMetadataIsNotSet
	}

	keyId, ok := metadata[MetadataEncryptionKeyId]
	if !ok {
		return nil, fmt.Errorf("%w: %s", converter.ErrUnableToDecode, "no encryption key id")
	}

	key := dc.getKey(string(keyId))
	serializedPayload, err := decrypt(payload.GetData(), key)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", converter.ErrUnableToDecode, err)
	}

	result := &commonpb.Payload{}
	err = result.Unmarshal(serializedPayload)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", converter.ErrUnableToDecode, err)
	}

	return result, nil
}

// FromPayload converts single value from payload.
func (dc *CryptDataConverter) FromPayload(payload *commonpb.Payload, valuePtr interface{}) error {
	return dc.dataConverter.FromPayload(payload, valuePtr)
}

// ToStrings converts payloads object into human readable strings.
func (dc *CryptDataConverter) ToStrings(payloads *commonpb.Payloads) []string {
	var result []string
	for _, payload := range payloads.GetPayloads() {
		result = append(result, dc.ToString(payload))
	}

	return result
}

// ToString converts payload object into human readable string.
func (dc *CryptDataConverter) ToString(payload *commonpb.Payload) string {
	if isEncryptedPayload(payload) {
		var err error
		payload, err = dc.decryptPayload(payload)
		if err != nil {
			return err.Error()
		}
	}

	return dc.dataConverter.ToString(payload)
}
