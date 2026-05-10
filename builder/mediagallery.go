package builder

import "github.com/streame-gg/go-discord-wrapper/types/components"

// MediaGalleryBuilder builds a MediaGalleryComponent (Components v2) using a fluent API.
//
//	gallery := builder.NewMediaGallery().
//	    AddItems(
//	        builder.NewMediaGalleryItem("https://example.com/a.png").SetDescription("First").Build(),
//	        builder.NewMediaGalleryItem("https://example.com/b.png").Build(),
//	    ).
//	    Build()
type MediaGalleryBuilder struct {
	mg components.MediaGalleryComponent
}

func NewMediaGallery() *MediaGalleryBuilder {
	items := []components.MediaGalleryItem{}
	return &MediaGalleryBuilder{mg: components.MediaGalleryComponent{Items: &items}}
}

func (b *MediaGalleryBuilder) AddItems(items ...components.MediaGalleryItem) *MediaGalleryBuilder {
	*b.mg.Items = append(*b.mg.Items, items...)
	return b
}

func (b *MediaGalleryBuilder) Build() *components.MediaGalleryComponent {
	return &b.mg
}

// MediaGalleryItemBuilder builds a single MediaGalleryItem.
//
//	item := builder.NewMediaGalleryItem("https://example.com/img.png").SetDescription("A photo").Build()
type MediaGalleryItemBuilder struct {
	item components.MediaGalleryItem
}

func NewMediaGalleryItem(url string) *MediaGalleryItemBuilder {
	return &MediaGalleryItemBuilder{item: components.MediaGalleryItem{
		Media: &components.UnfurledMediaItem{URL: url},
	}}
}

func (b *MediaGalleryItemBuilder) SetDescription(desc string) *MediaGalleryItemBuilder {
	b.item.Description = desc
	return b
}

func (b *MediaGalleryItemBuilder) SetSpoiler(spoiler bool) *MediaGalleryItemBuilder {
	b.item.Spoiler = spoiler
	return b
}

func (b *MediaGalleryItemBuilder) Build() components.MediaGalleryItem {
	return b.item
}
