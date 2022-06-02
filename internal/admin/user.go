package admin

import (
	"bytes"
	"cic-dw/pkg/address"
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/grassrootseconomics/cic-go/meta"
	"github.com/labstack/echo/v4"
	"github.com/mapaiva/vcard-go"
)

type metaRes struct {
	Person meta.PersonResponse `json:"person"`
	Name   string              `json:"name"`
}

func handlePhone2Address(c echo.Context) error {
	var (
		api   = c.Get("api").(*api)
		phone = c.Param("phone")

		address string
	)

	if err := api.db.QueryRow(context.Background(), api.q["phone-2-address"], phone).Scan(&address); err != nil {
		return c.String(http.StatusNotFound, "phone not found")
	}

	return c.String(http.StatusOK, address)
}

func handleAddress2Phone(c echo.Context) error {
	var (
		api = c.Get("api").(*api)

		phone string
	)

	sarafuAddress, err := address.SarafuAddress(c.Param("address"))
	if err != nil {
		return err
	}

	if err := api.db.QueryRow(context.Background(), api.q["address-2-phone"], sarafuAddress).Scan(&phone); err != nil {
		return c.String(http.StatusNotFound, "address not found")
	}

	return c.String(http.StatusOK, phone)
}

func handleMetaProxy(c echo.Context) error {
	var (
		api     = c.Get("api").(*api)
		address = c.Param("address")
	)

	person, err := api.m.GetPersonMetadata(address)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return c.String(http.StatusNotFound, "meta resource not found")
		} else {
			return err
		}
	}

	vCard, err := parseVCard(person.VCard)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &metaRes{
		Person: person,
		Name:   vCard.FormattedName,
	})
}

func parseVCard(encodedVCard string) (vcard.VCard, error) {
	data, err := base64.StdEncoding.DecodeString(encodedVCard)
	if err != nil {
		return vcard.VCard{}, err
	}

	reader := bytes.NewReader(data)

	vCards, err := vcard.GetVCardsByReader(reader)
	if err != nil {
		return vcard.VCard{}, nil
	}

	return vCards[0], nil
}
