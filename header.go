package cpio

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
)

const OldMagic = 070707
const NewMagic = "070701"
const Trailer = "TRAILER!!!"

// Magic
// The integer value octal 070707.  This value can be used to determine
// whether this archive is written with little-endian or big-endian
// integers.
//
// Dev, Ino
// The device and inode numbers from the disk.  These are used by
// programs that read cpio archives to determine when two entries
// refer to the same file.  Programs that synthesize cpio archives
// should be careful to set these to distinct values for each entry.
//
// Mode
// The mode specifies both the regular permissions and the file
// type.  It consists of several bit fields as follows:
// 	0170000  This masks the file type bits.
// 	0140000  File type value for sockets.
// 	0120000  File type value for symbolic links.  For symbolic links,
// the link body is stored as file data.
// 	0100000  File type value for regular files.
// 	0060000  File type value for block special devices.
// 	0040000  File type value for directories.
// 	0020000  File type value for character special devices.
// 	0010000  File type value for named pipes or FIFOs.
// 	0004000  SUID bit.
// 	0002000  SGID bit.
// 	0001000  Sticky bit.  On some systems, this modifies the behavior
// of executables and/or directories.
// 	0000777  The lower 9 bits specify read/write/execute permissions
// for world, group, and user following standard POSIX conventions.
//
// UID, GID
// The numeric user id and group id of the owner.
//
// NLink   The number of links to this file.	Directories always have a
// value of at least two here.  Note that hardlinked files include
// file data with every copy in the archive.
//
// RDev    For block special and character special entries, this field contains
// the associated device number.  For all other entry types,
// it should be set to zero by writers and ignored by readers.
//
// MTime   Modification time of the file, indicated as the number of seconds
// since the start of the epoch, 00:00:00 UTC January 1, 1970.  The
// four-byte integer is stored with the most-significant 16 bits
// first followed by the least-significant 16 bits.  Each of the two
// 16 bit values are stored in machine-native byte order.
//
// NameSize
// The number of bytes in the pathname that follows the header.
// This count includes the trailing NUL byte.
//
// FileSize
// The size of the file.  Note that this archive format is limited
// to four gigabyte file sizes.  See mtime above for a description
// of the storage of four-byte integers.
type HeaderOldCpio struct {
	Magic    uint64
	Dev      uint16
	Ino      uint16
	UID      uint16
	GID      uint16
	NLink    uint16
	RDev     uint16
	MTime    [2]uint16
	NameSize uint16
	FileSize [2]uint16
}

type CpioHeader struct {
	Magic     string
	Ino       uint64
	Mode      uint16
	UID       uint16
	GID       uint16
	NLink     uint16
	MTime     uint16
	FileSize  uint16
	DevMajor  uint16
	DevMinor  uint16
	RDevMajor uint16
	RDevMinor uint16
	NameSize  uint16
	Check     uint16
}

type BinaryReader struct {
	reader io.Reader
}

func (r *BinaryReader) Read(buf interface{}) error {
	return binary.Read(r.reader, binary.BigEndian, buf)
}

func (r *BinaryReader) ReadField() (int64, error) {
	bytes := make([]byte, 8)

	err := r.Read(&bytes)
	if err != nil {
		return 0, err
	}

	data, err := strconv.ParseInt(string(bytes), 16, 0)
	if err != nil {
		return 0, err
	}
	return data, nil
}

func GetHeader(r io.Reader) (*CpioHeader, error) {
	header := &CpioHeader{}

	//binaryReader := &BinaryReader{reader: r}

	magic := make([]byte, 6)
	_, err := r.Read(magic)
	if err != nil {
		return nil, err
	}
	if string(magic) != NewMagic {
		return nil, fmt.Errorf("Not cpio!\n")
	}
	//
	//header.Ino, err = binaryReader.ReadField()
	//if err != nil {
	//	return nil, err
	//}

	return header, nil
}
