package layout

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Option struct {
	Debug      bool             // Debug mode or not
	Folder     string           // Base folder for layout and html files, default is 'views'
	Ext        string           // File extension of layout and html, default is '.html'
	LeftDelim  string           // Left delimiter of template action, default is '{{'
	RightDelim string           // Right delimiter of template action, default is '}}'
	Funcs      template.FuncMap // Function map for template
}

type tmplmap map[string]*template.Template

type Layout struct {
	Option
	pages tmplmap
}

func (l *Layout) Render(w http.ResponseWriter, name string, data interface{}) {
	defer func() {
		if e := recover(); e != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	var t *template.Template

	if l.Debug {
		path := name
		if !strings.HasSuffix(name, l.Ext) {
			path += l.Ext
		}
		t = l.loadPage(path, nil)
	} else {
		if strings.HasSuffix(name, l.Ext) {
			name = name[:len(name)-len(l.Ext)]
		}
		t = l.pages[name]
	}

	if t == nil {
		w.WriteHeader(http.StatusNotFound)
	} else if e := t.ExecuteTemplate(w, "", data); e != nil {
		panic(e)
	}
}

func (l *Layout) getFileList(folder string) []string {
	var files []string

	dir, e := os.Open(folder)
	if e != nil {
		return nil
	}
	defer dir.Close()

	fis, e := dir.Readdir(-1)
	if e != nil {
		return nil
	}

	for _, fi := range fis {
		path := filepath.Join(folder, fi.Name())
		if fi.IsDir() {
			list := l.getFileList(path)
			files = append(files, list...)
			continue
		}
		if !strings.HasSuffix(path, l.Ext) {
			continue
		}
		if path, e = filepath.Rel(l.Folder, path); e == nil {
			files = append(files, path)
		}
	}

	return files
}

func (l *Layout) newTemplate() *template.Template {
	return template.New("").Delims(l.LeftDelim, l.RightDelim).Funcs(l.Funcs)
}

func (l *Layout) loadLayout(path string) *template.Template {
	path = filepath.Join(l.Folder, path)
	data, e := ioutil.ReadFile(path)
	if e != nil {
		panic(e)
	}

	t, e := l.newTemplate().Parse(string(data))
	if e != nil {
		panic(e)
	}

	return t
}

func (l *Layout) getLayoutName(data []byte) string {
	prefix, suffix := l.LeftDelim+"/*", "*/"+l.RightDelim
	buf := bytes.NewBuffer(data)
	for {
		line, e := buf.ReadString('\n')
		if e != nil {
			if e == io.EOF {
				break
			}
			panic(e)
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if !strings.HasPrefix(line, prefix) {
			break
		}
		if !strings.HasSuffix(line, suffix) {
			break
		}
		line = strings.TrimSpace(line[len(prefix) : len(line)-len(suffix)])
		if !strings.HasPrefix(line, "layout:") {
			break
		}
		line = strings.TrimSpace(line[7:])

		if strings.HasSuffix(line, l.Ext) {
			line = line[:len(line)-len(l.Ext)]
		}
		if strings.HasSuffix(line, "_layout") {
			line = line[:len(line)-7]
		}

		return line
	}

	return ""
}

func (l *Layout) loadPage(path string, layouts tmplmap) *template.Template {
	path = filepath.Join(l.Folder, path)
	data, e := ioutil.ReadFile(path)
	if e != nil {
		panic(e)
	}

	var t *template.Template

	ln := l.getLayoutName(data)
	if len(ln) == 0 {
		t = l.newTemplate()
	} else if layouts == nil {
		t = l.loadLayout(ln + "_layout" + l.Ext)
	} else if ll := layouts[ln]; ll != nil {
		t, _ = ll.Clone()
	} else {
		panic("cannot find Layout: " + ln)
	}

	if _, e = t.Parse(string(data)); e != nil {
		panic(e)
	}

	return t
}

func (l *Layout) Build(o Option) {
	l.Option = o

	if len(l.Folder) == 0 {
		l.Folder = "views"
	} else {
		l.Folder = filepath.Clean(l.Folder)
	}

	if len(l.Ext) == 0 {
		l.Ext = ".html"
	} else if l.Ext[0] != '.' {
		l.Ext = "." + l.Ext
	}

	if len(l.LeftDelim) == 0 {
		l.LeftDelim = "{{"
	}
	if len(l.RightDelim) == 0 {
		l.RightDelim = "}}"
	}

	files := l.getFileList(l.Folder)

	suffix := "_layout" + l.Ext
	layouts := make(tmplmap)
	for _, f := range files {
		if !strings.HasSuffix(f, suffix) {
			continue
		}
		name := f[:len(f)-len(suffix)]
		name = strings.Replace(name, "\\", "/", -1)
		layouts[name] = l.loadLayout(f)
	}

	l.pages = make(tmplmap)
	for _, f := range files {
		if strings.HasSuffix(f, suffix) {
			continue
		}
		name := f[:len(f)-len(l.Ext)]
		name = strings.Replace(name, "\\", "/", -1)
		l.pages[name] = l.loadPage(f, layouts)
	}
}
