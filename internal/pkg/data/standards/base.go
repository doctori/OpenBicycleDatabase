package standards

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"gorm.io/gorm"
)

// TODO remove this shit from here
// this should be controlled by the API not the datamodel
const defaultPerPage int = 30

// StandardInt interface define all the method that a standard need to have to be a
// real standard struct
type StandardInt interface {
	GetName() string
	//	GetCountry() string
	GetCode() string
	GetID() uint
	//	GetType() string
	//	Get() string
	IsNul() bool
	Get(db *gorm.DB, values url.Values, id int, adj string) (int, interface{})
	Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{})
	Put(db *gorm.DB, values url.Values, body io.ReadCloser) (int, interface{})
	Delete(db *gorm.DB, values url.Values, id int) (int, interface{})
	Save(db *gorm.DB) (err error)
}

// Standard define the generic common Standard properties
type Standard struct {
	// add basic ID/Created@/Updated@/Delete@ through Gorm
	gorm.Model `formType:"-"`
	// TODO : this is a embded struct on other structs,
	//  find a way to create indexes on each struct that embed this struct
	Name        string `formType:"string" `
	Country     string `formType:"country"`
	Code        string `formType:"string"`
	Type        string `formType:"string"`
	Description string `formType:"string"`
}

// FieldForm holds the defintion of the field on the Form side (UI)
type FieldForm struct {
	Name       string
	Label      string
	Type       string
	Unit       string
	NestedType string
}

// IsNul return true if the the standard is empty
func (s *Standard) IsNul() bool {
	if s.Name != "" && s.Type != "" {
		return false
	}
	return true
}

// GetName returns the Name
func (s *Standard) GetName() string {
	return s.Name
}

// GetCode return the Code
func (s *Standard) GetCode() string {
	return s.Code
}

// GetID return the ID
func (s *Standard) GetID() uint {
	return s.ID
}

// Delete Standard will remove the Standard struct in the database
func (Standard) Delete(db *gorm.DB, values url.Values, id int, standardType StandardInt) (int, interface{}) {
	// retrieve the Standard
	standard := reflect.New(reflect.TypeOf(standardType))
	err := db.First(&standard, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 404, "Standard not found"
	}
	// TODO : implement Delete method
	db.Delete(standard)
	return 200, ""
}

// GetFormFields will return the custom fields of a Standard to be used in the UI
func GetFormFields(s StandardInt) map[string]FieldForm {
	fields := make(map[string]FieldForm)
	t := reflect.TypeOf(s).Elem()
	log.Printf("Type: %s\n", t.Name())
	log.Printf("Kind: %s\n", t.Kind())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		typeTag := field.Tag.Get("formType")
		if typeTag == "-" {
			continue
		}
		unitTag := field.Tag.Get("formUnit")

		log.Printf("%d, %v (%v), tag : '%v'\n", i, field.Name, field.Type.Name(), typeTag)
		nestedType := ""
		// non standard type
		if typeTag == "nested" || typeTag == "nestedArray" {
			nestedType = field.Type.Name()
			if nestedType == "" {
				nestedType = field.Type.Elem().Name()
			}
		}

		if field.Type.Name() == "" {
			log.Println(field.Type.Elem())
			nestedType = field.Type.Elem().Name()

		}
		fields[field.Name] = FieldForm{
			Name:       field.Name,
			Type:       typeTag,
			Unit:       unitTag,
			NestedType: nestedType,
		}

	}
	return fields
}

// Get Standard return the requests Standards (given the type of standard requested)
func (Standard) Get(db *gorm.DB, values url.Values, id int, standardType StandardInt, adj string) (int, interface{}) {
	log.Printf("having Get for standard [%#v] with ID : %d", standardType, id)
	page := values.Get("page")
	perPage := values.Get("per_page")
	structOnly := values.Get("struct_only")
	// if the request just want the struct we'll respond an new struct only
	if structOnly != "" {
		return 200, GetFormFields(standardType)
	}
	// Let Display All that We Have
	// Someday Pagination will be there
	if id == 0 {
		return 200, GetAll(db, page, perPage, standardType)
	}
	err := db.Model(standardType).First(&standardType, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 404, "Standard not found"
	}
	return 200, standardType
}

// GetAll will returns al the standards
func GetAll(db *gorm.DB, page string, perPage string, standardType StandardInt) interface{} {
	ipage, err := strconv.Atoi(page)
	if err != nil {
		ipage = 0
	}
	// Retrieve the per_page arg, if not a number default to 30
	iperPage, err := strconv.Atoi(perPage)
	if err != nil {
		iperPage = defaultPerPage
	}

	//db.Preload("Components").Preload("Components.Standards").Find(&bikes) // Don't Need to load every Component for the main List
	sType := reflect.TypeOf(standardType)
	standards := reflect.New(reflect.SliceOf(sType)).Interface()

	//var cr []ChainRing
	log.Printf("%#v\n", standards)
	db.Model(standardType).
		Offset(ipage * iperPage).
		Limit(iperPage).
		Find(standards)
	log.Printf("%#v\n", standards)
	return standards
}

// Post will save the BBStandard
func (Standard) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string, standardType StandardInt) (int, interface{}) {
	body := request.Body
	log.Printf("Received args : \n\t %+v\n", values)
	decoder := json.NewDecoder(body)
	log.Println(reflect.TypeOf(standardType))
	standard := reflect.New(reflect.TypeOf(standardType).Elem()).Interface()

	err := decoder.Decode(standard)
	if err != nil {
		log.Printf("Could not unmarshal object : %s", err)
		return http.StatusBadRequest, fmt.Sprintf("Could not unmarshal object : %s", err)
	}

	standardTyped := standard.(StandardInt)
	if standardTyped.IsNul() {
		return http.StatusBadRequest, "The object is null"
	}
	err = standardTyped.Save(db)
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("Could Not Save the Standard : \n\t %s", err.Error())
	}
	return http.StatusAccepted, standardTyped
}

// Put updates Standard
func (Standard) Put(db *gorm.DB, values url.Values, body io.ReadCloser, standardType StandardInt) (int, interface{}) {
	fmt.Printf("Received args : \n\t %+v\n", values)
	decoder := json.NewDecoder(body)
	standard := reflect.New(reflect.TypeOf(standardType)).Elem().Interface().(StandardInt)
	err := decoder.Decode(&standard)
	if err != nil {
		log.Printf("Could not unmarshal object : %s", err)
		return http.StatusBadRequest, fmt.Sprintf("Could not unmarshal object : %s", err)

	}
	log.Println(standard)
	err = standard.Save(db)
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("Could Not Save the Standard : \n\t%s", err.Error())
	}
	return http.StatusOK, standard
}
