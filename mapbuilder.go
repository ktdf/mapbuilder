package mapbuilder

import (
	"github.com/ktdf/parser"
	"net/http"
	"net/url"
	"regexp"
)

func CollectUrls(siteUrl string, maxDepth int) (l *Links, err error) {
	urls := make(Links)
	strippedUrl, _, _, err := getStrippedHostname(siteUrl)
	if err != nil {
		return l, err
	}
	urls.AddLink(strippedUrl)
	l, err = someInternalFunc(&urls, siteUrl, 1, maxDepth)

	return l, err
}

func someInternalFunc(l *Links, siteUrl string, depth int, maxDepth int) (*Links, error) {
	strippedUrl, _, _, err := getStrippedHostname(siteUrl)
	if err != nil {
		return l, err
	}
	response, err := http.Get(siteUrl)
	if err != nil {
		return l, err
	}
	pagelinks := parser.ParseLinks(response.Body)
	for _, link := range pagelinks {
		strippedLink, strippedPath, strippedSchema, err := getStrippedHostname(link.Href)
		if err != nil {
			return l, err
		}
		if strippedSchema == "" || strippedSchema == "http" || strippedSchema == "https" {
			if strippedLink == "" {
				if strippedPath == "/" {
					strippedPath = ""
				}
				l.AddChild(strippedUrl, strippedUrl+strippedPath)
			} else {
				ok, err := regexp.Match(`^(www\.)?`+strippedUrl, []byte(strippedLink))
				if err != nil {
					return l, err
				}
				if ok {
					l.AddChild(strippedUrl+strippedPath, strippedUrl)
				}
			}
		}
	}

	return l, nil

}

//Link is a struct of Ulrs. Could have more than one parent and more than one child.
//Depth - lowest depth of the links
type Link struct {
	parent   bool
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
	var alreadyIn bool
	if (*l)[s] == nil {
		l.AddLink(s)
	}
	if (*l)[child] == nil {
		l.AddLink(child)
	}
	for _, childIterator := range (*l)[s].Children {
		if (*l)[child] == childIterator {
			alreadyIn = true
		}
	}
	if !alreadyIn {
		(*l)[s].Children = append((*l)[s].Children, (*l)[child])
		(*l)[child].Parents = append((*l)[child].Parents, (*l)[s])
		if (*l)[child].Depth == 0 || (*l)[s].Depth+1 < (*l)[child].Depth {
			(*l)[child].Depth = (*l)[s].Depth + 1
			(*l)[child].DepthUpdate()
		}
	}
	return nil
}

//DepthUpdate function could make help.
//It possible to call itself recursively more than once. Don't like it
func (l *Link) DepthUpdate() error {
	for _, child := range l.Children {
		if l.Depth+1 < child.Depth {
			child.Depth = l.Depth + 1
			child.DepthUpdate()
		}
	}
	return nil
}


//getStrippedHostname is used to strip hostname and path to work with it without any crap
func getStrippedHostname(s string) (string, string, string, error) {
	parsedUrl, err := url.Parse(s)
	if err != nil {
		return "", "", "", err
	}
	return parsedUrl.Hostname(), parsedUrl.Path, parsedUrl.Scheme, nil
}
