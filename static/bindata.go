package static

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// bindata_read reads the given file from disk. It returns an error on failure.
func bindata_read(path, name string) ([]byte, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset %s at %s: %v", name, path, err)
	}
	return buf, err
}

// bindata_go reads file data from disk. It returns an error on failure.
func bindata_go() ([]byte, error) {
	return bindata_read(
		"/home/cdecker/go/src/github.com/cdecker/kugelblitz/static/bindata.go",
		"bindata.go",
	)
}

// custom_css reads file data from disk. It returns an error on failure.
func custom_css() ([]byte, error) {
	return bindata_read(
		"/home/cdecker/go/src/github.com/cdecker/kugelblitz/static/custom.css",
		"custom.css",
	)
}

// d3_v3_min_js reads file data from disk. It returns an error on failure.
func d3_v3_min_js() ([]byte, error) {
	return bindata_read(
		"/home/cdecker/go/src/github.com/cdecker/kugelblitz/static/d3.v3.min.js",
		"d3.v3.min.js",
	)
}

// favicon_png reads file data from disk. It returns an error on failure.
func favicon_png() ([]byte, error) {
	return bindata_read(
		"/home/cdecker/go/src/github.com/cdecker/kugelblitz/static/favicon.png",
		"favicon.png",
	)
}

// favicon_xcf reads file data from disk. It returns an error on failure.
func favicon_xcf() ([]byte, error) {
	return bindata_read(
		"/home/cdecker/go/src/github.com/cdecker/kugelblitz/static/favicon.xcf",
		"favicon.xcf",
	)
}

// grid_css reads file data from disk. It returns an error on failure.
func grid_css() ([]byte, error) {
	return bindata_read(
		"/home/cdecker/go/src/github.com/cdecker/kugelblitz/static/grid.css",
		"grid.css",
	)
}

// index_html reads file data from disk. It returns an error on failure.
func index_html() ([]byte, error) {
	return bindata_read(
		"/home/cdecker/go/src/github.com/cdecker/kugelblitz/static/index.html",
		"index.html",
	)
}

// jquery_min_js reads file data from disk. It returns an error on failure.
func jquery_min_js() ([]byte, error) {
	return bindata_read(
		"/home/cdecker/go/src/github.com/cdecker/kugelblitz/static/jquery.min.js",
		"jquery.min.js",
	)
}

// main_js reads file data from disk. It returns an error on failure.
func main_js() ([]byte, error) {
	return bindata_read(
		"/home/cdecker/go/src/github.com/cdecker/kugelblitz/static/main.js",
		"main.js",
	)
}

// package_go reads file data from disk. It returns an error on failure.
func package_go() ([]byte, error) {
	return bindata_read(
		"/home/cdecker/go/src/github.com/cdecker/kugelblitz/static/package.go",
		"package.go",
	)
}

// semantic_css reads file data from disk. It returns an error on failure.
func semantic_css() ([]byte, error) {
	return bindata_read(
		"/home/cdecker/go/src/github.com/cdecker/kugelblitz/static/semantic.css",
		"semantic.css",
	)
}

// semantic_min_js reads file data from disk. It returns an error on failure.
func semantic_min_js() ([]byte, error) {
	return bindata_read(
		"/home/cdecker/go/src/github.com/cdecker/kugelblitz/static/semantic.min.js",
		"semantic.min.js",
	)
}

// themes_default_assets_fonts_icons_woff2 reads file data from disk. It returns an error on failure.
func themes_default_assets_fonts_icons_woff2() ([]byte, error) {
	return bindata_read(
		"/home/cdecker/go/src/github.com/cdecker/kugelblitz/static/themes/default/assets/fonts/icons.woff2",
		"themes/default/assets/fonts/icons.woff2",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"bindata.go": bindata_go,
	"custom.css": custom_css,
	"d3.v3.min.js": d3_v3_min_js,
	"favicon.png": favicon_png,
	"favicon.xcf": favicon_xcf,
	"grid.css": grid_css,
	"index.html": index_html,
	"jquery.min.js": jquery_min_js,
	"main.js": main_js,
	"package.go": package_go,
	"semantic.css": semantic_css,
	"semantic.min.js": semantic_min_js,
	"themes/default/assets/fonts/icons.woff2": themes_default_assets_fonts_icons_woff2,
}
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"bindata.go": &_bintree_t{bindata_go, map[string]*_bintree_t{
	}},
	"custom.css": &_bintree_t{custom_css, map[string]*_bintree_t{
	}},
	"d3.v3.min.js": &_bintree_t{d3_v3_min_js, map[string]*_bintree_t{
	}},
	"favicon.png": &_bintree_t{favicon_png, map[string]*_bintree_t{
	}},
	"favicon.xcf": &_bintree_t{favicon_xcf, map[string]*_bintree_t{
	}},
	"grid.css": &_bintree_t{grid_css, map[string]*_bintree_t{
	}},
	"index.html": &_bintree_t{index_html, map[string]*_bintree_t{
	}},
	"jquery.min.js": &_bintree_t{jquery_min_js, map[string]*_bintree_t{
	}},
	"main.js": &_bintree_t{main_js, map[string]*_bintree_t{
	}},
	"package.go": &_bintree_t{package_go, map[string]*_bintree_t{
	}},
	"semantic.css": &_bintree_t{semantic_css, map[string]*_bintree_t{
	}},
	"semantic.min.js": &_bintree_t{semantic_min_js, map[string]*_bintree_t{
	}},
	"themes": &_bintree_t{nil, map[string]*_bintree_t{
		"default": &_bintree_t{nil, map[string]*_bintree_t{
			"assets": &_bintree_t{nil, map[string]*_bintree_t{
				"fonts": &_bintree_t{nil, map[string]*_bintree_t{
					"icons.woff2": &_bintree_t{themes_default_assets_fonts_icons_woff2, map[string]*_bintree_t{
					}},
				}},
			}},
		}},
	}},
}}
