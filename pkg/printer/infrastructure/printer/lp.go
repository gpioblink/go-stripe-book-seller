package printer

import "log"

type PrinterInterface struct {
}

func NewPrinterInterface() *PrinterInterface {
	return &PrinterInterface{}
}

func (p *PrinterInterface) PrintTitle(query string) error {
	log.Printf("print title: %s", query)
	return nil
}

func (p *PrinterInterface) PrintLine(query string) error {
	log.Printf("print line: %s", query)
	return nil
}

func (p *PrinterInterface) PrintPicture(srcUrl string) error {
	log.Printf("print picture: %s", srcUrl)
	return nil
}

func (p *PrinterInterface) PrintBarCode(query string) error {
	log.Printf("print barcode: %s", query)
	return nil
}

func (p *PrinterInterface) Cut() error {
	log.Printf("print cut")
	return nil
}
