package mapbuilder

import (
	"fmt"
	"github.com/ktdf/parser"
	"net/http"
	"net/url"
	"regexp"
)

func CollectUrls(siteUrl string, maxDepth uint) (l *Links, err error) {
	urls := make(Links)
	strippedUrl, _, _, err := getStrippedHostname(siteUrl)
	if err != nil {
		return l, err
	}
	urls.AddLink(strippedUrl)
	urls[strippedUrl].Depth = 1
	l, err = recurseLinkParse(&urls, strippedUrl,1, maxDepth)
	return l, err
}

func recurseLinkParse(l *Links, strippedUrl string, depth uint, maxDepth uint) (*Links, error) {
	var notDoneYet bool
	fmt.Printf("Working on %v level\n", depth)
	for siteLink, linkData := range *l {
		if linkData.Depth == depth {
			notDoneYet = true
			response, err := http.Get(`http://`+siteLink)
			if err != nil {
				return l, err
			}
			fmt.Printf("Working with %v\n", siteLink)
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
						l.AddChild(siteLink, strippedUrl+strippedPath)
					} else {
						ok, err := regexp.Match(`^(www\.)?`+strippedUrl, []byte(strippedLink))
						if err != nil {
							return l, err
						}
						if ok {
							l.AddChild(siteLink, strippedUrl+strippedPath)
						}
					}
				}
			}
		}
	}
	if depth == maxDepth {
		return l, nil
	}
	depth++
	if notDoneYet {
		l, err := recurseLinkParse(l, strippedUrl, depth, maxDepth)
		if err != nil {
			return l, err
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
