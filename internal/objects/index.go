package objects

import (
	"bytes"
	"encoding/binary"
	"errors"
	"mcavazotti/git-go/internal/repo"
	"os"
)

type IndexEntry struct {
	Ctime_s           uint32
	Ctime_ns          uint32
	Mtime_s           uint32
	Mtime_ns          uint32
	Dev               uint32
	Ino               uint32
	Mode_type         uint32
	Mode_perms        uint32
	Uid               uint32
	Gid               uint32
	Fsize             uint32
	Sha               []byte
	Flag_assume_valid bool
	Flag_stage        byte
	Name              string
}

type IndexObject struct {
	Version uint32
	Entries []IndexEntry
}

func ReadIndex(repository *repo.Repository) (IndexObject, error) {
	indexPath := repository.RepoPath("index")

	if _, err := os.Stat(indexPath); errors.Is(err, os.ErrNotExist) {
		return IndexObject{Version: 2}, nil
	}
	indexObj := IndexObject{}

	data, err := os.ReadFile(indexPath)
	if err != nil {
		return indexObj, err
	}

	header := data[:12]

	if string(header[:4]) != "DIRC" {
		return indexObj, errors.New("Invalid index file")
	}

	indexObj.Version = binary.BigEndian.Uint32(header[4:8])

	numRecords := binary.BigEndian.Uint32(header[8:12])

	content := data[12:]
	r := bytes.NewReader(content)

	for i := uint32(0); i < numRecords; i++ {
		entry, err := readIndexEntry(r)
		if err != nil {
			return IndexObject{}, err
		}
		indexObj.Entries = append(indexObj.Entries, entry)
	}
	return indexObj, nil
}

func readIndexEntry(reader *bytes.Reader) (IndexEntry, error) {
	twoByteBuffer := make([]byte, 4)
	fourByteBuffer := make([]byte, 4)
	entry := IndexEntry{}
	bytesRead := 0

	n, err := reader.Read(fourByteBuffer)
	if err != nil {
		return entry, err
	}
	bytesRead += n
	entry.Ctime_s = binary.BigEndian.Uint32(fourByteBuffer)

	n, err = reader.Read(fourByteBuffer)
	if err != nil {
		return entry, err
	}
	bytesRead += n
	entry.Ctime_ns = binary.BigEndian.Uint32(fourByteBuffer)

	n, err = reader.Read(fourByteBuffer)
	if err != nil {
		return entry, err
	}
	bytesRead += n
	entry.Mtime_s = binary.BigEndian.Uint32(fourByteBuffer)

	n, err = reader.Read(fourByteBuffer)
	if err != nil {
		return entry, err
	}
	bytesRead += n
	entry.Mtime_ns = binary.BigEndian.Uint32(fourByteBuffer)

	n, err = reader.Read(fourByteBuffer)
	if err != nil {
		return entry, err
	}
	bytesRead += n
	entry.Dev = binary.BigEndian.Uint32(fourByteBuffer)

	n, err = reader.Read(fourByteBuffer)
	if err != nil {
		return entry, err
	}
	bytesRead += n
	entry.Ino = binary.BigEndian.Uint32(fourByteBuffer)

	// Ignore two bytes
	n, err = reader.Read(twoByteBuffer)
	if err != nil {
		return entry, err
	}
	bytesRead += n
	if binary.BigEndian.Uint32(twoByteBuffer) != 0 {
		return entry, errors.New("Invalid index file")
	}

	n, err = reader.Read(twoByteBuffer)
	if err != nil {
		return entry, err
	}
	bytesRead += n
	mode := binary.BigEndian.Uint32(twoByteBuffer)
	entry.Mode_type = mode >> 12
	entry.Mode_perms = mode & 0b0000000111111111

	n, err = reader.Read(fourByteBuffer)
	if err != nil {
		return entry, err
	}
	bytesRead += n
	entry.Uid = binary.BigEndian.Uint32(fourByteBuffer)

	n, err = reader.Read(twoByteBuffer)
	if err != nil {
		return entry, err
	}
	bytesRead += n
	entry.Gid = binary.BigEndian.Uint32(fourByteBuffer)

	n, err = reader.Read(twoByteBuffer)
	if err != nil {
		return entry, err
	}
	bytesRead += n
	entry.Fsize = binary.BigEndian.Uint32(fourByteBuffer)

	entry.Sha = make([]byte, 20)
	n, err = reader.Read(entry.Sha)
	if err != nil {
		return entry, err
	}
	bytesRead += n

	n, err = reader.Read(twoByteBuffer)
	if err != nil {
		return entry, err
	}
	bytesRead += n
	flags := binary.BigEndian.Uint32(twoByteBuffer)

	entry.Flag_assume_valid = (flags & 0b1000000000000000) != 0
	entry.Flag_stage = byte((flags & 0b0011000000000000) >> 8)

	nameLength := flags & 0b0000111111111111

	if nameLength < 0xfff {
		nameBuffer := make([]byte, nameLength+1)
		n, err = reader.Read(nameBuffer)
		if err != nil {
			return entry, err
		}
		bytesRead += n
		entry.Name = string(nameBuffer)
	} else {
		nameBuffer := bytes.NewBufferString("")
		for b, err := reader.ReadByte(); b != 0x00; {
			if err != nil {
				return entry, err
			}
			nameBuffer.WriteByte(b)
			bytesRead++
		}
		entry.Name = string(nameBuffer.String())
		bytesRead++
	}
	padding := 8 - bytesRead%8

	if padding < 0 {
		tmpBuffer := make([]byte, padding)
		_, err = reader.Read(tmpBuffer)
		if err != nil {
			return entry, err
		}
	}
	return entry, nil
}
