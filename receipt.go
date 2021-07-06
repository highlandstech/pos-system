package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"text/template"

	"github.com/phin1x/go-ipp"
	"github.com/skip2/go-qrcode"
	"github.com/wisepythagoras/pos-system/crypto"
)

// ReceiptTmplData describes the template data for the receipt.
type ReceiptTmplData struct {
	Total    float64
	Order    *OrderJSON
	Products *[]AggregateProduct
	Qrcode   string
	Name     string
	Address1 string
	Address2 string
}

// Receipt describes the receipt object. It should contain an order and the config.
type Receipt struct {
	Order  *OrderJSON
	Total  float64
	Config *Config
	Client *ipp.IPPClient
}

// ConnectToPrinter connects to the printer server.
func (r *Receipt) ConnectToPrinter() {
	server := r.Config.Printer.Server
	port := r.Config.Printer.Port
	username := r.Config.Printer.Username
	password := r.Config.Printer.Password

	// Connect to the printer server.
	r.Client = ipp.NewIPPClient(server, port, username, password, true)
}

// Print prints the receipt.
func (r *Receipt) Print() (int, error) {
	// Connect to the printer if the connection died.
	if r.Client.TestConnection() != nil {
		r.ConnectToPrinter()
	}

	// Create the file.
	data, err := ioutil.ReadFile("templates/receipt.html")

	if err != nil {
		return 99, err
	}

	t, err := template.New("receipt").Parse(string(data))

	if err != nil {
		return 99, err
	}

	fileName := "receipt-" + strconv.Itoa(int(r.Order.ID))

	// Create a file for the rendered receipt.
	receiptFile, err := os.OpenFile("receipts/"+fileName+".html", os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		return 99, err
	}

	defer receiptFile.Close()

	encryptJSON := `{"i": ` + strconv.Itoa(int(r.Order.ID)) + `}`
	encryptedId, _ := crypto.EncryptGCM([]byte(encryptJSON), []byte(r.Config.Key))
	hexEncryptedId := hex.EncodeToString(encryptedId)
	png, _ := qrcode.Encode(hexEncryptedId, qrcode.Medium, 160)

	// This is the new buffer that will contain the executedtemplate.
	buff := new(bytes.Buffer)

	var aggregateProducts []AggregateProduct
	var aggregateMap map[uint64]uint = make(map[uint64]uint)
	var products map[uint64]ProductJSON = make(map[uint64]ProductJSON)

	// Find how many of each products we have.
	for _, product := range r.Order.Products {
		if val, ok := aggregateMap[product.ID]; ok {
			aggregateMap[product.ID] = val + 1
		} else {
			aggregateMap[product.ID] = 1
			products[product.ID] = product
		}
	}

	// Now compose the aggregate product array.
	for productId, quantity := range aggregateMap {
		product := products[productId]
		aggregateProduct := AggregateProduct{
			Quantity: quantity,
			ID:       product.ID,
			Name:     product.Name,
			Price:    float64(quantity) * product.Price,
			Type:     product.Type,
		}

		aggregateProducts = append(aggregateProducts, aggregateProduct)
	}

	tmplData := ReceiptTmplData{
		Total:    r.Total,
		Order:    r.Order,
		Qrcode:   base64.StdEncoding.EncodeToString(png),
		Products: &aggregateProducts,
		Name:     r.Config.Name,
		Address1: r.Config.Address1,
		Address2: r.Config.Address2,
	}

	// Execute the receipt.
	err = t.Execute(buff, tmplData)

	receiptFile.WriteString(buff.String())

	cwd, _ := os.Getwd()

	input := cwd + "/receipts/" + fileName + ".html"
	output := cwd + "/receipts/" + fileName + ".pdf"

	// Now we want to convert the receipt to a PDF so that the printer can render it.
	cmd := exec.Command("wkhtmltopdf", "--page-width", "60", "--page-height", "200", input, output)
	err = cmd.Start()

	if err != nil {
		return 99, err
	}

	err = cmd.Wait()

	if err != nil {
		return 99, err
	}

	// Finally, call the printer to print the receipt.
	return r.Client.PrintFile(output, r.Config.Printer.Name, map[string]interface{}{})
}