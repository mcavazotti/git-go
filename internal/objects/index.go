package objects

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"mcavazotti/git-go/internal/repo"
	"mcavazotti/git-go/internal/shared"
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
		shared.VerbosePrintln("No signature DIRC")
		return indexObj, errors.New("Invalid index file")
	}

	indexObj.Version = binary.BigEndian.Uint32(header[4:8])
	shared.VerbosePrintln("Version", indexObj.Version)

	numRecords := binary.BigEndian.Uint32(header[8:12])
	shared.VerbosePrintln("Num records", numRecords)

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
	shared.VerbosePrintln("\nReading index entry")
	twoByteBuffer := make([]byte, 2)
	fourByteBuffer := make([]byte, 4)
	entry := IndexEntry{}
	bytesRead := 0

	n, err := reader.Read(fourByteBuffer)
	if err != nil {
		shared.VerbosePrintln("failed to read CTIME_s")
		return entry, err
	}
	bytesRead += n
	entry.Ctime_s = binary.BigEndian.Uint32(fourByteBuffer)
	shared.VerbosePrintln("CTIME_S", entry.Ctime_s)

	n, err = reader.Read(fourByteBuffer)
	if err != nil {
		shared.VerbosePrintln("failed to read CTIME_ns")
		return entry, err
	}
	bytesRead += n
	entry.Ctime_ns = binary.BigEndian.Uint32(fourByteBuffer)
	shared.VerbosePrintln("CTIME_NS", entry.Ctime_ns)

	n, err = reader.Read(fourByteBuffer)
	if err != nil {
		shared.VerbosePrintln("failed to read MTIME_s")
		return entry, err
	}
	bytesRead += n
	entry.Mtime_s = binary.BigEndian.Uint32(fourByteBuffer)
	shared.VerbosePrintln("MTIME_S", entry.Mtime_s)

	n, err = reader.Read(fourByteBuffer)
	if err != nil {
		shared.VerbosePrintln("failed to read MTIME_ns")
		return entry, err
	}
	bytesRead += n
	entry.Mtime_ns = binary.BigEndian.Uint32(fourByteBuffer)
	shared.VerbosePrintln("MTIME_NS", entry.Mtime_ns)

	n, err = reader.Read(fourByteBuffer)
	if err != nil {
		shared.VerbosePrintln("failed to read DEV id")
		return entry, err
	}
	bytesRead += n
	entry.Dev = binary.BigEndian.Uint32(fourByteBuffer)
	shared.VerbosePrintln("DEV", entry.Dev)

	n, err = reader.Read(fourByteBuffer)
	if err != nil {
		shared.VerbosePrintln("failed to read ino")
		return entry, err
	}
	bytesRead += n
	entry.Ino = binary.BigEndian.Uint32(fourByteBuffer)
	shared.VerbosePrintln("INO", entry.Ino)

	// Ignore two bytes
	n, err = reader.Read(fourByteBuffer)
	if err != nil {
		shared.VerbosePrintln("failed to read MODE")
		return entry, err
	}
	bytesRead += n
	mode := binary.BigEndian.Uint32(fourByteBuffer)
	shared.VerbosePrintf("Ignored part %b\n", mode>>16)
	shared.VerbosePrintf("Mode %b\n", mode)
	if mode>>16 != 0 {
		shared.VerbosePrintln("ignored bytes not zero")
		return entry, errors.New("Invalid index file")
	}

	entry.Mode_type = mode >> 12
	entry.Mode_perms = mode & 0b0000000111111111
	shared.VerbosePrintf("Mode type %b\n", entry.Mode_type)
	shared.VerbosePrintf("Mode type %o\n", entry.Mode_perms)

	n, err = reader.Read(fourByteBuffer)
	if err != nil {
		shared.VerbosePrintln("failed to read UID")
		return entry, err
	}
	bytesRead += n
	entry.Uid = binary.BigEndian.Uint32(fourByteBuffer)
	shared.VerbosePrintln("UID", entry.Uid)

	n, err = reader.Read(fourByteBuffer)
	if err != nil {
		shared.VerbosePrintln("failed to read GID")
		return entry, err
	}
	bytesRead += n
	entry.Gid = binary.BigEndian.Uint32(fourByteBuffer)
	shared.VerbosePrintln("GID", entry.Gid)

	n, err = reader.Read(fourByteBuffer)
	if err != nil {
		shared.VerbosePrintln("failed to read FSIZE")
		return entry, err
	}
	bytesRead += n
	entry.Fsize = binary.BigEndian.Uint32(fourByteBuffer)
	shared.VerbosePrintln("FSIZE", entry.Fsize)

	entry.Sha = make([]byte, 20)
	n, err = reader.Read(entry.Sha)
	if err != nil {
		shared.VerbosePrintln("failed to read obj name (SHA)")
		return entry, err
	}
	bytesRead += n
	shared.VerbosePrintln("SHA", hex.EncodeToString(entry.Sha))

	n, err = reader.Read(twoByteBuffer)
	if err != nil {
		shared.VerbosePrintln("failed to read flags")
		return entry, err
	}
	bytesRead += n
	flags := binary.BigEndian.Uint16(twoByteBuffer)
	shared.VerbosePrintf("Flags %b\n", flags)

	entry.Flag_assume_valid = (flags & 0b1000000000000000) != 0
	entry.Flag_stage = byte((flags & 0b0011000000000000) >> 8)

	nameLength := flags & 0b0000111111111111
	shared.VerbosePrintf("Name length %x\n", nameLength)

	if nameLength < 0xfff {
		shared.VerbosePrintln("Name length < 0xfff")
		nameBuffer := make([]byte, nameLength+1)
		n, err = reader.Read(nameBuffer)
		if err != nil {
			return entry, err
		}
		bytesRead += n
		entry.Name = string(nameBuffer)
	} else {
		shared.VerbosePrintln("Name length >= 0xfff")
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
	shared.VerbosePrintln("Name", entry.Name)
	if entry.Name[len(entry.Name)-1] == 0 {
		entry.Name = entry.Name[:len(entry.Name)-1]
	}

	padding := (8 - bytesRead%8) % 8
	shared.VerbosePrintln("\nBytes read", bytesRead)
	shared.VerbosePrintln("Padding", padding)

	if padding > 0 {
		tmpBuffer := make([]byte, padding)
		_, err = reader.Read(tmpBuffer)
		if err != nil {
			shared.VerbosePrintln("failed to read padding")
			return entry, err
		}
		shared.VerbosePrintf("bytes in padding %b\n", tmpBuffer)
	}
	return entry, nil
}
