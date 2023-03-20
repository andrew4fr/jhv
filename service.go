package treasure

import (
	"encoding/xml"
	"net/http"
	"treasure/internal/rest/model"
)

const individualType string = "Individual"

var busy bool

type Entry struct {
	UID       int    `xml:"uid"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	SdnType   string `xml:"sdnType"`
}

type State struct {
	Code   int
	Info   string
	Result string
}

type Storage interface {
	GetNames(searchName, searchType string) (*model.Persons, error)
	SaveEntry(uid int, firstName, lastName string) error
	IsEmpty() bool
}

type Service struct {
	xmlPath string
	storage Storage
}

func New(path string, stor Storage) *Service {
	return &Service{
		xmlPath: path,
		storage: stor,
	}
}

func (s *Service) GetNames(searchName, searchType string) (*model.Persons, error) {
	return nil, nil
}

func (s *Service) GetState() *State {
	if s.storage.IsEmpty() {
		return &State{
			Result: "false",
			Info:   "empty",
		}
	}

	if busy {
		return &State{
			Result: "false",
			Info:   "updating",
		}
	}

	return &State{
		Result: "true",
		Info:   "ok",
	}
}

func (s *Service) UpdateList() error {
	setBusy()
	defer setFree()

	resp, err := http.Get(s.xmlPath)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	var inElement string

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "sdnEntry" {
				var entry Entry
				decoder.DecodeElement(&entry, &se)
				if entry.SdnType != individualType {
					continue
				}

				err := s.storage.SaveEntry(entry.UID, entry.FirstName, entry.LastName)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func setBusy() {
	busy = true
}

func setFree() {
	busy = false
}
