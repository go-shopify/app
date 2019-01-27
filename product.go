package shopify

import "strings"

// ProductTag represents a product tag.
type ProductTag string

// ProductTags represents a list of product tags.
type ProductTags []ProductTag

func (t ProductTags) String() string {
	x := make([]string, len(t))

	for i, tag := range t {
		x[i] = string(tag)
	}

	return strings.Join(x, ",")
}

// ParseProductTags parses a product tags string.
//
// Empty tags are automatically removed as they are invalid.
func ParseProductTags(s string) (ProductTags, error) {
	parts := strings.Split(s, ",")

	result := make(ProductTags, 0, len(parts))

	for _, part := range parts {
		part := strings.TrimSpace(part)

		if len(part) > 0 {
			result = append(result, ProductTag(part))
		}
	}

	return result, nil
}
