package file

// Size represents file size in bytes
type Size int64

// Bytes returns the file size in bytes
func (s Size) Bytes() int64 {
	return int64(s)
}

// KB returns the file size in kilobytes
func (s Size) KB() float64 {
	return float64(s.Bytes()) / 1024
}

// MB returns the file size in megabytes
func (s Size) MB() float64 {
	return s.KB() / 1024
}

// GB returns the file size in gigabytes
func (s Size) GB() float64 {
	return s.MB() / 1024
}

// TB returns the file size in terabytes
func (s Size) TB() float64 {
	return s.GB() / 1024
}
