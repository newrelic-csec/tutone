package schema

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// filterDescription uses a regex to parse certain data out of the
// description of an item
func filterDescription(description string) string {
	var ret string

	re := regexp.MustCompile(`(?s)(.*)\n---\n`)
	desc := re.FindStringSubmatch(description)

	log.Tracef("description: %#v", desc)

	if len(desc) > 1 {
		ret = desc[1]
	} else {
		ret = description
	}

	return strings.TrimSpace(ret)
}

// typeNameInTypes determines if a name is already present in a set of TypeInfo.
func typeNameInTypes(s string, types []TypeInfo) bool {
	for _, t := range types {
		if t.Name == s {
			return true
		}
	}

	return false
}

// hasType determines if a Type is already present in a slice of Type objects.
func hasType(t *Type, types []*Type) bool {
	for _, tt := range types {
		if t.Name == tt.Name {
			return true
		}
	}

	return false
}

// ExpandType receives a Type which is used to determine the Type for all
// nested fields.
func ExpandType(s *Schema, t *Type) (*[]*Type, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to expand type from nil schema")
	}

	if t == nil {
		return nil, fmt.Errorf("unable to expand nil type")
	}

	var f []*Type

	// Collect the nested types from InputFields
	for _, i := range t.InputFields {
		if i.Type.OfType != nil {
			result, err := s.LookupTypeByName(i.Type.OfType.GetTypeName())
			if err != nil {
				log.Error(err)
			}

			if result != nil {
				f = append(f, result)
			}
		}
	}

	// Same as above, but for Fields
	for _, i := range t.Fields {
		if i.Type.OfType != nil {
			result, err := s.LookupTypeByName(i.Type.OfType.GetTypeName())
			if err != nil {
				log.Error(err)
			}

			if result != nil {
				f = append(f, result)
			}
		}
	}

	return &f, nil
}

// ExpandTypes receives a set of TypeInfo, which is then expanded to include
// all the nested types from the fields.
func ExpandTypes(s *Schema, types []TypeInfo) (*[]*Type, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to expand types from nil schema")
	}

	var expandedTypes []*Type

	for _, schemaType := range s.Types {
		if schemaType != nil {

			// Match the name of types we've resolve and append them to the list
			if typeNameInTypes(schemaType.GetName(), types) {
				expandedTypes = append(expandedTypes, schemaType)

				fieldTypes, err := ExpandType(s, schemaType)
				if err != nil {
					log.Error(err)
				}

				// Avoid duplicates, append the unique names to the set
				for _, f := range *fieldTypes {
					if !hasType(f, expandedTypes) {
						expandedTypes = append(expandedTypes, f)
					}
				}
			}
		}
	}

	return &expandedTypes, nil
}
