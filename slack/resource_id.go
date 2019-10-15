package slack

import (
	"bytes"
	"github.com/fluxcd/flux/pkg/resource"
)

type ResourceID struct {
	Namespace, Kind, Name string
}

func newResourceID(id resource.ID) *ResourceID {
	namespace, kind, name := id.Components()
	return &ResourceID{namespace, kind, name}
}

func (s *sender) getResourceURI(id resource.ID) string {
	buf := &bytes.Buffer{}
	rsid := newResourceID(id)
	if err := s.resourceURITmpl.Execute(buf, rsid); err != nil {
		s.log.Error("error on rendering resource ID through template", err)
		return id.String()
	}
	return buf.String()
}
