package mapbuilder

import (
	"github.com/ktdf/parser"
	"net/http"
)

func CollectUrls(url string, depth int) (l Links, err error) {
	//var siteMap Links
	//response, err := http.Get(url)
	//if err != nil {
	//	return nil, err
	//}
	//pagelinks := parser.ParseLinks(response.Body)
	//
	//return nil, nil
}

//Link is a struct of Ulrs. Could have more than one parent and more than one child.
//Depth - lowest depth of the links
type Link struct {
	Name     string
	Parents  []*Link
	Children []*Link
	Depth    uint
}

//Links - map of all links from the site
type Links map[string]*Link

//AddLink creates new element to Links if it no exists
func (l *Links) AddLink(s string) error {
	if (*l)[s] == nil {
		link := new(Link)
		link.Name = s
		(*l)[s] = link
	}
	return nil
}

func (l *Links) AddChild(s string, child string) error {
	if (*l)[s] == nil {
		l.AddLink(s)
	}
	if (*l)[child] == nil {
		l.AddLink(child)
	}
	(*l)[s].Children = append((*l)[s].Children, (*l)[child])
	(*l)[child].Parents = append((*l)[s].Parents, (*l)[s])
	if (*l)[child].Depth == 0 || (*l)[s].Depth+1 > (*l)[child].Depth {
		(*l)[child].Depth = (*l)[s].Depth + 1
		(*l)[child].DepthUpdate()
	}
	return nil
}

//DepthUpdate function could make help.
//It possible to call itself recursively more than once. Don't like it
func (l *Link) DepthUpdate() error {
	for _, child := range l.Children {
		if l.Depth+1 > child.Depth {
			child.Depth = l.Depth+1
			child.DepthUpdate()
		}
	}
	return nil
}
