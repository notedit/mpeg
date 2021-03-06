package ts_test

import "testing"
import "io"
import "github.com/32bitkid/bitreader"
import "github.com/32bitkid/mpeg/ts"

func TestPacketParsing(t *testing.T) {
	reader := bitreader.NewBitReader(nullPacketReader())
	packet, err := ts.NewPacket(reader)
	if err != nil {
		t.Fatal(err)
	}
	if packet.PID != nullPacketPID {
		t.Fatalf("unexpected PID. expected %x, got %x", nullPacketPID, packet.PID)
	}
}

func TestEOFAfterPacket(t *testing.T) {
	var err error
	reader := bitreader.NewBitReader(nullPacketReader())
	_, err = ts.NewPacket(reader)
	_, err = ts.NewPacket(reader)
	if err != io.EOF {
		t.Fatal(err)
	}
}

func TestIncompletePacket(t *testing.T) {
	reader := bitreader.NewBitReader(io.LimitReader(nullPacketReader(), 4))
	_, err := ts.NewPacket(reader)
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("unexpected error. expected %v, got %v", io.ErrUnexpectedEOF, err)
	}
}
