package main

import "github.com/veandco/go-sdl2/sdl"

type Layer struct {
	Name string
	Desc string
	Geometry
	Rect    sdl.Rect
	Surface *sdl.Surface
	Items   []*Drawable
}

func (L *Layer) AddItem(item Drawable) {
	L.Items = append(L.Items, &item)
	item.SetParentSurface(L.Surface)
}

func (L *Layer) AddItems(items []Drawable) {
	for _, item := range items {
		L.AddItem(item)
	}
}

func (L *Layer) Draw() *sdl.Surface {
	for _, item := range L.Items {
		i := (*item)
		if i.IsChanged() {
			i.Draw()
		}
	}
	// L.Surface.Blit(&L.Rect, s, &L.Rect)
	return L.Surface
}

func (L *Layer) GetChanged() bool {
	changed := false
	for _, item := range L.Items {
		i := (*item)
		ch := i.IsChanged()
		if !changed && ch {
			return true
		}
	}
	return changed
}

func (L *Layer) Destroy() {
	L.Surface.Free()
	for _, item := range L.Items {
		i := (*item)
		i.Destroy()
	}
}
