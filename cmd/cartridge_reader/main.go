package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jmontupet/gbcore/internal/pkg/cartridge"
)

func main() {
	// Parse command options
	flag.Parse()
	tails := flag.Args()
	if len(tails) != 1 {
		log.Fatalln("Need valid File name")
	}
	// Read GameData
	gameData, err := ioutil.ReadFile(tails[0])
	if err != nil {
		log.Fatal(err)
	}

	cart, err := cartridge.NewCartridge(gameData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Title:", cartridge.ReadTitle(cart))
	fmt.Println("Manufacturer Code:", cartridge.ReadManufacturerCode(cart))
	fmt.Printf("CGB Flag: 0x%0.2X\n", cartridge.ReadCGBFlag(cart))
	fmt.Printf("New Licensee Code: %s\n", cartridge.ReadNewLicenseeCode(cart))
	fmt.Printf("SGB Flag: 0x%0.2X\n", cartridge.ReadSGBFlag(cart))
	fmt.Printf("CartridgeType: 0x%0.2X\n", cartridge.ReadCartridgeType(cart))
	fmt.Printf("ROM Size: 0x%0.2X\n", cartridge.ReadROMSize(cart))
	fmt.Printf("RAM Size: 0x%0.2X\n", cartridge.ReadRAMSize(cart))
	fmt.Printf("Destination Code: 0x%0.2X\n", cartridge.ReadDestinationCode(cart))
	fmt.Printf("Old Licensee Code: 0x%0.2X\n", cartridge.ReadOldLicenseeCode(cart))
	fmt.Printf("ROM Version: 0x%0.2X\n", cartridge.ReadROMVersion(cart))
	fmt.Printf("Header Checksum (READ): 0x%0.2X\n", cartridge.ReadHeaderChecksum(cart))
	fmt.Printf("Header Checksum (CALC): 0x%0.2X\n", cartridge.CalcHeaderChecksum(cart))
	fmt.Printf("ROM Checksum: 0x%0.4X\n", cartridge.ReadROMChecksum(cart))
	fmt.Printf("Nintendo Logo: %X\n", cartridge.ReadNintendoLogo(cart))
}
