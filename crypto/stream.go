package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"hash"
	"io"

	"github.com/blend/go-sdk/ex"
)

// NewStreamEncryptor creates a new stream encrypter
func NewStreamEncryptor(key []byte, plainText io.Reader) (*StreamEncryptor, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ex.New(err)
	}
	iv := make([]byte, block.BlockSize())
	_, err = rand.Read(iv)
	if err != nil {
		return nil, ex.New(err)
	}
	stream := cipher.NewCTR(block, iv)
	mac := hmac.New(sha256.New, key)
	return &StreamEncryptor{
		Source: plainText,
		Block:  block,
		Stream: stream,
		Mac:    mac,
		IV:     iv,
	}, nil
}

// NewStreamDecryptor creates a new stream decrypter
func NewStreamDecryptor(key []byte, meta StreamMeta, cipherText io.Reader) (*StreamDecryptor, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ex.New(err)
	}
	stream := cipher.NewCTR(block, meta.IV)
	mac := hmac.New(sha256.New, key)
	return &StreamDecryptor{
		Source: cipherText,
		Block:  block,
		Stream: stream,
		Mac:    mac,
		Meta:   meta,
	}, nil
}

// StreamEncryptor is an encrypter for a stream of data with authentication
type StreamEncryptor struct {
	Source io.Reader
	Block  cipher.Block
	Stream cipher.Stream
	Mac    hash.Hash
	IV     []byte
}

// StreamDecryptor is a decrypter for a stream of data with authentication
type StreamDecryptor struct {
	Source io.Reader
	Block  cipher.Block
	Stream cipher.Stream
	Mac    hash.Hash
	Meta   StreamMeta
}

// Read encrypts the bytes of the inner reader and places them into p
func (s *StreamEncryptor) Read(p []byte) (int, error) {
	n, readErr := s.Source.Read(p)
	if n > 0 {
		err := writeHash(s.Mac, p[:n])
		if err != nil {
			return n, ex.New(err)
		}
		s.Stream.XORKeyStream(p[:n], p[:n])
		return n, readErr
	}
	return 0, io.EOF
}

// Meta returns the encrypted stream metadata for use in decrypting. This should only be called after the stream is finished
func (s *StreamEncryptor) Meta() StreamMeta {
	return StreamMeta{IV: s.IV, Hash: s.Mac.Sum(nil)}
}

// Read reads bytes from the underlying reader and then decrypts them
func (s *StreamDecryptor) Read(p []byte) (int, error) {
	n, readErr := s.Source.Read(p)
	if n > 0 {
		s.Stream.XORKeyStream(p[:n], p[:n])
		err := writeHash(s.Mac, p[:n])
		if err != nil {
			return n, ex.New(err)
		}
		return n, readErr
	}
	return 0, io.EOF
}

// Authenticate verifys that the hash of the stream is correct. This should only be called after processing is finished
func (s *StreamDecryptor) Authenticate() error {
	if !hmac.Equal(s.Meta.Hash, s.Mac.Sum(nil)) {
		return ex.New("authentication failed")
	}
	return nil
}

func writeHash(mac hash.Hash, p []byte) error {
	m, err := mac.Write(p)
	if err != nil {
		return ex.New(err)
	}
	if m != len(p) {
		return ex.New("could not write all bytes to hmac")
	}
	return nil
}

func checkedWrite(dst io.Writer, p []byte) (int, error) {
	n, err := dst.Write(p)
	if err != nil {
		return n, ex.New(err)
	}
	if n != len(p) {
		return n, ex.New("unable to write all bytes")
	}
	return len(p), nil
}

// StreamMeta is metadata about an encrypted stream
type StreamMeta struct {
	// IV is the initial value for the crypto function
	IV []byte
	// Hash is the sha256 hmac of the stream
	Hash []byte
}
