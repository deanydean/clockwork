package watches

import (
	"fmt"
	"os"
	"time"

	"github.com/deanydean/clockwork/core"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var snapLen = int32(65536)
var promisc = true
var timeout = time.Second

var is_reading = bool(true)
var not_reading = bool(false)

type NetWatch struct {
	layer      *string
	iface      *string
	pcapHandle *pcap.Handle
	decoder    *gopacket.Decoder
	isReading  *bool
}

func NewNetWatch(layer *string, iface *string) *NetWatch {
	netWatch := new(NetWatch)
	netWatch.iface = iface

	// Open the handle
	var handle, err = pcap.OpenLive(*netWatch.iface, snapLen, promisc, timeout)
	if err != nil {
		fmt.Printf("Failed to open pcap for %s : %s\n", *netWatch.iface, err)
		return nil
	}
	netWatch.pcapHandle = handle

	// Create a decoder
	decoder, ok := gopacket.DecodersByLayerName[*layer]
	if !ok {
		fmt.Printf("Failed to get decoder for layer=%s", *layer)
		return nil
	}
	netWatch.decoder = &decoder

	netWatch.isReading = &not_reading

	return netWatch
}

func (watch *NetWatch) startReading() {
	source := gopacket.NewPacketSource(watch.pcapHandle, *watch.decoder)
	source.Lazy = true
	source.NoCopy = true
	source.DecodeStreamsAsDatagrams = true
	fmt.Fprintln(os.Stderr, "Starting to read packets")

	count := 0
	bytes := int64(0)

	for packet := range source.Packets() {
		count++
		bytes += int64(len(packet.Data()))

		// TODO collect stats
		fmt.Printf("Got packet %s \n", packet)
	}

	fmt.Printf("Completed reading packets")
	watch.isReading = &not_reading
}

// Observe the network
func (watch *NetWatch) Observe() *core.WatchEvent {

	if !*watch.isReading {
		// If we're not reading, start watching now...
		watch.isReading = &is_reading
		go watch.startReading()
	}

	// TODO observe the stats

	return nil
}
