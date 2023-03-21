package treasure

import (
	"encoding/xml"
	"net/http"
	"strings"
	"treasure/internal/rest/model"
)

const (
	individualType   string = "Individual"
	sdnEntry         string = "sdnEntry"
	strongSearchType string = "strong"
)

var busy bool

// Entry represents sdnEntry node
type Entry struct {
	UID       int    `xml:"uid"`
	FirstName string `xml:"firstName"`
	LastName  string `xml:"lastName"`
	SdnType   string `xml:"sdnType"`
}

// State for internal service state
type State struct {
	Code   int
	Info   string
	Result string
}

// Storage defines methods for storage operations
type Storage interface {
	StrongGetNames(searchName string) (*model.Persons, error)
	WeakGetNames(searchName string) (*model.Persons, error)
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
	var (
		err     error
		persons *model.Persons
	)

	switch strings.ToLower(searchType) {
	case strongSearchType:
		persons, err = s.storage.StrongGetNames(searchName)
	default:
		persons, err = s.storage.WeakGetNames(searchName)

	}
	if err != nil {
		return nil, err
	}

	return persons, nil
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
			if inElement == sdnEntry {
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
