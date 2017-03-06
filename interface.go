package pixel

import "image/color"

// Target is something that can be drawn onto, such as a window, a canvas, and so on.
//
// You can notice, that there are no "drawing" methods in a Target. That's because all drawing
// happens indirectly through Triangles instance generated via MakeTriangles method.
type Target interface {
	// MakeTriangles generates a specialized copy of the provided Triangles.
	//
	// When calling Draw method on the returned TargetTriangles, the TargetTriangles will be
	// drawn onto the Target that generated them.
	//
	// Note, that not every Target has to recognize all possible types of Triangles. Some may
	// only recognize TrianglesPosition and TrianglesColor and ignore all other properties (if
	// present) when making new TargetTriangles. This varies from Target to Target.
	MakeTriangles(Triangles) TargetTriangles

	// MakePicture generates a specialized copy of the provided Picture.
	//
	// When calling Draw method on the returned TargetPicture, the TargetPicture will be drawn
	// onto the Target that generated it together with the TargetTriangles supplied to the Draw
	// method.
	MakePicture(Picture) TargetPicture
}

// BasicTarget is a Target with additional basic "adjustment" methods.
type BasicTarget interface {
	Target

	// SetMatrix sets a Matrix that every point will be projected by.
	SetMatrix(Matrix)

	// SetColorMask sets a color that will be multiplied with the TrianglesColor property of all
	// Triangles.
	SetColorMask(color.Color)
}

// Triangles represents a list of vertices, where each three vertices form a triangle. (First,
// second and third is the first triangle, fourth, fifth and sixth is the second triangle, etc.)
type Triangles interface {
	// Len returns the number of vertices. The number of triangles is the number of vertices
	// divided by 3.
	Len() int

	// SetLen resizes Triangles to len vertices. If Triangles B were obtained by calling Slice
	// method on Triangles A, the relationship between A and B is undefined after calling SetLen
	// on either one of them.
	SetLen(len int)

	// Slice returns a sub-Triangles of this Triangles, covering vertices in range [i, j).
	//
	// If Triangles B were obtained by calling Slice(4, 9) on Triangles A, then A and B must
	// share the same underlying data. Modifying B must change the contents of A in range
	// [4, 9). The vertex with index 0 at B is the vertex with index 4 in A, and so on.
	//
	// Returned Triangles must have the same underlying type.
	Slice(i, j int) Triangles

	// Update copies vertex properties from the supplied Triangles into this Triangles.
	//
	// Properies not supported by these Triangles should be ignored. Properties not supported by
	// the supplied Triangles should be left untouched.
	//
	// The two Triangles need to have the same Len.
	Update(Triangles)

	// Copy creates an exact independent copy of this Triangles (with the same underlying type).
	Copy() Triangles
}

// TargetTriangles are Triangles generated by a Target with MakeTriangles method. They can be drawn
// onto that (no other) Target.
type TargetTriangles interface {
	Triangles

	// Draw draws Triangles onto an associated Target.
	Draw()
}

// TrianglesPosition specifies Triangles with Position property.
type TrianglesPosition interface {
	Triangles
	Position(i int) Vec
}

// TrianglesColor specifies Triangles with Color property.
type TrianglesColor interface {
	Triangles
	Color(i int) NRGBA
}

// TrianglesPicture specifies Triangles with Picture propery.
//
// Note that this represents picture coordinates, not an actual picture.
type TrianglesPicture interface {
	Triangles
	Picture(i int) (pic Vec, intensity float64)
}

// Picture represents a rectangular area of raster data, such as a color. It has Bounds which
// specify the rectangle where data is located.
type Picture interface {
	// Bounds returns the rectangle of the Picture. All data is located witih this rectangle.
	// Querying properties outside the rectangle should return default value of that property.
	Bounds() Rect

	// Slice returns a sub-Picture with specified Bounds.
	Slice(Rect) Picture

	// Original returns the most original Picture (may be itself) that this Picture was created
	// from using Slice-ing.
	//
	// Since the Original and this Picture should share the underlying data and this Picture can
	// be obtained just by slicing the Original, this method can be used for more efficient
	// caching of Pictures.
	Original() Picture
}

// TargetPicture is a Picture generated by a Target using MakePicture method. This Picture can be drawn onto
// that (no other) Target together with a TargetTriangles generated by the same Target.
//
// The TargetTriangles specify where, shape and how the Picture should be drawn.
type TargetPicture interface {
	Picture

	// Draw draws the supplied TargetTriangles (which must be generated by the same Target as
	// this TargetPicture) with this TargetPicture. The TargetTriangles should utilize the data
	// from this TargetPicture in some way.
	Draw(TargetTriangles)
}

// PictureColor specifies Picture with Color property, so that every position inside the Picture's
// Bounds has a color.
//
// Positions outside the Picture's Bounds must return transparent black (NRGBA{R: 0, G: 0, B: 0, A: 0}).
type PictureColor interface {
	Picture
	Color(at Vec) NRGBA
}
